package base

import (
	"context"
	"core/errorx"
	"core/internal/svc"
	"core/internal/types"
	"core/models"
	"core/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp *types.UserRegisterResponse, err error) {
	// 判断验证码是否一致
	var code string
	// 从Redis连接池中获取一个redis连接
	conn := l.svcCtx.CacheDB.RedisPool.Get()
	defer conn.Close()
	reply, err := conn.Do("get", req.Email)
	code, _ = reply.(string)
	if err != nil {
		return nil, errorx.NewDefaultError("请先获取验证码")
	}
	if code != req.Code {
		return nil, errorx.NewDefaultError("验证码错误")
	}
	// 判断用户是否存在
	var cnt int64
	err = l.svcCtx.DB.Where("email = ?", req.Email).Table("user_basic").Count(&cnt).Error
	if err != nil {
		return nil, errorx.NewDefaultError("查询用户失败")
	}
	if cnt > 0 {
		return nil, errorx.NewDefaultError("用户已存在")
	}
	// 判断用户名称是否被占用
	l.svcCtx.DB.Where("name = ?", req.Name).Count(&cnt)
	if cnt > 0 {
		return nil, errorx.NewDefaultError("昵称已被占用！")
	}
	// 新建用户
	user := &models.UserBasic{
		Identity: tools.GetUUID(),
		Name:     req.Name,
		Password: tools.StringToMD5(req.Password),
		Email:    req.Email,
	}
	err = l.svcCtx.DB.Create(user).Count(&cnt).Error
	if err != nil {
		return nil, errorx.NewDefaultError("插入用户失败")
	}
	if cnt == 0 {
		return nil, errorx.NewDefaultError("注册失败")
	}
	resp = &types.UserRegisterResponse{}
	resp.Success = true
	return
}
