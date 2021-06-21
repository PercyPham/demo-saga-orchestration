package msg

import "errors"

const (
	HeaderReplyType      = "reply_type"
	HeaderReplyOfCommand = "reply_of_command"
	HeaderReplyResult    = "reply_result"

	ReplyResultSucceeded = "SUCCEEDED"
	ReplyResultFailed    = "FAILED"
)

type Reply interface {
	Message

	SagaID() string
	Command() (commandID string) // Command returns command ID that this reply is for
	Success() bool
	Type() string
}

type ReplyMeta struct {
	Command string
	Success bool
	Type    string
}

type reply struct {
	Message
}

func NewReply(meta ReplyMeta, payload string) Reply {
	var result string
	if meta.Success {
		result = ReplyResultSucceeded
	} else {
		result = ReplyResultFailed
	}

	headers := map[string]string{
		HeaderMessageType:    TypeReply,
		HeaderReplyOfCommand: meta.Command,
		HeaderReplyResult:    result,
		HeaderReplyType:      meta.Type,
	}
	message := NewMessage(headers, payload)
	return &reply{message}
}

func (r *reply) SagaID() string  { return r.Message.GetHeader(HeaderSagaID) }
func (r *reply) Command() string { return r.Message.GetHeader(HeaderReplyOfCommand) }
func (r *reply) Success() bool   { return r.Message.GetHeader(HeaderReplyResult) == ReplyResultSucceeded }
func (r *reply) Type() string    { return r.Message.GetHeader(HeaderReplyType) }

func ValidateReply(message Message) (Reply, error) {
	if message.Type() != TypeReply {
		return nil, errors.New("invalid reply message type")
	}
	return &reply{message}, nil
}
