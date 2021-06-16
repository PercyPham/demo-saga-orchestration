package create_order

import (
	"github.com/percypham/saga-go"
	"services.order/internal/appservice/port"
	"services.order/internal/domain"
	"services.shared/apperror"
)

func NewCreateOrderService(r port.Repo, sagaManager saga.Manager) *CreateOrderService {
	return &CreateOrderService{r, sagaManager}
}

type CreateOrderService struct {
	repo        port.Repo
	sagaManager saga.Manager
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

	createOrderSaga, err := newCreateOrderSaga(order)
	if err != nil {
		return nil, apperror.WithLog(err, "create CreateOrderSaga instance")
	}

	err = s.sagaManager.ExecuteFirstStep(*createOrderSaga)
	if err != nil {
		return nil, apperror.WithLog(err, "execute first step in CreateOrderSaga")
	}

	return order, nil
}
