package file

import (
	"context"
	"core/errorx"
	"core/tools"
	"fmt"

	"core/internal/svc"
	"core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFilePathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFilePathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFilePathLogic {
	return &GetFilePathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFilePathLogic) GetFilePath(req *types.GetFilePathRequest) (resp *types.GetFilePathResponse, err error) {
	path, err := tools.GetFilePath(req.RepositoryIdentity, req.Ext)
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	fmt.Println("获取到的文件路径：", path)
	return &types.GetFilePathResponse{
		Path:    path,
		Success: true,
	}, nil
}
