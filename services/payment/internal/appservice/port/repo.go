package port

import (
	"services.payment/internal/domain"
	"services.shared/saga"
)

type Repo interface {
	Ping() error

	saga.Repo
	PaymentRepo
}

type PaymentRepo interface {
	CreatePayment(payment *domain.Payment) error
	FindPaymentByOrderID(orderID int64) *domain.Payment
	UpdatePayment(payment *domain.Payment) error
}
