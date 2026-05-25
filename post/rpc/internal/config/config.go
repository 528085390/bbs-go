package config

import (
	"temp/common/env"
	"temp/common/mq"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	RabbitMQ   mq.RabbitMQConf
	UserRpc    zrpc.RpcClientConf
	SectionRpc zrpc.RpcClientConf
}

func (c *Config) LoadFromEnv() {
	env.OverrideRpcServerConf(&c.RpcServerConf)
	c.RabbitMQ.LoadFromEnv()
	env.OverrideRpcClientConf(&c.UserRpc, "USER_RPC")
	env.OverrideRpcClientConf(&c.SectionRpc, "SECTION_RPC")
	env.OverrideRpcClientConf(&c.InteractionRpc, "INTERACTION_RPC")
	env.OverrideRpcClientConf(&c.CommentRpc, "COMMENT_RPC")
}
