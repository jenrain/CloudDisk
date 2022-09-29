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

type UserFolderCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFolderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFolderCreateLogic {
	return &UserFolderCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFolderCreateLogic) UserFolderCreate(req *types.UserFolderCreateRequest, userIdentity string) (resp *types.UserFolderCreateResponse, err error) {
	// 判断当前文件夹下同名文件夹是否存在
	var count int
	l.svcCtx.DB.Table("user_repository").Where("name = ? AND parent_id = ? AND deleted_time is null", req.Name, req.ParentId).Count(&count)
	if count > 0 {
		return nil, errorx.NewDefaultError("该文件夹已经存在")
	}
	folder := &models.UserRepository{
		Identity:     tools.GetUUID(),
		UserIdentity: userIdentity,
		ParentId:     int(req.ParentId),
		Name:         req.Name,
	}
	err = l.svcCtx.DB.Create(folder).Error
	if err != nil {
		return nil, errorx.NewDefaultError("文件夹创建失败")
	}
	return &types.UserFolderCreateResponse{
		Identity: folder.Identity,
		Success:  true,
	}, nil
}
