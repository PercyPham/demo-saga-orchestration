package rabbitmq

import (
	"fmt"
	"strings"
	"time"

	"github.com/streadway/amqp"
	"services.order/internal/common/config"
)

func Connect(cfg config.RabbitMQConfig) (conn *amqp.Connection, close func() error, err error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port)

	retried := 0
	maxRetry := 10

	conn, err = amqp.Dial(url)
	for err != nil && is501Err(err) && retried < maxRetry {
		time.Sleep(2 * time.Second)
		retried++
		conn, err = amqp.Dial(url)
	}
	if err != nil {
		return nil, nil, err
	}

	return conn, conn.Close, nil
}

// 591 EOF error happens when rabbitmq is not completely startup
func is501Err(err error) bool {
	return strings.Contains(err.Error(), "(501)")
}
