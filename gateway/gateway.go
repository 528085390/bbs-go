package main

import (
	"flag"
	"fmt"
	"gateway/middleware"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/gateway"
)

var configFile = flag.String("f", "etc/gateway.yaml", "config file")

func main() {
	flag.Parse()

	var c gateway.GatewayConf
	conf.MustLoad(*configFile, &c)

	gw := gateway.MustNewServer(c,
		gateway.WithMiddleware(middleware.AuthMiddleware))

	defer gw.Stop()

	fmt.Printf("Starting gateway server at %s:%d...\n", c.Host, c.Port)
	gw.Start()
}
