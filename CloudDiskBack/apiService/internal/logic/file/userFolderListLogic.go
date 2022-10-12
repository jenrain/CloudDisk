package file

import (
	"context"
	"core/errorx"
	"fmt"

	"core/internal/svc"
	"core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFolderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFolderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFolderListLogic {
	return &UserFolderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFolderListLogic) UserFolderList(req *types.UserFolderListRequest, userIdentity string) (resp *types.UserFolderListResponse, err error) {
	// 查询该用户的所有文件夹
	folderList := make([]*types.UserFile, 0)
	fmt.Println("开始查询文件夹：")
	err = l.svcCtx.DB.Table("user_repository").Where("user_identity = ? AND ext is null AND deleted_time is null", userIdentity).Select("id, parent_id, name").Find(&folderList).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	//fmt.Println("folderList：", folderList)
	resp = &types.UserFolderListResponse{}
	resp.FolderList = folderList
	resp.Success = true
	return
}
