package file

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"fmt"
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
	// 先删除repository_pool中的文件
	ur := &models.UserRepository{}
	err = l.svcCtx.DB.Table("user_repository").Where("identity = ?", req.Identity).Find(ur).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	fmt.Println("rp_identity: ", ur.RepositoryIdentity)
	err = l.svcCtx.DB.Where("identity = ?", ur.RepositoryIdentity).Delete(&models.RepositoryPool{}).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	// 再删除user_repository中的文件
	err = l.svcCtx.DB.Where("user_identity = ? AND identity = ?", userIdentity, req.Identity).Delete(&models.UserRepository{}).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	resp = &types.UserFileDeleteResponse{Success: true}
	return
}
