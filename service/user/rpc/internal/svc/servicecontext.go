package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-mall/service/user/model"
	"go-mall/service/user/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	// 增加数据库
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 获取连接
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(conn, c.CacheRedis),
	}
}
