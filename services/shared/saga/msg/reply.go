package msg

import "errors"

const (
	HeaderReplyType             = "reply_type"
	HeaderReplyOfCommand        = "reply_of_command"
	HeaderReplySagaReplyChannel = "reply_saga_reply_channel"
	HeaderReplyResult           = "reply_result"

	ReplyResultSucceeded = "SUCCEEDED"
	ReplyResultFailed    = "FAILED"
)

type Reply interface {
	Message

	CommandID() (commandID string) // Command returns command ID that this reply is for
	SagaID() string
	Success() bool
	Type() string

	SetCommandID(commandID string)
	SetSagaID(sagaID string)
	SetSuccess(bool)
}

type ReplyMeta struct {
	Type string
}

type reply struct {
	Message
}

func NewReply(meta ReplyMeta, payload string) Reply {
	var result string
	headers := map[string]string{
		HeaderMessageType: TypeReply,
		HeaderReplyResult: result,
		HeaderReplyType:   meta.Type,
	}
	message := NewMessage(headers, payload)
	return &reply{message}
}

func (r *reply) CommandID() string { return r.Message.Header(HeaderReplyOfCommand) }
func (r *reply) SagaID() string    { return r.Message.Header(HeaderSagaID) }
func (r *reply) Success() bool     { return r.Message.Header(HeaderReplyResult) == ReplyResultSucceeded }
func (r *reply) Type() string      { return r.Message.Header(HeaderReplyType) }

func (r *reply) SetCommandID(commandID string) { r.Message.SetHeader(HeaderReplyOfCommand, commandID) }
func (r *reply) SetSagaID(sagaID string)       { r.Message.SetHeader(HeaderSagaID, sagaID) }
func (r *reply) SetSuccess(isSuccess bool) {
	if isSuccess {
		r.Message.SetHeader(HeaderReplyResult, ReplyResultSucceeded)
	} else {
		r.Message.SetHeader(HeaderReplyResult, ReplyResultFailed)
	}
}

func ValidateReply(message Message) (Reply, error) {
	if message.Type() != TypeReply {
		return nil, errors.New("invalid reply message type")
	}
	return &reply{message}, nil
}
