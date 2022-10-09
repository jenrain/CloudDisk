package file

import (
	"context"
	"core/errorx"
	"core/tools"
	"fmt"
	"net/http"
	"time"

	"core/internal/svc"
	"core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadChunkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadChunkLogic {
	return &FileUploadChunkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadChunkLogic) FileUploadChunk(req *types.FileUploadChunkRequest, r *http.Request) (resp *types.FileUploadChunkResponse, err error) {
	etag, err := tools.ObsPartUpload(r)
	// 防止报500错误，主协程需要睡眠几秒钟
	time.Sleep(time.Duration(1) * time.Second)
	if err != nil {
		return nil, errorx.NewDefaultError("分片上传失败")
	}
	fmt.Println("返回的etag: ", etag)
	return &types.FileUploadChunkResponse{
		Etag:    etag,
		Success: true,
	}, nil
}
