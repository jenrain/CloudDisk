package logic

import (
	"context"

	"CloudDIsk/dataService6/internal/svc"
	"CloudDIsk/dataService6/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataService6Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDataService6Logic(ctx context.Context, svcCtx *svc.ServiceContext) *DataService6Logic {
	return &DataService6Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DataService6Logic) DataService6(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
