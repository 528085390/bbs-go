package logic

import (
	"context"
	"encoding/json"
	"temp/common/models"
	"temp/common/mq/events"
	"temp/common/valid"

	"temp/post/rpc/internal/svc"
	"temp/post/rpc/post"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostLogic) GetPost(in *post.IdPathReq) (*post.PostResp, error) {
	// 参数校验
	postId := in.Id
	err := valid.IsValidInt(postId)
	if err != nil {
		logx.Errorf("get post invalid params: %v", err)
		return nil, err
	}

	// 查询
	var postRes post.PostResp
	res := l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).First(&postRes)
	if res.Error != nil {
		logx.Errorf("get post failed: %v", res.Error)
		return nil, res.Error
	}

	// 增加文章浏览数
	err = increasePostView(l, postId)
	if err != nil {
		logx.Errorf("increase post view failed: %v", err)
	}

	logx.Infof("get post success: id=%d", postId)

	// 返回结果
	return &post.PostResp{
		Id:        postRes.Id,
		Title:     postRes.Title,
		Content:   postRes.Content,
		SectionId: postRes.SectionId,
		AuthorId:  postRes.AuthorId,
		ViewCount: postRes.ViewCount,
		Pinned:    postRes.Pinned,
		Featured:  postRes.Featured,
		CreatedAt: postRes.CreatedAt,
		UpdatedAt: postRes.UpdatedAt,
	}, nil

}

func increasePostView(l *GetPostLogic, postId int64) error {
	// 构造事件
	event := events.PostViewEvent{
		PostId: postId,
	}
	body, err := json.Marshal(event)
	if err != nil {
		logx.Errorf("marshal event failed: %v", err)
		return err
	}

	client, ok := l.svcCtx.RabbitMQ.Get()
	if !ok {
		logx.Info("RabbitMQ not ready, skip publish")
		return nil
	}

	// 发布消息到 RabbitMQ
	err = client.MQChannel.PublishWithContext(
		l.ctx,
		client.ExchangeName,
		client.RoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 持久化消息（防止重启丢失）
		},
	)
	if err != nil {
		logx.Errorf("publish message to mq failed: %v", err)
		return err
	}

	logx.Infof("published post view event to mq: %v", event)
	return nil
}
