package mq

import "fmt"

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

func (c RabbitMQConf) GetAMQPUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		c.User, c.Password, c.Host, c.Port, c.VHost)
}
