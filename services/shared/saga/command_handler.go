package saga

import (
	"fmt"
	"services.shared/apperror"
	"services.shared/saga/msg"
)

type CommandHandler interface {
	Handle(commandType string, handler func(HandlerContext) error)
	Serve()

	ReplySuccess(commandID string, reply msg.Reply) error
	ReplyFailure(commandID string, reply msg.Reply) error
}

func NewCommandHandler(c CommandHandlerConfig) (CommandHandler, error) {
	if err := c.validate(); err != nil {
		return nil, apperror.Wrap(err, "validate SagaCommandHandler's config")
	}
	return &commandHandler{
		commandChannel: c.CommandChannel,
		producer:       c.Producer,
		consumer:       c.Consumer,
		messageRepo:    c.MessageRepo,
		handlers:       make(map[string]func(HandlerContext) error),
		uuidGenerator:  new(defaultUUIDGenerator),
		logger:         new(defaultLogger),
	}, nil
}

type CommandHandlerConfig struct {
	CommandChannel string
	Producer       Producer
	Consumer       Consumer
	MessageRepo    MessageRepo
}

func (c *CommandHandlerConfig) validate() error {
	if c.CommandChannel == "" {
		return apperror.New("CommandChannel in SagaCommandHandler's config must not be empty")
	}
	if c.Producer == nil {
		return apperror.New("Producer in SagaCommandHandler's config must not be nil")
	}
	if c.Consumer == nil {
		return apperror.New("Consumer in SagaCommandHandler's config must not be nil")
	}
	if c.MessageRepo == nil {
		return apperror.New("MessageRepo in SagaCommandHandler's config must not be nil")
	}
	return nil
}

type commandHandler struct {
	commandChannel string
	producer       Producer
	consumer       Consumer
	messageRepo    MessageRepo
	handlers       map[string]func(HandlerContext) error
	uuidGenerator  UUIDGenerator
	logger         Logger
	isServing      bool
}

func (ch *commandHandler) Handle(commandType string, handler func(HandlerContext) error) {
	if commandType == "" {
		panic("empty command type")
	}
	if _, ok := ch.handlers[commandType]; ok {
		panic("duplicate handler for command type " + commandType)
	}
	ch.handlers[commandType] = handler
}

func (ch *commandHandler) Serve() {
	if ch.isServing {
		panic("saga manager serve run more than one")
	}
	ch.isServing = true
	commandChan, _, err := ch.consumer.Consume(ch.commandChannel)
	if err != nil {
		panic("cannot handle saga commands: " + err.Error())
	}
	ch.logf("Start handling commands from MessageQueue channel: %s", ch.commandChannel)
	for d := range commandChan {
		go ch.handleCommandDelivery(d)
	}
}

func (ch *commandHandler) handleCommandDelivery(d msg.Delivery) {
	command, err := msg.ValidateCommand(d.Message)
	if err != nil {
		ch.logf("Error: invalid saga command message: %s, reason: %v", d.Message.ID(), err)
		d.Nack()
		return
	}

	err = ch.handleCommand(command)
	if err != nil {
		ch.logf("Error: failed to handle command %s:%s, reason: %v", command.Type(), command.ID(), err)
		d.Nack()
		return
	}

	ch.logf("Handled saga command %s:%s of saga %s", command.Type(), command.ID(), command.SagaID())
	d.Ack()
}

func (ch *commandHandler) handleCommand(command msg.Command) error {
	if command.ID() == "" {
		return apperror.New("empty command ID")
	}

	if command.Type() == "" {
		return apperror.New("empty command type")
	}

	handler, ok := ch.handlers[command.Type()]
	if !ok {
		return apperror.Newf("command handler for command type %s not found", command.Type())
	}

	processedMessage := ch.messageRepo.GetProcessedMessageByID(command.ID())
	if processedMessage != nil {
		return nil
	}

	handlerCtx := HandlerContext{
		Command:       command,
		producer:      ch.producer,
		uuidGenerator: ch.uuidGenerator,
	}

	if err := handler(handlerCtx); err != nil {
		return apperror.Wrap(err, "handle command")
	}

	if err := ch.messageRepo.CreateProcessedMessage(command); err != nil {
		return apperror.Wrap(err, "record message")
	}

	return nil
}

func (ch *commandHandler) ReplySuccess(commandID string, reply msg.Reply) error {
	reply.SetSuccess(true)
	return ch.reply(commandID, reply)
}

func (ch *commandHandler) ReplyFailure(commandID string, reply msg.Reply) error {
	reply.SetSuccess(false)
	return ch.reply(commandID, reply)
}

func (ch *commandHandler) reply(commandID string, reply msg.Reply) error {
	message := ch.messageRepo.GetProcessedMessageByID(commandID)
	if message == nil {
		return apperror.New("cannot reply to unrecorded command id " + commandID)
	}

	command, err := msg.ValidateCommand(message)
	if err != nil {
		return apperror.Wrap(err, "validate command")
	}

	reply.SetID(ch.uuidGenerator.NewUUID())
	reply.SetCommandID(commandID)
	reply.SetSagaID(command.SagaID())

	err = ch.producer.Send(command.ReplyChannel(), reply)
	if err != nil {
		return apperror.Wrap(err, "send reply")
	}

	return nil
}

func (ch *commandHandler) logf(format string, args ...interface{}) {
	ch.log(fmt.Sprintf(format, args...))
}

func (ch *commandHandler) log(args ...interface{}) {
	args = append([]interface{}{"[SagaCommandHandler]"}, args...)
	ch.logger.Log(args...)
}
