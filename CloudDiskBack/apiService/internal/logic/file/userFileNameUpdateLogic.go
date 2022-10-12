package file

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileNameUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileNameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileNameUpdateLogic {
	return &UserFileNameUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileNameUpdateLogic) UserFileNameUpdate(req *types.UserFileNameUpdateRequest, userIdentity string) (resp *types.UserFileNameUpdateResponse, err error) {
	// 判断当前名称是否在当前层级下存在
	var count int
	l.svcCtx.DB.Model(&models.UserRepository{}).Where("name = ? AND parent_id = (SELECT parent_id FROM user_repository ur WHERE ur.identity = ?)", req.Name, req.Identity).Count(&count)
	if count > 0 {
		return nil, errorx.NewDefaultError("该名称已经存在")
	}
	data := &models.UserRepository{Name: req.Name}
	err = l.svcCtx.DB.Table("user_repository").Where("identity = ? AND user_identity = ?", req.Identity, userIdentity).Update(data).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	//// 修改repository_pool中的名称
	//ur := &models.UserRepository{}
	//err = l.svcCtx.DB.Table("user_repository").Where("identity = ?", req.Identity).Find(ur).Error
	//if err != nil {
	//	return nil, errorx.NewDefaultError(err.Error())
	//}
	//err = l.svcCtx.DB.Table("repository_pool").Where("identity = ?", ur.RepositoryIdentity).Update("name", req.Name).Error
	//if err != nil {
	//	return nil, errorx.NewDefaultError(err.Error())
	//}
	resp = &types.UserFileNameUpdateResponse{
		Success: true,
		Name:    req.Name,
	}
	return
}
