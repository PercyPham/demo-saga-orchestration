package pubsub

import (
	"github.com/percypham/saga-go"
	"github.com/percypham/saga-go/msg"
	"github.com/streadway/amqp"
)

func NewSubscriber(conn *amqp.Connection) saga.PSSubscriber {
	return &subscriber{conn}
}

type subscriber struct {
	conn *amqp.Connection
}

func (s *subscriber) Subscribe(topic string) (mChan <-chan msg.Delivery, close func() error, err error) {
	ch, err := s.conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	q, err := ch.QueueDeclare(
		"xemmenu-order-service-"+topic, // name
		true,                           // durable
		false,                          // delete when unused
		false,                          // exclusive
		false,                          // no-wait
		nil,                            // arguments
	)
	if err != nil {
		return nil, nil, err
	}

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		topic,  // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, nil, err
	}

	messageChan := make(chan msg.Delivery)

	go func() {
		for d := range msgs {
			message, err := msg.Unmarshal(d.Body)
			if err != nil {
				_ = d.Nack(false, false)
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
