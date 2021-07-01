package pubsub

import (
	"github.com/streadway/amqp"
	"services.shared/apperror"
	"services.shared/saga"
	"services.shared/saga/msg"
)

// TODO: cloudamqp.com/blog/part4-rabbitmq-13-common-errors.html
//  + reduce open/close channel repeatedly too many times

func NewPublisher(conn *amqp.Connection) saga.Publisher {
	return &publisher{conn}
}

type publisher struct {
	conn *amqp.Connection
}

func (p *publisher) Publish(topic string, event msg.Event) error {
	jsonEvent, err := msg.Marshal(event)
	if err != nil {
		return apperror.Wrap(err, "marshal event into json")
	}

	ch, err := p.conn.Channel()
	if err != nil {
		return apperror.Wrap(err, "create channel to RabbitMQ")
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		topic,    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		topic, // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonEvent,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
