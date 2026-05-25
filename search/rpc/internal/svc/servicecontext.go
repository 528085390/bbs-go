package svc

import (
	"temp/post/rpc/postservice"
	"temp/search/rpc/internal/config"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	PostRpc postservice.PostService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		PostRpc: postservice.NewPostService(zrpc.MustNewClient(c.PostRpc)),
	}
}
