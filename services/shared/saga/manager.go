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
	Serve() error
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

func (m *manager) ExecuteFirstStep(saga Saga) error {
	if saga.Type == "" {
		return errors.New("empty saga type")
	}

	stateMachine, ok := m.stateMachineMap[saga.Type]
	if !ok {
		return errors.New("state machine for saga type " + saga.Type + " not found")
	}

	newSaga, nextCmd, err := stateMachine.ExecuteFirstStep(saga)
	if err != nil {
		return errors.New("invoke first step in state machine: " + err.Error())
	}
	saga = newSaga
	saga.ID = m.uuidGenerator.NewUUID()

	tx := m.sagaRepo.BeginTransaction()
	err = tx.Repo().CreateSaga(&saga)
	if err != nil {
		tx.RollbackTransaction()
		return errors.New("cannot save saga: " + err.Error())
	}

	id := m.uuidGenerator.NewUUID()

	nextCmd.SetID(id)
	nextCmd.SetSagaID(saga.ID)
	saga.LastCommandID = id

	err = tx.Repo().UpdateSaga(&saga)
	if err != nil {
		tx.RollbackTransaction()
		return errors.New("cannot update saga: " + err.Error())
	}

	err = m.producer.Send(nextCmd.Destination(), nextCmd)
	if err != nil {
		tx.RollbackTransaction()
		return errors.New("cannot send saga command: " + err.Error())
	}

	tx.CommitTransaction()
	return nil
}

func (m *manager) Serve() error {
	if m.isServing {
		return errors.New("already started serving")
	}

	m.isServing = true

	replyChan, _, err := m.consumer.Consume(m.replyChannel)
	if err != nil {
		return err
	}
	go func() {
		for d := range replyChan {
			go m.handleReplyDelivery(d)
		}
	}()
	m.logf("Start receiving replies from MessageQueue channel: %s", m.replyChannel)

	commandChan, _, err := m.consumer.Consume(m.commandChannel)
	if err != nil {
		return err
	}
	go func() {
		for d := range commandChan {
			go m.handleCommandDelivery(d)
		}
	}()
	m.logf("Start handling commands from MessageQueue channel: %s", m.commandChannel)

	forever := make(chan interface{})
	<-forever

	return nil // never reach this
}

func (m *manager) handleReplyDelivery(d msg.Delivery) {
	reply, err := msg.ValidateReply(d.Message)
	if err != nil {
		m.logf("Error: invalid saga reply message: %s, reason: %v", genJsonMessage(d.Message), err)
		d.Nack()
		return
	}

	tx := m.sagaRepo.BeginTransaction()
	err = m.handleReply(tx, reply)
	if err != nil {
		m.logf("Error: cannot handle reply message: %s, reason: %v", genJsonMessage(d.Message), err)
		d.Nack()
		tx.RollbackTransaction()
		return
	}

	m.logf("Handled saga reply %s of saga %s", reply.ID(), reply.SagaID())
	d.Ack()
	tx.CommitTransaction()
}

func (m *manager) handleReply(tx Transaction, reply msg.Reply) error {
	if reply.ID() == "" {
		return errors.New("empty reply ID")
	}

	if reply.SagaID() == "" {
		return errors.New("empty reply's saga ID")
	}

	processed := tx.Repo().CheckIfMessageProcessed(reply.ID())
	if processed {
		return nil
	}

	saga := tx.Repo().FindSagaByID(reply.SagaID())
	if saga == nil {
		return errors.New("cannot find saga with ID " + reply.SagaID())
	}

	if saga.EndState {
		return errors.New("current saga reached end state")
	}

	stateMachine, ok := m.stateMachineMap[saga.Type]
	if !ok {
		return errors.New("state machine for saga type " + saga.Type + " not found")
	}

	newSaga, nextCommand, err := stateMachine.Process(*saga, reply)
	*saga = newSaga
	if err != nil {
		return err
	}

	if nextCommand != nil {
		nextCommand.SetID(m.uuidGenerator.NewUUID())
		nextCommand.SetSagaID(saga.ID)
		nextCommand.SetReplyChannel(m.replyChannel)
		saga.LastCommandID = nextCommand.ID()
	}

	err = tx.Repo().UpdateSaga(saga)
	if err != nil {
		return err
	}

	if nextCommand != nil {
		err = m.producer.Send(nextCommand.Destination(), nextCommand)
		if err != nil {
			return err
		}
	}

	err = tx.Repo().RecordMessageAsProcessed(reply.ID())
	if err != nil {
		return err
	}

	return nil
}

func (m *manager) handleCommandDelivery(d msg.Delivery) {
	command, err := msg.ValidateCommand(d.Message)
	if err != nil {
		m.logf("Error: invalid saga command message: %s, reason: %v", d.Message.ID(), err)
		d.Nack()
		return
	}

	tx := m.sagaRepo.BeginTransaction()
	err = m.handleCommand(tx, command)
	if err != nil {
		m.logf("Error: failed to handle command message: %s, reason: %v", d.Message.ID(), err)
		d.Nack()
		tx.RollbackTransaction()
		return
	}

	m.logf("Handled saga command %s of saga %s", command.ID(), command.SagaID())
	d.Ack()
	tx.CommitTransaction()
}

func (m *manager) handleCommand(tx Transaction, command msg.Command) error {
	if command.ID() == "" {
		return errors.New("empty command ID")
	}

	if command.Type() == "" {
		return errors.New("empty command type")
	}

	handler, ok := m.commandHandlerMap[command.Type()]
	if !ok {
		return errors.New("command handler for command type " + command.Type() + " not found")
	}

	processed := tx.Repo().CheckIfMessageProcessed(command.ID())
	if processed {
		return nil
	}

	// TODO: check the return. Maybe it should be reply message
	err := handler(command)
	if err != nil {
		return errors.New("handle command: " + err.Error())
	}

	err = tx.Repo().RecordMessageAsProcessed(command.ID())
	if err != nil {
		return errors.New("record message as processed: " + err.Error())
	}

	return nil
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
