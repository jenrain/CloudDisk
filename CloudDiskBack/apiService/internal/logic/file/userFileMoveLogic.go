package file

import (
	"context"
	"core/internal/svc"
	"core/internal/types"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileMoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileMoveLogic {
	return &UserFileMoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileMoveLogic) UserFileMove(req *types.UserFileMoveRequest) (resp *types.UserFileMoveResponse, err error) {
	// 更新记录的 ParentId
	err = l.svcCtx.DB.Table("user_repository").Where("id = ?", req.Id).Update("parent_id", req.ParentId).Error
	if err != nil {
		return nil, errors.New("更新失败")
	}
	resp = &types.UserFileMoveResponse{Success: true}
	return
}
