package env

import (
	"os"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

func GetEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func GetEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

func GetEnvSlice(key string, defaultVal []string) []string {
	if v := os.Getenv(key); v != "" {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	return defaultVal
}

// OverrideRpcServerConf 用环境变量覆盖 RpcServerConf（服务自身注册）
func OverrideRpcServerConf(c *zrpc.RpcServerConf) {
	if hosts := GetEnvSlice("ETCD_HOSTS", nil); len(hosts) > 0 {
		c.Etcd.Hosts = hosts
	}
	if listenOn := GetEnv("LISTEN_ON", ""); listenOn != "" {
		c.ListenOn = listenOn
	}
}

// OverrideRpcClientConf 用环境变量覆盖 RpcClientConf（调用下游 RPC）
// 环境变量命名规则: {PREFIX}_TARGET -> 直连地址, {PREFIX}_ETCD_HOSTS -> etcd 地址列表
// 若均未设置, 回退到通用 ETCD_HOSTS
func OverrideRpcClientConf(c *zrpc.RpcClientConf, prefix string) {
	targetKey := prefix + "_TARGET"
	if target := GetEnv(targetKey, ""); target != "" {
		c.Target = target
		c.Etcd = discov.EtcdConf{}
		return
	}
	etcdHostsKey := prefix + "_ETCD_HOSTS"
	if hosts := GetEnvSlice(etcdHostsKey, nil); len(hosts) > 0 {
		c.Etcd.Hosts = hosts
	} else if hosts := GetEnvSlice("ETCD_HOSTS", nil); len(hosts) > 0 {
		c.Etcd.Hosts = hosts
	}
}
