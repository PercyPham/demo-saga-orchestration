package mq

import (
	"github.com/percypham/saga-go"
	"github.com/percypham/saga-go/msg"
	"github.com/streadway/amqp"
	"services.shared/apperror"
)

func NewConsumer(conn *amqp.Connection) saga.Consumer {
	return &consumer{conn}
}

type consumer struct {
	conn *amqp.Connection
}

func (c *consumer) Consume(channel string) (mChan <-chan msg.Delivery, close func() error, err error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, nil, apperror.WithLog(err, "open rabbitmq channel")
	}

	q, err := ch.QueueDeclare(
		channel, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, nil, apperror.WithLog(err, "declare queue with name "+channel)
	}

	msgQueueChan, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, nil, apperror.WithLog(err, "consume message from channel "+channel)
	}

	messageChan := make(chan msg.Delivery, 1)

	go func() {
		for d := range msgQueueChan {
			message, err := msg.Unmarshal(d.Body)
			if err != nil {
				_ = d.Nack(false, false)
				continue
			}
			ack := func() {
				_ = d.Ack(false)
			}
			nack := func() {
				_ = d.Nack(false, true)
			}
			delivery := msg.Delivery{
				Message: message,
				Ack:     ack,
				Nack:    nack,
			}
			messageChan <- delivery
		}
	}()

	return messageChan, ch.Close, nil
}
