package msg

import "errors"

const (
	HeaderSagaID = "saga_id"

	HeaderCommandType         = "command_type"
	HeaderCommandDestination  = "command_destination"
	HeaderCommandReplyChannel = "command_reply_channel"
)

type Command interface {
	Message

	SagaID() string
	Type() string
	Destination() string
	ReplyChannel() string

	SetSagaID(id string)
	SetReplyChannel(channel string)
}

type CommandMeta struct {
	Destination string
	Type        string
}

type command struct {
	Message
}

func NewCommand(meta CommandMeta, payload string) Command {
	headers := map[string]string{
		HeaderMessageType:        TypeCommand,
		HeaderCommandType:        meta.Type,
		HeaderCommandDestination: meta.Destination,
	}
	return &command{NewMessage(headers, payload)}
}

func (c *command) SagaID() string                 { return c.Header(HeaderSagaID) }
func (c *command) Type() string                   { return c.Header(HeaderCommandType) }
func (c *command) Destination() string            { return c.Header(HeaderCommandDestination) }
func (c *command) ReplyChannel() string           { return c.Header(HeaderCommandReplyChannel) }
func (c *command) SetSagaID(id string)            { c.SetHeader(HeaderSagaID, id) }
func (c *command) SetReplyChannel(channel string) { c.SetHeader(HeaderCommandReplyChannel, channel) }

func ValidateCommand(message Message) (Command, error) {
	if message.Type() != TypeCommand {
		return nil, errors.New("invalid command message type")
	}
	return &command{message}, nil
}
