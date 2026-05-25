package config

import (
	"temp/common/env"

	"github.com/zeromicro/go-zero/gateway"
)

type PublicRoute struct {
	Method string `json:",optional"`
	Path   string `json:",optional"`
}

type Config struct {
	gateway.GatewayConf
	JwtSecret    string        `json:",optional"`
	PublicRoutes []PublicRoute `json:",optional"`
}

func (c *Config) LoadFromEnv() {
	c.JwtSecret = env.GetEnv("JWT_SECRET", c.JwtSecret)
}
