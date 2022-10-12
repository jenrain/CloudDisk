package logic

import (
	"context"

	"CloudDIsk/dataService3/internal/svc"
	"CloudDIsk/dataService3/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataService3Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDataService3Logic(ctx context.Context, svcCtx *svc.ServiceContext) *DataService3Logic {
	return &DataService3Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DataService3Logic) DataService3(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
