package logic

import (
	"context"

	"CloudDIsk/dataService5/internal/svc"
	"CloudDIsk/dataService5/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataService5Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDataService5Logic(ctx context.Context, svcCtx *svc.ServiceContext) *DataService5Logic {
	return &DataService5Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DataService5Logic) DataService5(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
