package svc

import (
	"temp/comment/rpc/internal/config"
	"temp/common/db"
	"temp/post/rpc/postservice"
	"temp/user/userclient"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	Db          *gorm.DB
	PostService postservice.PostService
	UserService userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		Db:          db.GetDB(),
		PostService: postservice.NewPostService(zrpc.MustNewClient(c.PostRpc)),
		UserService: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
