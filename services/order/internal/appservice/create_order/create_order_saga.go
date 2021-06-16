package create_order

import (
	"encoding/json"
	"github.com/percypham/saga-go"
	"github.com/percypham/saga-go/msg"
	"services.order/internal/appservice/proxy/kitchenproxy"
	"services.order/internal/appservice/proxy/orderproxy"
	"services.order/internal/appservice/proxy/paymentproxy"
	"services.order/internal/domain"
	"services.shared/apperror"
)

const (
	SagaTypeCreateOrder = "CreateOrder"
)

func newCreateOrderSaga(order *domain.Order) (*saga.Saga, error) {
	jsonOrder, err := json.Marshal(order)
	if err != nil {
		return nil, apperror.WithLog(err, "marshal order to json")
	}
	return saga.NewSaga(SagaTypeCreateOrder, jsonOrder), nil
}

func NewCreateOrderStateMachine() saga.StateMachine {
	m := new(createOrderStateMachine)
	sm, err := saga.StateMachineBuilder().
		ForSagaType(SagaTypeCreateOrder).
		WithCompensation(m.rejectOrder).
		InvokeParticipant(m.createTicket).
		WithCompensation(m.rejectTicket).
		InvokeParticipant(m.authorizePayment).
		InvokeParticipant(m.approveTicket).
		InvokeParticipant(m.approveOrder).
		Build()
	if err != nil {
		panic(err)
	}
	return sm
}

type createOrderStateMachine struct{}

func (m *createOrderStateMachine) rejectOrder(sagaData []byte) (msg.Command, error) {
	order, err := m.deserialize(sagaData)
	if err != nil {
		return nil, apperror.WithLog(err, "deserialize saga data")
	}
	rejectOrderCommand := orderproxy.GenRejectOrderCommand(order.ID)
	return rejectOrderCommand, nil
}

func (m *createOrderStateMachine) createTicket(sagaData []byte) (msg.Command, error) {
	order, err := m.deserialize(sagaData)
	if err != nil {
		return nil, apperror.WithLog(err, "deserialize saga data")
	}
	createTicketCommand, err := kitchenproxy.GenCreateTicketCommand(order)
	if err != nil {
		return nil, apperror.WithLog(err, "generate Kitchen's CreateTicket command")
	}
	return createTicketCommand, nil
}

func (m *createOrderStateMachine) rejectTicket(sagaData []byte) (msg.Command, error) {
	order, err := m.deserialize(sagaData)
	if err != nil {
		return nil, apperror.WithLog(err, "deserialize saga data")
	}
	rejectTicketCommand := kitchenproxy.GenRejectTicketCommand(order.ID)
	return rejectTicketCommand, nil
}

func (m *createOrderStateMachine) authorizePayment(sagaData []byte) (msg.Command, error) {
	order, err := m.deserialize(sagaData)
	if err != nil {
		return nil, apperror.WithLog(err, "deserialize saga data")
	}
	authorizePaymentCommand, err := paymentproxy.GenAuthorizePaymentCommand(order)
	if err != nil {
		return nil, apperror.WithLog(err, "generate AuthorizePayment command")
	}
	return authorizePaymentCommand, nil
}

func (m *createOrderStateMachine) approveTicket(sagaData []byte) (msg.Command, error) {
	order, err := m.deserialize(sagaData)
	if err != nil {
		return nil, apperror.WithLog(err, "deserialize saga data")
	}
	approveTicketCommand := kitchenproxy.GenApproveTicketCommand(order.ID)
	return approveTicketCommand, nil
}

func (m *createOrderStateMachine) approveOrder(sagaData []byte) (msg.Command, error) {
	order, err := m.deserialize(sagaData)
	if err != nil {
		return nil, apperror.WithLog(err, "deserialize saga data")
	}
	approveOrderCommand := orderproxy.GenApproveOrderCommand(order.ID)
	return approveOrderCommand, nil
}

func (m *createOrderStateMachine) deserialize(sagaData []byte) (order *domain.Order, err error) {
	err = json.Unmarshal(sagaData, order)
	if err != nil {
		return nil, apperror.WithLog(err, "unmarshal sagaData to Order")
	}
	return order, nil
}
