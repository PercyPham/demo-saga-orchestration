package postgresql

import (
	"encoding/json"
	"services.kitchen/internal/domain"
	"services.shared/apperror"
)

type Ticket struct {
	OrderID   int64  `json:"order_id" gorm:"primaryKey"`
	Status    string `json:"status"`
	CommandID string `json:"command_id"`
	Vendor    string `json:"vendor"`
	LineItems []byte `json:"line_items" sql:"type:json"`
}

func (t *Ticket) toDomainTicket() *domain.Ticket {
	lineItems := make([]domain.LineItem, 0)
	_ = json.Unmarshal(t.LineItems, &lineItems)
	return &domain.Ticket{
		OrderID:   t.OrderID,
		Status:    t.Status,
		CommandID: t.CommandID,
		Vendor:    t.Vendor,
		LineItems: lineItems,
	}
}

func convertTicketToGorm(t *domain.Ticket) (*Ticket, error) {
	lineItems, err := json.Marshal(t.LineItems)
	if err != nil {
		return nil, apperror.Wrap(err, "marshal order line items")
	}
	return &Ticket{
		OrderID:   t.OrderID,
		Status:    t.Status,
		CommandID: t.CommandID,
		Vendor:    t.Vendor,
		LineItems: lineItems,
	}, nil
}

func (r *repoImpl) CreateTicket(ticket *domain.Ticket) error {
	gormTicket, err := convertTicketToGorm(ticket)
	if err != nil {
		return apperror.Wrap(err, "convert ticket to gorm ticket")
	}
	result := r.db.Create(gormTicket)
	if result.Error != nil {
		return apperror.Wrap(result.Error, "create ticket in db using gorm")
	}
	return nil
}

func (r *repoImpl) FindTicketByOrderID(orderID int64) *domain.Ticket {
	gormTicket := new(Ticket)
	result := r.db.Where("order_id = ?", orderID).First(gormTicket)
	if result.Error != nil {
		return nil
	}
	return gormTicket.toDomainTicket()
}

func (r *repoImpl) FindTickets() ([]*domain.Ticket, error) {
	gormTickets := make([]*Ticket, 0)
	result := r.db.Find(&gormTickets)
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, "find tickets using gorm")
	}
	ticket := make([]*domain.Ticket, result.RowsAffected)
	for i, gormTicket := range gormTickets {
		ticket[i] = gormTicket.toDomainTicket()
	}
	return ticket, nil
}

func (r *repoImpl) UpdateTicket(ticket *domain.Ticket) error {
	gormTicket, err := convertTicketToGorm(ticket)
	if err != nil {
		return apperror.Wrap(err, "convert domain ticket to gorm ticket")
	}
	result := r.db.Updates(gormTicket)
	if result.Error != nil {
		return apperror.Wrap(result.Error, "update ticket using gorm")
	}
	return nil
}
