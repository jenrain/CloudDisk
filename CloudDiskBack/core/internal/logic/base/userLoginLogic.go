package base

import (
	"context"
	"core/define"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"core/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 从数据库中查询用户
	user := &models.UserBasic{}
	err = l.svcCtx.DB.Where("name = ? AND password = ?", req.Name, tools.MD5(req.Password)).Find(user).Error
	if err != nil {
		//err = errors.New("查询用户信息失败")
		return nil, errorx.NewDefaultError("查询用户信息失败")
	}
	if user.Name != req.Name {
		//return nil, errors.New("用户名或密码错误")
		return nil, errorx.NewDefaultError("用户名或密码错误")
	}
	// 生成token返回给用户
	token, err := tools.GenerateToken(user.Id, user.Identity, user.Name, define.TokenExpireTime)
	if err != nil {
		return nil, errorx.NewDefaultError("生成token错误")
	}
	// 生成用户刷新token的token
	refreshToken, err := tools.GenerateToken(user.Id, user.Identity, user.Name, define.RefreshTokenExpireTime)
	if err != nil {
		return nil, errorx.NewDefaultError("生成用户刷新的token错误")
	}
	resp = &types.LoginResponse{}
	resp.Token = token
	resp.RefreshToken = refreshToken
	resp.Success = true
	return
}
