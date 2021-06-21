package saga

import (
	"fmt"
	"services.shared/saga/msg"

	GoogleUUID "github.com/google/uuid"
)

// Publisher represents Publisher of PubSub
type Publisher interface {
	Publish(topic string, event msg.Event) error
}

// Subscriber represents Subscriber of PubSub
type Subscriber interface {
	Subscribe(topic string) (c <-chan msg.Delivery, close func() error, err error)
}

// Producer represents Message Queue's Producer
type Producer interface {
	Send(channel string, message msg.Message) error
}

// Consumer represents Message Queue's Consumer
type Consumer interface {
	Consume(channel string) (dChan <-chan msg.Delivery, close func() error, err error)
}

// Repo defines methods needed for saga to work
type Repo interface {
	CreateSaga(*Saga) error
	UpdateSaga(*Saga) error
	FindSagaByID(id string) *Saga

	CreateProcessedMessage(msg.Message) error
	GetProcessedMessageByID(id string) msg.Message

	BeginTransaction() Transaction
}

// Transaction defines methods related to transaction
type Transaction interface {
	Repo() Repo

	RollbackTransaction()
	CommitTransaction()
}

// Logger represents saga's logger
type Logger interface {
	Log(args ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Log(args ...interface{}) {
	fmt.Println(args...)
}

// UUIDGenerator is used for generating uuid for command and reply messages
type UUIDGenerator interface {
	NewUUID() string
}

type defaultUUIDGenerator struct {
}

func (g *defaultUUIDGenerator) NewUUID() string {
	uuid, _ := GoogleUUID.NewUUID()
	return uuid.String()
}
