package logic

import (
	"context"

	"CloudDIsk/dataService2/internal/svc"
	"CloudDIsk/dataService2/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataService2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDataService2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *DataService2Logic {
	return &DataService2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DataService2Logic) DataService2(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
