package file

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentParentIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrentParentIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentParentIdLogic {
	return &GetCurrentParentIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentParentIdLogic) GetCurrentParentId(req *types.GetCurrentParentIdRequest) (resp *types.GetCurrentParentIdResponse, err error) {
	// 当前目录为根目录
	if req.Id == 0 {
		return &types.GetCurrentParentIdResponse{
			ParentId: 0,
			Success:  true,
		}, nil
	}
	ur := &models.UserRepository{}
	err = l.svcCtx.DB.Table("user_repository").Where("id = ?", req.Id).Find(ur).Error
	if err != nil {
		return nil, errorx.NewDefaultError("查询文件ParentId失败")
	}
	resp = &types.GetCurrentParentIdResponse{}
	resp.ParentId = ur.ParentId
	resp.Success = true
	return
}
