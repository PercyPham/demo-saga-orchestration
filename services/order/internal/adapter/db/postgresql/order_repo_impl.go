package postgresql

import (
	"encoding/json"
	"services.order/internal/domain"
	"services.shared/apperror"
)

type Order struct {
	ID        int64  `json:"id"`
	Status    string `json:"status"`
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
		Status:    o.Status,
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
		Status:    o.Status,
		Vendor:    o.Vendor,
		Location:  o.Location,
		LineItems: lineItems,
	}
}

func (r *repoImpl) CreateOrder(order *domain.Order) error {
	gormOrder, err := convertOrderToGorm(order)
	if err != nil {
		return apperror.WithLog(err, "convert order to gorm order")
	}
	result := r.db.Create(gormOrder)
	if result.Error != nil {
		return apperror.WithLog(result.Error, "create order in db using gorm")
	}
	order.ID = gormOrder.ID
	return nil
}

func (r *repoImpl) FindOrderByID(id int64) *domain.Order {
	oGorm := new(Order)
	result := r.db.Where("id = ?", id).First(oGorm)
	if result.Error != nil {
		return nil
	}
	return oGorm.toDomainOrder()
}

func (r *repoImpl) FindOrders() ([]*domain.Order, error) {
	gormOrders := make([]*Order, 0)
	result := r.db.Find(&gormOrders)
	if result.Error != nil {
		return nil, apperror.WithLog(result.Error, "find orders using gorm")
	}
	orders := make([]*domain.Order, result.RowsAffected)
	for i, gormOrder := range gormOrders {
		orders[i] = gormOrder.toDomainOrder()
	}
	return orders, nil
}

func (r *repoImpl) UpdateOrder(order *domain.Order) error {
	gormOrder, err := convertOrderToGorm(order)
	if err != nil {
		return apperror.WithLog(err, "convert domain order to gorm order")
	}
	result := r.db.Where("id = ?", order.ID).Updates(gormOrder)
	if result.Error != nil {
		return apperror.WithLog(result.Error, "update order using gorm")
	}
	return nil
}
