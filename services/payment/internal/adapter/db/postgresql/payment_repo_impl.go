package postgresql

import (
	"services.payment/internal/domain"
	"services.shared/apperror"
)

type Payment struct {
	OrderID   int64  `json:"order_id" gorm:"primaryKey"`
	Total     int64  `json:"total"`
	Status    string `json:"status"`
	CommandID string `json:"command_id"`
}

func (p *Payment) toDomainPayment() *domain.Payment {
	return &domain.Payment{
		OrderID:   p.OrderID,
		Total:     p.Total,
		Status:    p.Status,
		CommandID: p.CommandID,
	}
}

func convertPaymentToGorm(p *domain.Payment) *Payment {
	return &Payment{
		OrderID:   p.OrderID,
		Total:     p.Total,
		Status:    p.Status,
		CommandID: p.CommandID,
	}
}

func (r *repoImpl) CreatePayment(ticket *domain.Payment) error {
	gormPayment := convertPaymentToGorm(ticket)
	result := r.db.Create(gormPayment)
	if result.Error != nil {
		return apperror.Wrap(result.Error, "create payment in db using gorm")
	}
	return nil
}

func (r *repoImpl) FindPaymentByOrderID(orderID int64) *domain.Payment {
	gormPayment := new(Payment)
	result := r.db.Where("order_id = ?", orderID).First(gormPayment)
	if result.Error != nil {
		return nil
	}
	return gormPayment.toDomainPayment()
}

func (r *repoImpl) UpdatePayment(payment *domain.Payment) error {
	gormPayment := convertPaymentToGorm(payment)
	result := r.db.Updates(gormPayment)
	if result.Error != nil {
		return apperror.Wrap(result.Error, "update payment using gorm")
	}
	return nil
}
