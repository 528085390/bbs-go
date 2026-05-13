package config

import (
	"temp/common/mq"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	RabbitMQ   mq.RabbitMQConf
	UserRpc    zrpc.RpcClientConf
	SectionRpc zrpc.RpcClientConf
}
