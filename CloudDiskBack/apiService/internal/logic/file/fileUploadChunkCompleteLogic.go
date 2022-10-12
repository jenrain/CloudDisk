package file

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"core/mq"
	"core/tools"
	"encoding/json"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"path"
	"time"
)

type FileUploadChunkCompleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadChunkCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadChunkCompleteLogic {
	return &FileUploadChunkCompleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadChunkCompleteLogic) FileUploadChunkComplete(req *types.FileUploadChunkCompleteRequest, userIdentity string, r *http.Request) (resp *types.FileUploadChunkCompleteResponse, err error) {
	fmt.Println("上传的文件uuid：", req.FileIdentity)
	obsParts := make([]obs.Part, 0)
	for _, v := range req.ObsObjects {
		obsParts = append(obsParts, obs.Part{
			ETag:       v.ETag,
			PartNumber: v.PartNumber,
		})
	}

	var node models.Obs
	json.Unmarshal([]byte(req.Node), &node)
	fmt.Println("向节点：", node.HuaweiObsBucket, "上传文件成功")
	url, err := tools.ObsPartUploadComplete(req.Key, req.UploadId, obsParts, node)
	// 获取剩下的节点
	nodes := mq.ChooseRandomDataServersExcept(node)
	// 开始异步向剩下的节点上传文件
	go tools.ObsPartUploadTheRestPart(req.Key, url, node, nodes)
	// 防止报500错误，主协程需要睡眠几秒钟
	time.Sleep(time.Duration(1) * time.Second)
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	// 记录文件上传到了哪些节点
	nodeList := node.HuaweiObsBucket + ","
	for _, o := range nodes {
		nodeList += o.HuaweiObsBucket + ","
	}
	// 将文件信息存入repository_pool
	rp := &models.RepositoryPool{
		Identity: req.FileIdentity,
		Hash:     req.Hash,
		Name:     req.Name,
		Ext:      path.Ext(req.Name),
		Size:     req.Size,
		Node:     nodeList,
	}
	err = l.svcCtx.DB.Create(rp).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	// 将文件信息存入user_repository
	ur := &models.UserRepository{
		Identity:           tools.GetUUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: rp.Identity,
		Ext:                rp.Ext,
		Name:               rp.Name,
	}
	err = l.svcCtx.DB.Table("user_repository").Create(ur).Error
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	return &types.FileUploadChunkCompleteResponse{Success: true}, nil
}
