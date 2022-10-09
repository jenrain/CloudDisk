package file

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"core/tools"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/zeromicro/go-zero/core/logx"
	"path"
	"time"
)

type FileUploadChunkCompleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadChunkCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadChunkCompleteLogic {
	return &FileUploadChunkCompleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadChunkCompleteLogic) FileUploadChunkComplete(req *types.FileUploadChunkCompleteRequest, userIdentity string) (resp *types.FileUploadChunkCompleteResponse, err error) {
	// 告诉Obs文件已上传完成
	// 将req.CosObjects转化为cos.Object
	//fmt.Println("开始操作")
	//fmt.Println("文件的etag数组：", req.ObsObjects)
	obsParts := make([]obs.Part, 0)
	for _, v := range req.ObsObjects {
		obsParts = append(obsParts, obs.Part{
			ETag:       v.ETag,
			PartNumber: v.PartNumber,
		})
	}
	filePath, err := tools.ObsPartUploadComplete(req.Key, req.UploadId, obsParts)
	// 防止报500错误，主协程需要睡眠几秒钟
	time.Sleep(time.Duration(1) * time.Second)
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	// 将文件信息存入repository_pool
	rp := &models.RepositoryPool{
		Identity: tools.GetUUID(),
		Hash:     req.Hash,
		Name:     req.Name,
		Ext:      path.Ext(req.Name),
		Size:     req.Size,
		Path:     filePath,
	}
	err = l.svcCtx.DB.Create(rp).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	// 将文件信息存入user_repository
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
	return &types.FileUploadChunkCompleteResponse{Success: true}, nil
}
