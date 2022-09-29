package user

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailLogic {
	return &UserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDetailLogic) UserDetail(req *types.UserDetailRequest, userIdentity string) (resp *types.UserDetailResponse, err error) {
	// 根据identity查找用户
	user := &models.UserBasic{}
	err = l.svcCtx.DB.Where("identity = ?", userIdentity).Find(user).Error
	if err != nil {
		return nil, errorx.NewDefaultError("查询用户信息失败")
	}
	if user.Identity != userIdentity {
		return nil, errorx.NewDefaultError("用户不存在")
	}
	resp = &types.UserDetailResponse{}
	resp.Name = user.Name
	resp.Email = user.Email
	return
}
