package rabbitmq

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"services.order/internal/common/config"
	"services.shared/apperror"
)

func Connect(cfg config.RabbitMQConfig) (outflowConn, inflowConn *amqp.Connection, close func() error, err error) {
	var closeOutflowConn func() error
	var outflowConnErr error

	var closeInflowConn func() error
	var inflowConnErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		outflowConn, closeOutflowConn, outflowConnErr = connect(cfg)
		wg.Done()
	}()

	go func() {
		inflowConn, closeInflowConn, inflowConnErr = connect(cfg)
		wg.Done()
	}()

	wg.Wait()

	if outflowConnErr != nil || inflowConnErr != nil {
		if outflowConnErr == nil {
			_ = closeOutflowConn()
		}
		if inflowConnErr == nil {
			_ = closeInflowConn()
		}
		err := combineTwoErrors("outflowConnErr", outflowConnErr, "inflowConnErr", inflowConnErr)
		return nil, nil, nil, apperror.WithLog(err, "connect RabbitMQ to get MessageQueue")
	}

	close = func() error {
		err1 := closeOutflowConn()
		err2 := closeInflowConn()
		if err1 != nil || err2 != nil {
			err = combineTwoErrors("closeOutflowConn", err1, "closeInflowConn", err2)
			return apperror.WithLog(err, "close message queue's connections")
		}
		return nil
	}

	return outflowConn, inflowConn, close, nil
}

func connect(cfg config.RabbitMQConfig) (conn *amqp.Connection, close func() error, err error) {
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

func combineTwoErrors(err1FuncName string, err1 error, err2FuncName string, err2 error) error {
	errMes := ""
	if err1 != nil {
		errMes += err1FuncName + "() error: " + err1.Error()
	}
	if err2 != nil {
		errMes += err2FuncName + "() error: " + err2.Error()
	}
	if errMes == "" {
		return nil
	}
	return errors.New(errMes)
}
