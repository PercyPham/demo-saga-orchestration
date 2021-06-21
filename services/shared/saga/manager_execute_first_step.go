package saga

import "errors"

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
	nextCmd.SetReplyChannel(m.replyChannel)
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
