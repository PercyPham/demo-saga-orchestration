package postgresql

import (
	"encoding/json"
	"services.order/internal/domain"
	"services.shared/apperror"
)


type Order struct {
	ID        int64  `json:"id"`
	State     string `json:"state"`
	Vendor    string `json:"vendor"`
	Location  string `json:"location"`
	LineItems []byte `json:"line_items" sql:"type:json"`
}

func convertOrderToGorm(o *domain.Order) (*Order, error) {
	lineItems, err := json.Marshal(o.LineItems)
	if err != nil {
		return nil, apperror.WithLog(err, "marshal order line items")
	}
	return &Order{
		State:     o.State,
		Vendor:    o.Vendor,
		Location:  o.Location,
		LineItems: lineItems,
	}, nil
}

func (o *Order) toDomainOrder() *domain.Order {
	lineItems := make([]*domain.OrderLineItem, 0)
	_ = json.Unmarshal(o.LineItems, &lineItems)
	return &domain.Order{
		ID:        o.ID,
		State:     o.State,
		Vendor:    o.Vendor,
		Location:  o.Location,
		LineItems: lineItems,
	}
}

func (r *repoImpl) CreateOrder(order *domain.Order) error {
	oGorm, err := convertOrderToGorm(order)
	if err != nil {
		return apperror.WithLog(err, "convert order to gorm order")
	}
	result := r.db.Create(oGorm)
	if result.Error != nil {
		return apperror.WithLog(result.Error, "create order in db using gorm")
	}
	order.ID = oGorm.ID
	return nil
}
