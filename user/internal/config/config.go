package config

import (
	"temp/common/env"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
}

func (c *Config) LoadFromEnv() {
	env.OverrideRpcServerConf(&c.RpcServerConf)
}
