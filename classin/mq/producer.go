package mq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
	"meishiedu.com/classin/internal/config"
)

type RabbitProducer struct {
	logx.Logger
	ctx    context.Context
	MqConf config.RabbitMqConf
}

func NewRabbitProducer(ctx context.Context, mqConf config.RabbitMqConf) *RabbitProducer {
	return &RabbitProducer{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		MqConf: mqConf,
	}
}
func (mq *RabbitProducer) SendClassInMsgToMq(topic config.Topic, msg []byte) error {
	conf := amqp.Config{
		Vhost:      "/",
		Properties: amqp.NewConnectionProperties(),
	}
	conn, err := amqp.DialConfig(mq.MqConf.Uri, conf)
	if err != nil {
		mq.Errorf("connect to mq error: %v,conf:%v", err, mq.MqConf)
		return err
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		mq.Errorf("create channel error: %v,conf:%v", err, mq.MqConf)
		return err
	}
	defer channel.Close()
	if err := channel.ExchangeDeclare(
		mq.MqConf.Exchange,
		mq.MqConf.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		mq.Errorf("exchange declare error: %v,conf:%v", err, mq.MqConf)
		return err
	}
	mq.Info(nil, fmt.Sprintf("producer: declaring queue '%s'", topic.Queue))
	queue, err := channel.QueueDeclare(
		topic.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		mq.Errorf("queue declare error: %v,conf:%v", err, mq.MqConf)
		return err
	}

	if err := channel.QueueBind(queue.Name, topic.RoutingKey, mq.MqConf.Exchange, false, nil); err != nil {
		return err
	}
	_, err = channel.PublishWithDeferredConfirm(
		mq.MqConf.Exchange,
		topic.RoutingKey,
		true,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
			AppId:           "beisen-sync-producer",
			Body:            msg,
		},
	)
	if err != nil {
		mq.Errorf("publish error: %v,conf:%v,msg：%v", err, mq.MqConf, msg)
		return err
	}
	return nil
}
