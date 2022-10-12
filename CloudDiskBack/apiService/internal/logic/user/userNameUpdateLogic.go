package user

import (
	"context"
	"core/errorx"
	"fmt"

	"core/internal/svc"
	"core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserNameUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserNameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserNameUpdateLogic {
	return &UserNameUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserNameUpdateLogic) UserNameUpdate(req *types.UserNameUpdateRequest, userIdentity string) (resp *types.UserNameUpdateResponse, err error) {
	var cnt int
	fmt.Println("name = ?", req.NewName)
	l.svcCtx.DB.Table("user_basic").Where("name = ?", req.NewName).Count(&cnt)
	if cnt > 0 {
		return nil, errorx.NewDefaultError("昵称已被占用！")
	}
	err = l.svcCtx.DB.Table("user_basic").Where("identity = ?", userIdentity).Update("name", req.NewName).Error
	if err != nil {
		return nil, errorx.NewDefaultError("修改失败！")
	}
	resp = &types.UserNameUpdateResponse{}
	resp.Success = true
	return
}
