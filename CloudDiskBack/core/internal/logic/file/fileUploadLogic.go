package file

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"core/tools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"path"
	"strconv"
	"time"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest, r *http.Request) (resp *types.FileUploadResponse, err error) {
	// 解析文件
	f, header, err := r.FormFile("file")
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	fmt.Println("开始解析文件：", header.Filename)
	// 获取parent_id
	r.ParseForm()
	parentId, err := strconv.Atoi(r.FormValue("parent_id"))
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}

	// 判断当前目录下是否有同名文件
	exist := IsFileExistInUserRepository(header.Filename, parentId, l.svcCtx)
	if exist {
		return nil, errorx.NewDefaultError("文件已存在")
	}

	// 判断repository_pool是否有相同hash的文件，有就直接秒传
	// 计算文件的hash
	hash, err := tools.FileToMD5(f, header.Size)
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	// 从token中获取userIdentity
	userIdentity := r.Header.Get("UserIdentity")
	var rp *models.RepositoryPool
	exist, rp = IsFileExistInRepositoryPool(hash, l.svcCtx)
	// 存在，直接秒传
	if exist {
		return storeInUserRepository(userIdentity, rp.Identity, rp.Ext, rp.Name, parentId, l.svcCtx)
	}
	// 不存在，需要上传到OBS
	cosPath, err := tools.ObsUpload(r)
	// 防止报500错误，主协程需要睡眠几秒钟
	time.Sleep(time.Duration(2) * time.Second)
	if err != nil {
		return
	}
	// 存入repository_pool中
	rp = &models.RepositoryPool{
		Identity: tools.GetUUID(),
		Hash:     hash,
		Name:     header.Filename,
		Ext:      path.Ext(header.Filename),
		Size:     header.Size,
		Path:     cosPath,
	}
	err = l.svcCtx.DB.Create(rp).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	// 存入user_repository中
	return storeInUserRepository(userIdentity, rp.Identity, rp.Ext, rp.Name, parentId, l.svcCtx)
}

// 将上传文件记录存入user_repository中
func storeInUserRepository(userIdentity, repositoryPoolIdentity, ext, name string, parentId int, svcCtx *svc.ServiceContext) (resp *types.FileUploadResponse, err error) {
	ur := &models.UserRepository{
		Identity:           tools.GetUUID(),
		UserIdentity:       userIdentity,
		ParentId:           parentId,
		RepositoryIdentity: repositoryPoolIdentity,
		Ext:                ext,
		Name:               name,
	}
	err = svcCtx.DB.Table("user_repository").Create(ur).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	return &types.FileUploadResponse{
		Identity: ur.Identity,
		Ext:      ur.Ext,
		Name:     ur.Name,
		Success:  true,
	}, nil
}

// IsFileExistInRepositoryPool 判断文件在repository_pool中是否存在
func IsFileExistInRepositoryPool(hash string, svcCtx *svc.ServiceContext) (bool, *models.RepositoryPool) {
	var cnt int
	rp := &models.RepositoryPool{}
	svcCtx.DB.Table("repository_pool").Where("hash = ?", hash).Count(&cnt).Find(&rp)
	// 存在
	if cnt > 0 {
		return true, rp
	}
	return false, nil
}

// IsFileExistInUserRepository 判断文件在user_repository是否存在
func IsFileExistInUserRepository(name string, parentId int, svcCtx *svc.ServiceContext) bool {
	var cnt int
	svcCtx.DB.Table("user_repository").Where("parent_id = ? AND name = ? AND deleted_time is null", parentId, name).Count(&cnt)
	// 存在
	if cnt > 0 {
		return true
	}
	return false
}
