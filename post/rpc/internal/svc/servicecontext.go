package svc

import (
	"encoding/json"
	"sync"
	"time"

	"temp/comment/rpc/commentservice"
	"temp/common/db"
	"temp/common/models"
	"temp/common/mq"
	"temp/common/mq/events"
	"temp/interaction/rpc/interactionclient"
	"temp/post/rpc/internal/config"
	"temp/section/rpc/sectionservice"
	"temp/user/userclient"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config         config.Config
	RabbitMQ       *RabbitMQHolder
	UserRpc        userclient.User
	SectionRpc     sectionservice.SectionService
	InteractionRpc interactionclient.Interaction
	CommentRpc     commentservice.CommentService
	Db             *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	rabbitMQHolder := &RabbitMQHolder{}
	initRabbitMQAsync(c.RabbitMQ, rabbitMQHolder)
	return &ServiceContext{
		Config:         c,
		UserRpc:        userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		SectionRpc:     sectionservice.NewSectionService(zrpc.MustNewClient(c.SectionRpc)),
		InteractionRpc: interactionclient.NewInteraction(zrpc.MustNewClient(c.InteractionRpc)),
		CommentRpc:     commentservice.NewCommentService(zrpc.MustNewClient(c.CommentRpc)),
		Db:             db.GetDB(),
		RabbitMQ:       rabbitMQHolder,
	}
}

type RabbitMQHolder struct {
	mu    sync.RWMutex
	mq    mq.RabbitMQ
	ready bool
}

func (h *RabbitMQHolder) Set(client mq.RabbitMQ) {
	h.mu.Lock()
	h.mq = client
	h.ready = true
	h.mu.Unlock()
}

func (h *RabbitMQHolder) Get() (mq.RabbitMQ, bool) {
	h.mu.RLock()
	client := h.mq
	ready := h.ready
	h.mu.RUnlock()
	return client, ready
}

func initRabbitMQAsync(conf mq.RabbitMQConf, holder *RabbitMQHolder) {
	go func() {
		backoff := time.Second
		maxBackoff := 30 * time.Second
		for {
			client, err := mq.GetRabbitMQ(conf)
			if err == nil {
				holder.Set(client)
				logx.Info("RabbitMQ ready")
				return
			}
			logx.Errorf("init RabbitMQ failed, retrying in %s: %v", backoff, err)
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

func StartConsumer(ctx *ServiceContext) {
	logx.Info("starting consumer...")
	var client mq.RabbitMQ
	for {
		var ok bool
		client, ok = ctx.RabbitMQ.Get()
		if ok {
			break
		}
		logx.Info("RabbitMQ not ready, waiting...")
		time.Sleep(2 * time.Second)
	}
	msgs, err := client.MQChannel.Consume(
		client.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logx.Error(err)
		return
	}
	var e events.PostViewEvent
	var eventsList []events.PostViewEvent
	for msg := range msgs {
		body := msg.Body

		err := json.Unmarshal(body, &e)
		if err != nil {
			logx.Error(err)
			err := msg.Nack(false, true)
			if err != nil {
				return
			}
		}
		err = msg.Ack(false)
		eventsList = append(eventsList, e)
		logx.Infof("events: %v", eventsList)
		if len(eventsList) >= 10 {
			err := processEvents(eventsList, ctx)
			if err != nil {
				logx.Error(err)
			}
			eventsList = []events.PostViewEvent{}
		}
	}
	logx.Infof("consumer stopped")
}

func processEvents(viewEvents []events.PostViewEvent, ctx *ServiceContext) error {
	logx.Infof("processing view events: %v", viewEvents)
	views := make(map[int64]int64)
	for _, e := range viewEvents {
		views[e.PostId] += 1
	}
	for postId, count := range views {
		res := ctx.Db.Model(&models.Post{}).Where("id = ?", postId).Update("view_count", gorm.Expr("view_count + ?", count))
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}
