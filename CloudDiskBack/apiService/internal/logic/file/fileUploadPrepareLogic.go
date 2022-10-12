package file

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"core/mq"
	"core/tools"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"path"
)

type FileUploadPrepareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadPrepareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadPrepareLogic {
	return &FileUploadPrepareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadPrepareLogic) FileUploadPrepare(req *types.FileUploadPrepareRequest, userIdentity string) (resp *types.FileUploadPrepareResponse, err error) {
	// 先判断当前目录下是否有该文件
	exist := IsFileExistInUserRepository(req.Name, req.ParentId, l.svcCtx)
	if exist {
		return nil, errorx.NewDefaultError("文件已存在")
	}
	// 判断公共池中有没有该文件
	var rp *models.RepositoryPool
	exist, rp = IsFileExistInRepositoryPool(req.MD5, l.svcCtx)
	// 存在，直接秒传
	if exist {
		fmt.Println("开始秒传")
		ur := &models.UserRepository{
			Identity:           tools.GetUUID(),
			UserIdentity:       userIdentity,
			ParentId:           req.ParentId,
			RepositoryIdentity: rp.Identity,
			Ext:                rp.Ext,
			Name:               rp.Name,
		}
		err = l.svcCtx.DB.Table("user_repository").Create(ur).Error
		if err != nil {
			return nil, errorx.NewDefaultError(err.Error())
		}
		return nil, errorx.NewDefaultError("文件上传成功")
	}
	// 进行分片上传的准备
	// 生成该文件的uuid
	fileIdentity := tools.GetUUID()
	// 随机选择一个节点
	node := mq.ChooseRandomDataServer()
	key, uploadId, err := tools.ObsInitPartUpload(path.Ext(req.Name), fileIdentity, node)

	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	fmt.Println("返回的文件uuid：", fileIdentity)
	respNode, _ := json.Marshal(node)
	return &types.FileUploadPrepareResponse{
		Node:         string(respNode),
		FileIdentity: fileIdentity,
		UploadId:     uploadId,
		Key:          key,
		Success:      true,
	}, nil
}
