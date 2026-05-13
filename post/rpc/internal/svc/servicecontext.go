package svc

import (
	"encoding/json"
	"fmt"
	"temp/common/db"
	"temp/common/models"
	"temp/common/mq"
	"temp/common/mq/events"
	"temp/post/rpc/internal/config"
	"temp/section/rpc/sectionservice"
	"temp/user/userclient"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	RabbitMQ   mq.RabbitMQ
	UserRpc    userclient.User
	SectionRpc sectionservice.SectionService
	Db         *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:     c,
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		SectionRpc: sectionservice.NewSectionService(zrpc.MustNewClient(c.SectionRpc)),
		Db:         db.GetDB(),
		RabbitMQ:   mq.GetRabbitMQ(c.RabbitMQ),
	}
}

func StartConsumer(ctx *ServiceContext) {

	fmt.Println("starting consumer...")
	msgs, err := ctx.RabbitMQ.MQChannel.Consume(
		ctx.RabbitMQ.QueueName,
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
		fmt.Println(eventsList)
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
	fmt.Println(viewEvents)
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
