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

	if hosts := env.GetEnvSlice("ETCD_HOSTS", nil); len(hosts) > 0 {
		for i := range c.Upstreams {
			if c.Upstreams[i].Grpc != nil {
				c.Upstreams[i].Grpc.Etcd.Hosts = hosts
			}
		}
	}

	for i := range c.Upstreams {
		if c.Upstreams[i].Grpc == nil {
			continue
		}
		name := c.Upstreams[i].Name
		if name == "" {
			name = c.Upstreams[i].Grpc.Etcd.Key
		}
		if name == "" {
			continue
		}
		env.OverrideRpcClientConf(c.Upstreams[i].Grpc, toEnvPrefix(name))
	}
}

func toEnvPrefix(s string) string {
	b := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			b = append(b, c-'a'+'A')
		} else if c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
			b = append(b, c)
		} else {
			b = append(b, '_')
		}
	}
	return string(b)
}
