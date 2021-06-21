package saga

import (
	"errors"
	"fmt"
	"services.shared/saga/msg"
)

type Manager interface {
	Handle(commandType string, handler CommandHandler)
	Register(StateMachine)
	ExecuteFirstStep(Saga) error
	ReplySuccess(commandID string, reply msg.Reply) error
	ReplyFailure(commandID string, reply msg.Reply) error
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
		commandChannel:    config.CommandChannel,
		replyChannel:      config.ReplyChannel,
		isServing:         false,
		commandHandlerMap: make(map[string]CommandHandler),
		stateMachineMap:   make(map[string]StateMachine),
		logger:            new(defaultLogger),
		uuidGenerator:     new(defaultUUIDGenerator),
	}, nil
}

type Config struct {
	SagaRepo       Repo
	Producer       Producer
	Consumer       Consumer
	CommandChannel string
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
	if c.CommandChannel == "" {
		return errors.New("CommandChannel must not be empty")
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
	commandChannel    string
	replyChannel      string
	isServing         bool
	logger            Logger
	uuidGenerator     UUIDGenerator
	commandHandlerMap map[string]CommandHandler
	stateMachineMap   map[string]StateMachine
}

func (m *manager) SetLogger(logger Logger) {
	m.logger = logger
}

func (m *manager) SetUUIDGenerator(uuidGenerator UUIDGenerator) {
	m.uuidGenerator = uuidGenerator
}

func (m *manager) Handle(commandType string, handler CommandHandler) {
	if commandType == "" {
		panic("empty command type")
	}
	if _, ok := m.commandHandlerMap[commandType]; ok {
		panic("duplicate command handler for command type: " + commandType)
	}
	m.commandHandlerMap[commandType] = handler
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

	go m.serveHandlingReplies()

	go m.serveHandlingCommands()

	forever := make(chan interface{})
	<-forever
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
