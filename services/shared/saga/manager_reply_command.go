package saga

import (
	"services.shared/apperror"
	"services.shared/saga/msg"
)

func (m *manager) ReplySuccess(commandID string, reply msg.Reply) error {
	reply.SetSuccess(true)
	return m.reply(commandID, reply)
}

func (m *manager) ReplyFailure(commandID string, reply msg.Reply) error {
	reply.SetSuccess(false)
	return m.reply(commandID, reply)
}

func (m *manager) reply(commandID string, reply msg.Reply) error {
	message := m.sagaRepo.GetProcessedMessageByID(commandID)
	if message == nil {
		return apperror.New(apperror.InternalServerError, "cannot reply to unrecorded command id "+commandID)
	}

	command, err := msg.ValidateCommand(message)
	if err != nil {
		return apperror.WithLog(err, "validate command")
	}

	reply.SetID(m.uuidGenerator.NewUUID())
	reply.SetCommandID(commandID)
	reply.SetSagaID(command.SagaID())
	reply.SetSagaID(command.SagaID())

	err = m.producer.Send(command.ReplyChannel(), reply)
	if err != nil {
		return apperror.WithLog(err, "send reply")
	}

	return nil
}
