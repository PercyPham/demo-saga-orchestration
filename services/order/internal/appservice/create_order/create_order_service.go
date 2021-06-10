package create_order

import (
	"github.com/percypham/saga-go"
	"services.order/internal/appservice/port"
	"services.order/internal/domain"
	"services.shared/apperror"

	"github.com/percypham/saga-go/msg"
)

func NewCreateOrderService(r port.Repo, p saga.PSPublisher) *CreateOrderService {
	return &CreateOrderService{r, p}
}

type CreateOrderService struct {
	repo      port.Repo
	publisher saga.PSPublisher
}

type CreateOrderInput struct {
	Vendor    string                  `json:"vendor"`
	Location  string                  `json:"location"`
	LineItems []*domain.OrderLineItem `json:"line_items"`
}

func (s *CreateOrderService) CreateOrder(input CreateOrderInput) (*domain.Order, error) {
	order, err := domain.NewOrder(input.Vendor, input.Location, input.LineItems...)
	if err != nil {
		return nil, apperror.WithLog(err, "create order from input")
	}

	err = s.repo.CreateOrder(order)
	if err != nil {
		return nil, apperror.WithLog(err, "create order in database")
	}

	//orderCreatedEvent := newOrderCreatedEvent(order.ID)
	//s.publisher.Publish(orderCreatedEvent.Topic(), orderCreatedEvent)

	return order, nil
}

func newOrderCreatedEvent(orderID int64) msg.Event {
	// TODO
	return nil
}
