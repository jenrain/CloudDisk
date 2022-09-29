package file

import (
	"context"
	"core/tools"
	"github.com/tencentyun/cos-go-sdk-v5"

	"core/internal/svc"
	"core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *FileUploadChunkCompleteLogic) FileUploadChunkComplete(req *types.FileUploadChunkCompleteRequest) (resp *types.FileUploadChunkCompleteResponse, err error) {
	// 将req.CosObjects转化为cos.Object
	co := make([]cos.Object, 0)
	for _, v := range req.CosObjects {
		co = append(co, cos.Object{
			ETag:       v.ETag,
			PartNumber: v.PartNumber,
		})
	}
	err = tools.CosPartUploadComplete(req.Name, req.UploadId, co)
	return
}
