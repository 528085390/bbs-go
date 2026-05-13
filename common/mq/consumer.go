package mq

//
//type OrderCreatedEvent struct {
//	OrderId   int64 `json:"orderId"`
//	UserId    int64 `json:"userId"`
//	ProductId int64 `json:"productId"`
//	Quantity  int32 `json:"quantity"`
//}
//
////func startConsumer(ch *amqp.Channel, queueName string) (amqp.Delivery, error) {
////	consumer, err := ch.Consume(
////		queueName, // 队列名称
////		"",        // 消费者唯一标识（自动生成）
////		false,     // autoAck=false，手动确认
////		false,     // exclusive
////		false,     // noLocal
////		false,     // noWait
////		nil,
////	)
////	if err != nil {
////		logx.Errorf("failed to register consumer: %v", err)
////		return nil, err
////	}
////	return consumer, nil
////
////}
//
//// handleOrderCreated 实际业务处理：发送通知
