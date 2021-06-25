package saga

import (
	"errors"
	"services.shared/saga/msg"
)

func (m *manager) serveHandlingCommands() {
	commandChan, _, err := m.consumer.Consume(m.commandChannel)
	if err != nil {
		panic("cannot handle saga commands: "+err.Error())
	}
	m.logf("Start handling commands from MessageQueue channel: %s", m.commandChannel)
	for d := range commandChan {
		go m.handleCommandDelivery(d)
	}
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

	m.logf("Handled saga command %s:%s of saga %s", command.Type(), command.ID(), command.SagaID())
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

	processedMessage := tx.Repo().GetProcessedMessageByID(command.ID())
	if processedMessage != nil {
		return nil
	}

	err := handler(command)
	if err != nil {
		return errors.New("handle command: " + err.Error())
	}

	err = tx.Repo().CreateProcessedMessage(command)
	if err != nil {
		return errors.New("record message as processed: " + err.Error())
	}

	return nil
}
