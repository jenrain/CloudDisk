package file

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"core/tools"
	"github.com/zeromicro/go-zero/core/logx"
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

func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest, userIdentity string) (resp *types.FileUploadResponse, err error) {
	rp := &models.RepositoryPool{
		Identity: tools.GetUUID(),
		Hash:     req.Hash,
		Name:     req.Name,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     req.Path,
	}
	// 存入repository_pool中
	err = l.svcCtx.DB.Create(rp).Error
	if err != nil {
		return nil, errorx.NewDefaultError("存入repository_pool失败")
	}

	// 存入user_repository中
	ur := &models.UserRepository{
		Identity:           tools.GetUUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: rp.Identity,
		Ext:                req.Ext,
		Name:               req.Name,
	}
	err = l.svcCtx.DB.Create(ur).Error
	if err != nil {
		return nil, errorx.NewDefaultError("存入user_repository失败")
	}
	resp = &types.FileUploadResponse{}
	resp.Identity = rp.Identity
	resp.Ext = rp.Ext
	resp.Name = rp.Name
	resp.Success = true
	return
}
