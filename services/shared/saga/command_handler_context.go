package saga

import (
	"services.shared/apperror"
	"services.shared/saga/msg"
)

type HandlerContext struct {
	Command       msg.Command
	producer      Producer
	uuidGenerator UUIDGenerator
	hasReplied    bool
}

func (c *HandlerContext) ReplySuccess(reply msg.Reply) error {
	reply.SetSuccess(true)
	return c.reply(reply)
}

func (c *HandlerContext) ReplyFailure(reply msg.Reply) error {
	reply.SetSuccess(false)
	return c.reply(reply)
}

func (c *HandlerContext) reply(reply msg.Reply) error {
	if c.hasReplied {
		return apperror.New("reply multiple times")
	}
	c.hasReplied = true

	reply.SetID(c.uuidGenerator.NewUUID())
	reply.SetCommandID(c.Command.ID())
	reply.SetSagaID(c.Command.SagaID())

	if err := c.producer.Send(c.Command.ReplyChannel(), reply); err != nil {
		return apperror.Wrap(err, "send reply")
	}

	return nil
}
