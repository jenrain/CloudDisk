package file

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"strings"
)

var (
	HuaweiObsEndPoint         = "https://obs.cn-east-3.myhuaweicloud.com"
	HuaweiObsBucketRootFolder = "cloud-disk"
	ENTER                     = "/"
	DOT                       = "."
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
	filesList := make([]types.UserFile, 0)
	// 文件夹列表
	folderList := make([]types.UserFile, 0)
	// 文件数目
	var filesNum int
	// 文件夹数目
	var folderNum int

	// 先判断缓存里面有没有值
	conn := l.svcCtx.CacheDB.RedisPool.Get()
	parentId := strconv.Itoa(int(req.Id))
	_fileExist, _ := conn.Do("HEXISTS", userIdentity, parentId+"file")
	_folderExist, _ := conn.Do("HEXISTS", userIdentity, parentId+"folder")
	fileExist := _fileExist.(int64)
	folderExist := _folderExist.(int64)
	//fmt.Println(fileExist, folderExist)

	// 缓存中有值
	if fileExist > 0 && folderExist > 0 {
		bytes1, _ := redis.Bytes(conn.Do("HGET", userIdentity, parentId+"file"))
		json.Unmarshal(bytes1, &filesList)
		filesNum = len(filesList)
		bytes2, _ := redis.Bytes(conn.Do("HGET", userIdentity, parentId+"folder"))
		//fmt.Println(string(bytes2))
		json.Unmarshal(bytes2, &folderList)
		folderNum = len(folderList)
		//fmt.Println("fileList: ", filesList)
		//fmt.Println("folderList: ", folderList)
		//fmt.Println("取缓存数据")
		logx.Info("向缓存中取值：", userIdentity, ", ", parentId)
	} else {
		// 查询当前目录下的文件
		err = l.svcCtx.DB.Table("user_repository").Where("parent_id = ? AND user_identity = ? AND user_repository.deleted_time is null AND user_repository.ext is not null", req.Id, userIdentity).
			Select("user_repository.id, user_repository.identity, user_repository.parent_id, user_repository.repository_identity, user_repository.name, user_repository.ext, " +
				"repository_pool.size, repository_pool.node").
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
		// 得到每个文件的URL
		for i := 0; i < len(filesList); i++ {
			// 如果某份文件的Node字段为0，就代表所有节点中都没有该文件的备份
			if len(filesList[i].Node) > 0 {
				node := filesList[i].Node[:strings.Index(filesList[i].Node, ",")]
				filesList[i].Path = HuaweiObsEndPoint[:8] + node + DOT + HuaweiObsEndPoint[8:] + ENTER + HuaweiObsBucketRootFolder + ENTER + filesList[i].RepositoryIdentity + filesList[i].Ext
			}
		}

		// 将数据存入缓存
		//conn := l.svcCtx.CacheDB.RedisPool.Get()
		//parentId := strconv.Itoa(int(req.Id))
		//先存文件列表
		_fileList, _ := json.Marshal(filesList)
		conn.Do("Hset", userIdentity, parentId+"file", string(_fileList))
		// 再存文件夹列表
		_folderList, _ := json.Marshal(folderList)
		conn.Do("Hset", userIdentity, parentId+"folder", string(_folderList))
		logx.Info("向数据库中取值：", userIdentity, ", ", parentId)
	}
	resp = &types.UserFileListResponse{}
	resp.FilesList = filesList
	resp.FilesNum = filesNum
	resp.FolderList = folderList
	resp.FolderNum = folderNum
	resp.Success = true
	return
}
