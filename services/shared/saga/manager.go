package saga

import (
	"errors"
	"fmt"
	"services.shared/saga/msg"
)

type Manager interface {
	Register(StateMachine)
	ExecuteFirstStep(Saga) error
	Serve()
}

func NewManager(config Config) (Manager, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}
	return &manager{
		sagaRepo:          config.SagaRepo,
		producer:          config.Producer,
		consumer:          config.Consumer,
		replyChannel:      config.ReplyChannel,
		isServing:         false,
		stateMachineMap:   make(map[string]StateMachine),
		logger:            new(defaultLogger),
		uuidGenerator:     new(defaultUUIDGenerator),
	}, nil
}

type Config struct {
	SagaRepo       Repo
	Producer       Producer
	Consumer       Consumer
	ReplyChannel   string
}

func (c *Config) validate() error {
	if c.SagaRepo == nil {
		return errors.New("SagaRepo must not be nil")
	}
	if c.Producer == nil {
		return errors.New("Producer must not be nil")
	}
	if c.Consumer == nil {
		return errors.New("Consumer must not be nil")
	}
	if c.ReplyChannel == "" {
		return errors.New("ReplyChannel must not be empty")
	}
	return nil
}

type manager struct {
	sagaRepo          Repo
	producer          Producer
	consumer          Consumer
	replyChannel      string
	isServing         bool
	logger            Logger
	uuidGenerator     UUIDGenerator
	stateMachineMap   map[string]StateMachine
}

func (m *manager) SetLogger(logger Logger) {
	m.logger = logger
}

func (m *manager) SetUUIDGenerator(uuidGenerator UUIDGenerator) {
	m.uuidGenerator = uuidGenerator
}

func (m *manager) Register(stateMachine StateMachine) {
	if stateMachine.SagaType() == "" {
		panic("cannot register state machine with empty saga type")
	}
	if _, ok := m.stateMachineMap[stateMachine.SagaType()]; ok {
		panic("cannot register duplicate state machine with same saga type " + stateMachine.SagaType())
	}
	m.stateMachineMap[stateMachine.SagaType()] = stateMachine
}

func (m *manager) Serve() {
	if m.isServing {
		panic("saga manager serve run more than one")
	}
	m.isServing = true

	m.serveHandlingReplies()
}

func (m *manager) logf(format string, args ...interface{}) {
	m.log(fmt.Sprintf(format, args...))
}

func (m *manager) log(args ...interface{}) {
	args = append([]interface{}{"[SagaManager]"}, args...)
	m.logger.Log(args...)
}

func genJsonMessage(m msg.Message) string {
	jsonMsg, _ := msg.Marshal(m)
	return string(jsonMsg)
}
