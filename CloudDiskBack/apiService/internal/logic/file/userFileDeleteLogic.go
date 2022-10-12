package file

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileDeleteLogic {
	return &UserFileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileDeleteLogic) UserFileDelete(req *types.UserFileDeleteRequest, userIdentity string) (resp *types.UserFileDeleteResponse, err error) {
	// 只删除user_repository中的文件
	err = l.svcCtx.DB.Where("user_identity = ? AND identity = ?", userIdentity, req.Identity).Delete(&models.UserRepository{}).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	resp = &types.UserFileDeleteResponse{Success: true}
	return
}
