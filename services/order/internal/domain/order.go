package domain

import (
	"strconv"

	"services.shared/apperror"
)

type Order struct {
	ID        int64            `json:"id,omitempty"`
	State     string           `json:"state,omitempty"`
	Vendor    string           `json:"vendor"`
	Location  string           `json:"location"`
	LineItems []*OrderLineItem `json:"line_items"`
	Total     int64            `json:"total"`
}

const (
	OrderStateInit = "init"
)

type OrderLineItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
	Note     string `json:"note"`
}

func (o *Order) validate() error {
	if o.Vendor == "" {
		return apperror.New(apperror.BadRequest, "vendor must not be empty")
	}
	if o.Location == "" {
		return apperror.New(apperror.BadRequest, "location must not be empty")
	}
	if len(o.LineItems) == 0 {
		return apperror.New(apperror.BadRequest, "items must not be empty")
	}
	m := map[string]bool{}
	for idx, item := range o.LineItems {
		if item.ID == "" {
			return apperror.New(apperror.BadRequest, "empty item id at index "+strconv.Itoa(idx))
		}
		if item.Quantity < 1 {
			return apperror.New(apperror.BadRequest, "item quantity must be greater than zero")
		}
		if m[item.ID] {
			return apperror.New(apperror.BadRequest, "duplicate items "+item.ID)
		}
		m[item.ID] = true
	}
	return nil
}

func NewOrder(vendor, location string, items ...*OrderLineItem) (*Order, error) {
	order := &Order{
		State:     OrderStateInit,
		Vendor:    vendor,
		Location:  location,
		LineItems: items,
	}
	if err := order.validate(); err != nil {
		return nil, apperror.WithLog(err, "validate order")
	}
	return order, nil
}
