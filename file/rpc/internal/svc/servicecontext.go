package svc

import (
	"sync"
	"temp/common/db"
	"temp/file/rpc/internal/config"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	OSS    *OSSHolder
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	c.OSS.LoadFromEnv()
	ossHolder := &OSSHolder{}
	initOSSAsync(c.OSS, ossHolder)

	return &ServiceContext{
		Config: c,
		OSS:    ossHolder,
		DB:     db.GetDB(),
	}
}

type OSSHolder struct {
	mu     sync.RWMutex
	client *oss.Client
	bucket *oss.Bucket
	ready  bool
}

func (h *OSSHolder) Set(client *oss.Client, bucket *oss.Bucket) {
	h.mu.Lock()
	h.client = client
	h.bucket = bucket
	h.ready = true
	h.mu.Unlock()
}

func (h *OSSHolder) Get() (*oss.Client, *oss.Bucket, bool) {
	h.mu.RLock()
	client := h.client
	bucket := h.bucket
	ready := h.ready
	h.mu.RUnlock()
	return client, bucket, ready
}

func initOSSAsync(conf config.OSSConfig, holder *OSSHolder) {
	go func() {
		backoff := time.Second
		maxBackoff := 30 * time.Second
		for {
			client, err := oss.New(conf.Endpoint, conf.AccessKeyId, conf.AccessKeySecret)
			if err == nil {
				bucket, bErr := client.Bucket(conf.BucketName)
				if bErr == nil {
					holder.Set(client, bucket)
					logx.Info("OSS ready")
					return
				}
				err = bErr
			}
			logx.Errorf("init OSS failed, retrying in %s: %v", backoff, err)
			time.Sleep(backoff)
			if backoff < maxBackoff {
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}
		}
	}()
}
