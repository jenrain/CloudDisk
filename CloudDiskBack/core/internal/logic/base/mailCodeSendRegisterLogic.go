package base

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/tools"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type MailCodeSendRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMailCodeSendRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MailCodeSendRegisterLogic {
	return &MailCodeSendRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MailCodeSendRegisterLogic) MailCodeSendRegister(req *types.MailCodeSendRequest) (resp *types.MailCodeSendResponse, err error) {
	// 判断邮箱是否已经被注册
	var cnt int64
	fmt.Println("email: ", req.Email)
	err = l.svcCtx.DB.Where("email = ?", req.Email).Table("user_basic").Count(&cnt).Error
	if err != nil {
		return nil, errorx.NewDefaultError("查询用户信息失败")
	}
	if cnt > 0 {
		return nil, errorx.NewDefaultError("该邮箱已经被注册")
	}
	code := tools.RandCode()
	fmt.Println("code: ", code)
	// 储存验证码
	err = l.svcCtx.CacheDB.Set(req.Email, code)
	if err != nil {
		return nil, errorx.NewDefaultError("储存验证码失败")
	}
	// 发送验证码
	err = tools.MailSendCode(req.Email, code)
	if err != nil {
		return nil, errorx.NewDefaultError("发送验证码失败")
	}
	resp = &types.MailCodeSendResponse{}
	resp.Success = true
	return
}
