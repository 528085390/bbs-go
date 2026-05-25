package config

import (
	"temp/common/env"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	OSS OSSConfig `json:",optional"`
}

// OSSConfig OSS 配置
type OSSConfig struct {
	Endpoint        string `json:",optional"`
	AccessKeyId     string `json:",optional"`
	AccessKeySecret string `json:",optional"`
	BucketName      string `json:",optional"`
	BaseURL         string `json:",optional"`
	UploadTimeout   string `json:",default=15m"`
}

func (c *Config) LoadFromEnv() {
	env.OverrideRpcServerConf(&c.RpcServerConf)
	c.OSS.LoadFromEnv()
}

func (c *OSSConfig) LoadFromEnv() {
	c.Endpoint = env.GetEnv("OSS_ENDPOINT", c.Endpoint)
	c.AccessKeyId = env.GetEnv("OSS_ACCESS_KEY_ID", c.AccessKeyId)
	c.AccessKeySecret = env.GetEnv("OSS_ACCESS_KEY_SECRET", c.AccessKeySecret)
	c.BucketName = env.GetEnv("OSS_BUCKET_NAME", c.BucketName)
	c.BaseURL = env.GetEnv("OSS_BASE_URL", c.BaseURL)
	c.UploadTimeout = env.GetEnv("OSS_UPLOAD_TIMEOUT", c.UploadTimeout)
}
