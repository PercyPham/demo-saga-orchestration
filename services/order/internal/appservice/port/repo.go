package port

import (
	"services.order/internal/domain"
	"services.shared/saga"
)

type Repo interface {
	Ping() error

	OrderRepo
	saga.Repo
}

type OrderRepo interface {
	// CreateOrder creates order record in database and assign new created ID to order input
	CreateOrder(order *domain.Order) error

	// FindOrderByID finds and returns order in database, return nil if not found
	FindOrderByID(id int64) *domain.Order
	// FindOrders finds and returns orders
	// 	TODO: pagination by param
	FindOrders() ([]*domain.Order, error)
}
