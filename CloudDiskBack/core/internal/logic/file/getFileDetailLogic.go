package file

import (
	"context"
	"core/errorx"
	"core/models"
	"fmt"

	"core/internal/svc"
	"core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFileDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileDetailLogic {
	return &GetFileDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileDetailLogic) GetFileDetail(req *types.GetFileDetailRequest) (resp *types.GetFileDetailResponse, err error) {
	rp := &models.RepositoryPool{}
	err = l.svcCtx.DB.Table("user_repository").Where("user_repository.id = ?", req.Id).
		Select("repository_pool.name, repository_pool.size, repository_pool.ext, repository_pool.created_time, repository_pool.updated_time").
		Joins("left join repository_pool on user_repository.repository_identity = repository_pool.identity").
		Find(rp).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	fmt.Println("rp.CreatedAt：", rp.CreatedAt, ", rp.CreatedAt：", rp.UpdatedAt)
	resp = &types.GetFileDetailResponse{
		Name:       rp.Name,
		Size:       rp.Size,
		CreateTime: rp.CreatedAt.Format("2006-01-02 15-04-05"),
		UpdateTime: rp.UpdatedAt.Format("2006-01-02 15-04-05"),
		Success:    true,
	}
	switch rp.Ext {
	case ".txt":
		resp.FileType = "文本"
	case ".jpg":
		resp.FileType = "图片"
	case ".png":
		resp.FileType = "图片"
	case ".mp3":
		resp.FileType = "音乐"
	case ".mp4":
		resp.FileType = "视频"
	default:
		resp.FileType = "未知类型"
	}
	return
}
