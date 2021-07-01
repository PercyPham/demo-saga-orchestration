package mq

import (
	"github.com/streadway/amqp"
	"services.shared/apperror"
	"services.shared/saga"
	"services.shared/saga/msg"
)

func NewProducer(conn *amqp.Connection) saga.Producer {
	return &producer{conn}
}

type producer struct {
	conn *amqp.Connection
}

func (p *producer) Send(channel string, message msg.Message) error {
	jsonMsg, err := msg.Marshal(message)
	if err != nil {
		return apperror.Wrap(err, "marshal message into json")
	}

	ch, err := p.conn.Channel()
	if err != nil {
		return apperror.Wrap(err, "open rabbitmq channel")
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
		return apperror.Wrap(err, "declare queue with name "+channel)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonMsg,
		},
	)
	if err != nil {
		return apperror.Wrap(err, "publish message to queue "+channel)
	}

	return nil
}
