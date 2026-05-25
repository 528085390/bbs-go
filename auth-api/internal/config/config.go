package config

import (
	"temp/common/env"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	UserRpc   zrpc.RpcClientConf
	JwtSecret string
}

func (c *Config) LoadFromEnv() {
	env.OverrideRpcServerConf(&c.RpcServerConf)
	env.OverrideRpcClientConf(&c.UserRpc, "USER_RPC")
	c.JwtSecret = env.GetEnv("JWT_SECRET", c.JwtSecret)
}
