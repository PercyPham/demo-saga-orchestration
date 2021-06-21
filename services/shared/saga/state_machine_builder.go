package saga

import (
	"errors"
	"strconv"
)

func StateMachineBuilder() *builder {
	steps := make([]*stateMachineStep, 0)
	return &builder{"", steps}
}

type builder struct {
	sagaType     string
	machineSteps []*stateMachineStep
}

func (b *builder) For(sagaType string) *builder {
	b.sagaType = sagaType
	return b
}

func (b *builder) WithCompensation(f compensateProcessor) *builder {
	b.machineSteps = append(b.machineSteps, &stateMachineStep{
		stepType: stepTypeCompensation,
		process:  f,
	})
	return b
}

func (b *builder) InvokeParticipant(f invokeParticipantProcessor) *builder {
	b.machineSteps = append(b.machineSteps, &stateMachineStep{
		stepType: stepTypeInvokeParticipant,
		process:  f,
	})
	return b
}

func (b *builder) OnReply(replyType string, f onReplyProcessor) *builder {
	b.machineSteps = append(b.machineSteps, &stateMachineStep{
		stepType: stepTypeOnReply,
		meta:     stateMachineStepMeta{replyType: replyType},
		process:  f,
	})
	return b
}

func (b *builder) Build() (StateMachine, error) {
	if b.sagaType == "" {
		return nil, errors.New("SagaType is not set, use builder.For(sagaType) to set SagaType for StateMachine")
	}
	for i, step := range b.machineSteps {
		if step.stepType == stepTypeOnReply && step.meta.replyType == "" {
			return nil, errors.New("OnReply at [" + strconv.Itoa(i) + "] should not have empty ReplyType")
		}
	}
	return &stateMachine{
		sagaType: b.sagaType,
		steps:    b.machineSteps,
	}, nil
}
