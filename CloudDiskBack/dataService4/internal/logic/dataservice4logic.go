package logic

import (
	"context"

	"CloudDIsk/dataService4/internal/svc"
	"CloudDIsk/dataService4/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataService4Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDataService4Logic(ctx context.Context, svcCtx *svc.ServiceContext) *DataService4Logic {
	return &DataService4Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DataService4Logic) DataService4(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
