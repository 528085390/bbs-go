package svc

import (
	"temp/auth/internal/config"
	"temp/common/db"
	"temp/user/userclient"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	Db      *gorm.DB
	UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Db:      db.GetDB(),
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
