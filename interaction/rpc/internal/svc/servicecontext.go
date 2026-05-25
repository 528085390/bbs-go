package svc

import (
	"temp/common/db"
	"temp/interaction/rpc/internal/config"
	"temp/post/rpc/postservice"
	"temp/user/userclient"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	Db      *gorm.DB
	UserRpc userclient.User
	PostRpc postservice.PostService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Db:      db.GetDB(),
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		PostRpc: postservice.NewPostService(zrpc.MustNewClient(c.PostRpc)),
	}
}
