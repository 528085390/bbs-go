package types

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type (
	GatewayConf struct {
		rest.RestConf
		Upstreams []Upstream
		Timeout   time.Duration `json:",default=5s"`
	}

	RouteMapping struct {
		Method  string
		Path    string
		RpcPath string
	}

	Upstream struct {
		Name      string `json:",optional"`
		Grpc      zrpc.RpcClientConf
		ProtoSets []string       `json:",optional"`
		Mappings  []RouteMapping `json:",optional"`
	}
)
