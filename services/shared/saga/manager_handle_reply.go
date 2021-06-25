package saga

import (
	"errors"
	"services.shared/saga/msg"
)

func (m *manager) serveHandlingReplies() {
	replyChan, _, err := m.consumer.Consume(m.replyChannel)
	if err != nil {
		panic("cannot handle saga replies: " + err.Error())
	}
	m.logf("Start receiving replies from MessageQueue channel: %s", m.replyChannel)
	go func() {
		for d := range replyChan {
			go m.handleReplyDelivery(d)
		}
	}()
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
		m.logf("Error: cannot handle reply message: %s, reason: %v\nreply message: %s", reply.Type(), err, genJsonMessage(d.Message))
		d.Nack()
		tx.RollbackTransaction()
		return
	}

	m.logf("Handled saga reply %s:%s of saga %s", reply.Type(), reply.ID(), reply.SagaID())
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

	processedMessage := tx.Repo().GetProcessedMessageByID(reply.ID())
	if processedMessage != nil {
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

	err = tx.Repo().CreateProcessedMessage(reply)
	if err != nil {
		return err
	}

	return nil
}
