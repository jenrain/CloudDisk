package logic

import (
	"context"

	"CloudDIsk/dataService1/internal/svc"
	"CloudDIsk/dataService1/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataServiceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDataServiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataServiceLogic {
	return &DataServiceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DataServiceLogic) DataService(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
