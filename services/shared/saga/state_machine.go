package saga

import (
	"errors"
	"services.shared/saga/msg"
)

type StateMachine interface {
	SagaType() string
	ExecuteFirstStep(currentSaga Saga) (newSaga Saga, nextCommand msg.Command, err error)
	Process(currentSaga Saga, reply msg.Reply) (newSaga Saga, nextCommand msg.Command, err error)
}

type stateMachine struct {
	sagaType string
	steps    []*stateMachineStep
}

type stateMachineStep struct {
	stepType int
	process  interface{}
	meta     stateMachineStepMeta
}

const (
	stepTypeCompensation = iota
	stepTypeInvokeParticipant
	stepTypeOnReply
)

type compensateProcessor func(sagaData []byte) (nextCommand msg.Command, err error)
type invokeParticipantProcessor func(sagaData []byte) (nextCommand msg.Command, err error)
type onReplyProcessor func(sagaData []byte, reply msg.Reply) error

type stateMachineStepMeta struct {
	replyType string
}

func (m *stateMachine) SagaType() string {
	return m.sagaType
}

func (m *stateMachine) ExecuteFirstStep(currentSaga Saga) (newSaga Saga, nextCommand msg.Command, err error) {
	if currentSaga.EndState {
		return currentSaga, nil, nil
	}
	if currentSaga.CurrentStep != 0 {
		return Saga{}, nil, errors.New("saga " + currentSaga.ID + " executed first step before")
	}

	process, stepIdx, endState := m.getForwardProcessor(0)
	if endState {
		newSaga = currentSaga
		newSaga.CurrentStep = stepIdx
		newSaga.EndState = true
		return newSaga, nil, nil
	}

	nextCommand, err = process(currentSaga.Data)
	if err != nil {
		return Saga{}, nil, errors.New("cannot execute first step of saga " + currentSaga.ID + ": " + err.Error())
	}
	newSaga = currentSaga
	newSaga.CurrentStep = stepIdx
	return newSaga, nextCommand, nil
}

func (m *stateMachine) Process(currentSaga Saga, reply msg.Reply) (newSaga Saga, nextCommand msg.Command, err error) {
	if currentSaga.EndState {
		return currentSaga, nil, nil
	}
	if currentSaga.Compensating {
		return m.processBackward(currentSaga, reply)
	}
	return m.processForward(currentSaga, reply)
}

func (m *stateMachine) processBackward(currentSaga Saga, reply msg.Reply) (newSaga Saga, nextCommand msg.Command, err error) {
	if !reply.Success() {
		return Saga{}, nil, errors.New("unexpected failure when compensating: " + convertToJson(reply))
	}
	process, stepIdx, ensState := m.getBackwardProcessor(currentSaga.CurrentStep)
	if ensState {
		newSaga = currentSaga
		newSaga.CurrentStep = stepIdx
		newSaga.EndState = true
		return newSaga, nil, nil
	}
	nextCommand, err = process(currentSaga.Data)
	if err != nil {
		return Saga{}, nil, errors.New("process saga data: " + err.Error())
	}
	newSaga = currentSaga
	newSaga.CurrentStep = stepIdx
	return newSaga, nextCommand, nil
}

func (m *stateMachine) getBackwardProcessor(currentSagaStep int) (process compensateProcessor, stepIdx int, endState bool) {
	for i := currentSagaStep - 1; i >= 0; i-- {
		step := m.steps[i]
		if step.stepType == stepTypeCompensation {
			s, _ := step.process.(compensateProcessor)
			return s, i, false
		}
	}
	return nil, -1, true
}

func (m *stateMachine) processForward(currentSaga Saga, reply msg.Reply) (newSaga Saga, nextCommand msg.Command, err error) {
	m.processReplyIfNecessary(currentSaga, reply)
	if !reply.Success() {
		newSaga = currentSaga
		newSaga.Compensating = true
		process, stepIdx, endState := m.getBackwardProcessor(newSaga.CurrentStep)
		if endState {
			newSaga.EndState = true
			newSaga.CurrentStep = stepIdx
			return newSaga, nil, nil
		}
		nextCommand, err = process(newSaga.Data)
		if err != nil {
			return Saga{}, nil, errors.New("failed to process backward in response of a failure reply: " + convertToJson(reply))
		}
		return newSaga, nextCommand, nil
	}

	process, stepIdx, endState := m.getForwardProcessor(currentSaga.CurrentStep)
	if endState {
		newSaga = currentSaga
		newSaga.EndState = true
		newSaga.CurrentStep = stepIdx
		return newSaga, nil, nil
	}
	nextCommand, err = process(currentSaga.Data)
	if err != nil {
		return Saga{}, nil, errors.New("failed to process forward in response of success reply: " + convertToJson(reply))
	}

	newSaga = currentSaga
	newSaga.CurrentStep = stepIdx

	return newSaga, nextCommand, nil
}

func (m *stateMachine) getForwardProcessor(currentSagaStep int) (process invokeParticipantProcessor, stepIdx int, endState bool) {
	for i := currentSagaStep + 1; i < len(m.steps); i++ {
		step := m.steps[i]
		if step.stepType == stepTypeInvokeParticipant {
			s, _ := step.process.(invokeParticipantProcessor)
			return s, i, false
		}
	}
	return nil, len(m.steps), true
}

func (m *stateMachine) processReplyIfNecessary(currentSaga Saga, reply msg.Reply) error {
	if currentSaga.CurrentStep >= len(m.steps) {
		return nil
	}
	for i := currentSaga.CurrentStep + 1; i < len(m.steps); i++ {
		step := m.steps[i]
		if step.stepType != stepTypeOnReply {
			return nil
		}
		if reply.Type() == step.meta.replyType {
			process, _ := step.process.(onReplyProcessor)
			return process(currentSaga.Data, reply)
		}
	}
	return nil
}

func convertToJson(m msg.Message) string {
	j, _ := msg.Marshal(m)
	return string(j)
}
