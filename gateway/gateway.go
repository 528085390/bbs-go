package main

import (
	"context"
	"flag"
	"net/http"
	"strconv"
	"strings"
	"temp/common/errs/errorcode"

	"temp/common/response"
	"temp/gateway/config"
	"temp/gateway/middleware"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

var configFile = flag.String("f", "etc/gateway.yaml", "config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	c.LoadFromEnv()

	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		if st, ok := status.FromError(err); ok && st.Message() != "" {
			errMsg := st.Message()
			msgSlice := strings.Split(errMsg, "|")
			codeMsg := msgSlice[0]
			msg := msgSlice[1]
			data := msgSlice[2]

			code, _ := strconv.Atoi(codeMsg)

			return http.StatusOK, response.Error(code, msg, data)
		}
		return http.StatusOK, response.Error(int(errorcode.ServerError.Code), err.Error(), nil)
	})

	gw := gateway.MustNewServer(c.GatewayConf,
		gateway.WithMiddleware(middleware.AuthMiddleware(c)))

	defer gw.Stop()

	logx.Infof("Starting gateway server at %s:%d...", c.Host, c.Port)
	gw.Start()
}
