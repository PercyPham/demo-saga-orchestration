package postgresql

import (
	"encoding/json"
	"services.kitchen/internal/domain"
	"services.shared/apperror"
)

type Ticket struct {
	OrderID   int64  `json:"order_id" gorm:"primaryKey"`
	Vendor    string `json:"vendor"`
	SagaID    string `json:"saga_id"`
	CommandID string `json:"command_id"`
	Status    string `json:"status"`
	LineItems []byte `json:"line_items" sql:"type:json"`
}

func (t *Ticket) toDomainTicket() *domain.Ticket {
	lineItems := make([]domain.LineItem, 0)
	_ = json.Unmarshal(t.LineItems, &lineItems)
	return &domain.Ticket{
		OrderID:   t.OrderID,
		Vendor:    t.Vendor,
		SagaID:    t.SagaID,
		CommandID: t.CommandID,
		Status:    t.Status,
		LineItems: lineItems,
	}
}

func convertTicketToGorm(t *domain.Ticket) (*Ticket, error) {
	lineItems, err := json.Marshal(t.LineItems)
	if err != nil {
		return nil, apperror.WithLog(err, "marshal order line items")
	}
	return &Ticket{
		OrderID:   t.OrderID,
		Vendor:    t.Vendor,
		SagaID:    t.SagaID,
		CommandID: t.CommandID,
		Status:    t.Status,
		LineItems: lineItems,
	}, nil
}

func (r *repoImpl) CreateTicket(ticket *domain.Ticket) error {
	tGorm, err := convertTicketToGorm(ticket)
	if err != nil {
		return apperror.WithLog(err, "convert ticket to gorm ticket")
	}
	result := r.db.Create(tGorm)
	if result.Error != nil {
		return apperror.WithLog(result.Error, "create ticket in db using gorm")
	}
	return nil
}

func (r *repoImpl) FindTicketByOrderID(orderID int64) *domain.Ticket {
	tGorm := new(Ticket)
	result := r.db.Where("order_id = ?", orderID).First(tGorm)
	if result.Error != nil {
		return nil
	}
	return tGorm.toDomainTicket()
}

func (r *repoImpl) FindTickets() ([]*domain.Ticket, error) {
	gormTickets := make([]*Ticket, 0)
	result := r.db.Find(&gormTickets)
	if result.Error != nil {
		return nil, apperror.WithLog(result.Error, "find tickets using gorm")
	}
	ticket := make([]*domain.Ticket, result.RowsAffected)
	for i, gormTicket := range gormTickets {
		ticket[i] = gormTicket.toDomainTicket()
	}
	return ticket, nil
}
