package file

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileListLogic {
	return &UserFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileListLogic) UserFileList(req *types.UserFileListRequest, userIdentity string) (resp *types.UserFileListResponse, err error) {
	// 文件列表
	filesList := make([]*types.UserFile, 0)
	// 文件夹列表
	folderList := make([]*types.UserFile, 0)
	// 文件数目
	var filesNum int
	// 文件夹数目
	var folderNum int
	// 查询当前目录下的文件
	err = l.svcCtx.DB.Table("user_repository").Where("parent_id = ? AND user_identity = ? AND user_repository.deleted_time is null AND user_repository.ext is not null", req.Id, userIdentity).
		Select("user_repository.id, user_repository.identity, user_repository.parent_id, user_repository.repository_identity, user_repository.name, user_repository.ext, " +
			"repository_pool.path, repository_pool.size").
		Joins("left join repository_pool on user_repository.repository_identity = repository_pool.identity").
		Find(&filesList).Error
	if err != nil {
		return nil, errorx.NewDefaultError("查询文件失败")
	}

	// 查询当前目录下的文件夹
	err = l.svcCtx.DB.Table("user_repository").Where("parent_id = ? AND user_identity = ? AND user_repository.deleted_time is null AND user_repository.ext is null", req.Id, userIdentity).Find(&folderList).Error
	if err != nil {
		return nil, errorx.NewDefaultError("查询文件夹失败")
	}
	// 查询列表中的文件数
	l.svcCtx.DB.Table("user_repository").Where("parent_id = ? AND user_identity = ? AND deleted_time is null AND user_repository.ext is not null", req.Id, userIdentity).Count(&filesNum)
	// 查询列表中的文件夹数目
	l.svcCtx.DB.Table("user_repository").Where("parent_id = ? AND user_identity = ? AND deleted_time is null AND user_repository.ext is null", req.Id, userIdentity).Count(&folderNum)
	for i := 0; i < len(filesList); i++ {
		switch filesList[i].Ext {
		case "":
			filesList[1].FileType = "folder"
		case ".txt":
			filesList[i].FileType = "txt"
		case ".jpg":
			filesList[i].FileType = "image"
		case ".png":
			filesList[i].FileType = "image"
		case ".mp3":
			filesList[i].FileType = "audio"
		case ".mp4":
			filesList[i].FileType = "video"
		default:
			filesList[i].FileType = "unknown"
		}
	}
	resp = &types.UserFileListResponse{}
	resp.FilesList = filesList
	resp.FilesNum = filesNum
	resp.FolderList = folderList
	resp.FolderNum = folderNum
	resp.Success = true
	return
}
