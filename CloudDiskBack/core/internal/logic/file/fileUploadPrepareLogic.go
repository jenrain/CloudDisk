package file

import (
	"context"
	"core/models"
	"core/tools"

	"core/internal/svc"
	"core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *FileUploadPrepareLogic) FileUploadPrepare(req *types.FileUploadPrepareRequest) (resp *types.FileUploadPrepareResponse, err error) {
	rp := &models.RepositoryPool{}
	var cnt int
	l.svcCtx.DB.Where("hash = ?", req.MD5).Find(rp).Count(&cnt)
	resp = &types.FileUploadPrepareResponse{}
	// 公共池中有这个文件，可以走秒传
	if cnt > 0 {
		resp.Identity = rp.Identity
	} else {
		// 进行分片上传的准备，获取分片上传的ID
		name, uploadId, err := tools.CosInitPartUpload(req.Ext)
		if err != nil {
			return nil, err
		}
		resp.Name = name
		resp.UploadId = uploadId
	}
	return
}
