package mq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
	"meishiedu.com/classin/internal/config"
	"time"
)

type RabbitMqConsumer struct {
	logx.Logger
	MqConf   config.RabbitMqConf
	Ctx      context.Context
	conn     *amqp.Connection
	ConnName string
	channel  *amqp.Channel
	Topic    config.Topic
	CTag     string
	done     chan *amqp.Error
}
type processHandler func([]byte) error

func NewRabbitMqConsumer(ctx context.Context, mqConf config.RabbitMqConf, topic config.Topic) *RabbitMqConsumer {
	return &RabbitMqConsumer{
		Logger: logx.WithContext(ctx),
		MqConf: mqConf,
		Topic:  topic,
		Ctx:    ctx,
	}
}

func (consumer *RabbitMqConsumer) Consumer(processor processHandler) {
	defer func() {
		if err := recover(); err != nil {
			consumer.Errorf("Consumer encountered an error and needs to be restarted, error: %v", err)
			if consumer.channel != nil {
				_ = consumer.channel.Close()
			}
			if consumer.conn != nil {
				_ = consumer.conn.Close()
			}
			time.Sleep(time.Second * 3)
			consumer.Consumer(processor)
		}
	}()
	var err error
	conf := amqp.Config{Properties: amqp.NewConnectionProperties()}
	conf.Properties.SetClientConnectionName(consumer.ConnName)
	consumer.Infof("dialing %q", consumer.MqConf.Uri)
	consumer.conn, err = amqp.DialConfig(consumer.MqConf.Uri, conf)
	if err != nil {
		consumer.Errorf("Dial: %s", err)
		panic(err)
	}
	consumer.Info("got Connection, getting Channel")
	consumer.channel, err = consumer.conn.Channel()
	if err != nil {
		consumer.Errorf("Channel: %s", err)
		panic(err)
	}
	consumer.Infof("got Channel, declaring Exchange (%q)", consumer.MqConf.Exchange)
	if err = consumer.channel.ExchangeDeclare(
		consumer.MqConf.Exchange,
		consumer.MqConf.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		consumer.Errorf("Exchange Declare: %s", err)
		panic(err)
	}

	consumer.Infof("declared Exchange, declaring Queue %q", consumer.Topic.Queue)
	queue, err := consumer.channel.QueueDeclare(
		consumer.Topic.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		consumer.Errorf("Queue Declare: %s", err)
		panic(err)
	}

	consumer.Infof("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, consumer.Topic.RoutingKey)

	if err = consumer.channel.QueueBind(
		queue.Name,
		consumer.Topic.RoutingKey,
		consumer.MqConf.Exchange,
		false,
		nil,
	); err != nil {
		consumer.Errorf("Queue Bind: %s", err)
		panic(err)
	}

	consumer.Infof("Queue bound to Exchange, starting Consume (consumer tag %q)", consumer.Topic.CTag)
	consumer.done = consumer.channel.NotifyClose(make(chan *amqp.Error, 1))
	deliveries, err := consumer.channel.Consume(
		queue.Name,
		consumer.Topic.CTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		consumer.Errorf("Queue Consume: %s", err)
		panic(err)
	}
	for {
		select {
		case err = <-consumer.done:
			close(consumer.done)
			panic(err)
		case msg := <-deliveries:
			err = processor(msg.Body)
			if err != nil {
				consumer.Error(err.Error())
			} else {
				err = msg.Ack(true)
				if err != nil {
					consumer.Errorf("Ack error，error：%v", err)
				}

			}
		}
	}
}
