package svc

import (
	"core/internal/config"
	"core/internal/middleware"
	"core/models"
	"github.com/jinzhu/gorm"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	CacheDB models.CacheDb
	Auth    rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		DB:      models.InitDB(c.Mysql.DataSource),
		CacheDB: models.InitCacheDB(c),
		Auth:    middleware.NewAuthMiddleware().Handle,
	}
}
