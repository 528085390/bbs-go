package mq

import (
	"fmt"
	"temp/common/env"
)

type RabbitMQConf struct {
	Host         string
	Port         int
	User         string
	Password     string
	VHost        string
	ExchangeName string
	QueueName    string
	RoutingKey   string
}

func (c *RabbitMQConf) LoadFromEnv() {
	c.Host = env.GetEnv("MQ_HOST", c.Host)
	c.Port = env.GetEnvInt("MQ_PORT", c.Port)
	c.User = env.GetEnv("MQ_USER", c.User)
	c.Password = env.GetEnv("MQ_PASSWORD", c.Password)
	c.VHost = env.GetEnv("MQ_VHOST", c.VHost)
	c.ExchangeName = env.GetEnv("MQ_EXCHANGE", c.ExchangeName)
	c.QueueName = env.GetEnv("MQ_QUEUE", c.QueueName)
	c.RoutingKey = env.GetEnv("MQ_ROUTING_KEY", c.RoutingKey)
}

func (c *RabbitMQConf) GetAMQPUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		c.User, c.Password, c.Host, c.Port, c.VHost)
}
