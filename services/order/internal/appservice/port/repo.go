package port

import (
	"services.order/internal/domain"
)

type Repo interface {
	Ping() error

	// CreateOrder creates order record in database and assign new created ID to order input
	CreateOrder(order *domain.Order) error
}
