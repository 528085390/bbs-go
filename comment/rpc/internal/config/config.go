package config

import (
	"temp/common/env"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	PostRpc zrpc.RpcClientConf
	UserRpc zrpc.RpcClientConf
}

func (c *Config) LoadFromEnv() {
	env.OverrideRpcServerConf(&c.RpcServerConf)
	env.OverrideRpcClientConf(&c.PostRpc, "POST_RPC")
	env.OverrideRpcClientConf(&c.UserRpc, "USER_RPC")
}
