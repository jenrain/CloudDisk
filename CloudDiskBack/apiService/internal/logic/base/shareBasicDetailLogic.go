package base

import (
	"context"
	"core/internal/svc"
	"core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareBasicDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicDetailLogic {
	return &ShareBasicDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicDetailLogic) ShareBasicDetail(req *types.ShareBasicDetailRequest) (resp *types.ShareBasicDetailResponse, err error) {
	// 对分享记录的点击次数进行 + 1
	err = l.svcCtx.DB.Exec("UPDATE share_basic SET click_num = click_num + 1 WHERE identity = ?", req.Identity).Error
	if err != nil {
		return
	}
	// 获取资源的详细信息
	resp = &types.ShareBasicDetailResponse{}
	err = l.svcCtx.DB.Table("share_basic").Where("share_basic.identity = ?", req.Identity).
		Select("share_basic.repository_identity, user_repository.name, repository_pool.ext, repository_pool.size, repository_pool.path").
		Joins("left join user_repository on share_basic.user_repository_identity = user_repository.identity").
		Joins("left join repository_pool on share_basic.repository_identity = repository_pool.identity").
		Find(resp).Error
	if err != nil {
		return nil, err
	}
	return
}
