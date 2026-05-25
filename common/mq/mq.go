package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type RabbitMQ struct {
	MQChannel    *amqp.Channel
	ExchangeName string
	QueueName    string
	RoutingKey   string
}

// InitRabbitMQ 初始化 RabbitMQ 连接并声明交换机和队列
func InitRabbitMQ(amqpUrl, exchangeName, queueName, routingKey string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 声明持久化的 Topic 类型交换机
	err = ch.ExchangeDeclare(
		exchangeName, // 交换机名称
		"topic",      // 交换机类型
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部
		false,        // 不等待
		nil,
	)
	if err != nil {
		return nil, err
	}

	// 声明持久化队列
	_, err = ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 非独占
		false,     // 不等待
		nil,
	)
	if err != nil {
		return nil, err
	}

	// 绑定队列到交换机
	err = ch.QueueBind(
		queueName,    // 队列名称
		routingKey,   // 路由键
		exchangeName, // 交换机名称
		false,        // 不等待
		nil,
	)
	if err != nil {
		return nil, err
	}

	logx.Infof("RabbitMQ initialized. Exchange: %s, Queue: %s, RoutingKey: %s",
		exchangeName, queueName, routingKey)

	return ch, nil
}

func GetRabbitMQ(c RabbitMQConf) (RabbitMQ, error) {
	amqpUrl := c.GetAMQPUrl()
	ch, err := InitRabbitMQ(amqpUrl, c.ExchangeName,
		c.QueueName, c.RoutingKey)
	if err != nil {
		return RabbitMQ{}, err
	}
	return RabbitMQ{
		MQChannel:    ch,
		ExchangeName: c.ExchangeName,
		QueueName:    c.QueueName,
		RoutingKey:   c.RoutingKey,
	}, nil
}
