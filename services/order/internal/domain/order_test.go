package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"services.order/internal/domain"
)

func TestNewOrderWithEmptyVendor(t *testing.T) {
	item := &domain.OrderLineItem{"id", 1, "note"}
	_, err := domain.NewOrder("", "location", item)
	if assert.Error(t, err) {
		assert.Equal(t, "vendor must not be empty", err.Error())
	}
}

func TestNewOrderWithEmptyLocation(t *testing.T) {
	item := &domain.OrderLineItem{"id", 1, "note"}
	_, err := domain.NewOrder("vendor", "", item)
	if assert.Error(t, err) {
		assert.Equal(t, "location must not be empty", err.Error())
	}
}

func TestNewOrderWithDuplicateItems(t *testing.T) {
	item1 := &domain.OrderLineItem{ID: "1", Quantity: 1, Note: ""}
	item2 := &domain.OrderLineItem{ID: "1", Quantity: 1, Note: "note2"}
	_, err := domain.NewOrder("vendor", "location", item1, item2)
	if assert.Error(t, err) {
		assert.Equal(t, "duplicate items 1", err.Error())
	}
}

func TestNewOrderWithEmptyItemID(t *testing.T) {
	item1 := &domain.OrderLineItem{ID: "1", Quantity: 1, Note: "note1"}
	item2 := &domain.OrderLineItem{ID: "", Quantity: 1, Note: "note2"}
	item3 := &domain.OrderLineItem{ID: "3", Quantity: 1, Note: "note2"}
	_, err := domain.NewOrder("vendor", "location", item1, item2, item3)
	if assert.Error(t, err) {
		assert.Equal(t, "empty item id at index 1", err.Error())
	}
}

func TestNewOrderWithZeroQuantityItem(t *testing.T) {
	item := &domain.OrderLineItem{"id", 0, "note"}
	_, err := domain.NewOrder("vendor", "location", item)
	if assert.Error(t, err) {
		assert.Equal(t, "item quantity must be greater than zero", err.Error())
	}
}

func TestNewOrderWithNegativeQuantityItem(t *testing.T) {
	item := &domain.OrderLineItem{"id", -1, "note"}
	_, err := domain.NewOrder("vendor", "location", item)
	if assert.Error(t, err) {
		assert.Equal(t, "item quantity must be greater than zero", err.Error())
	}
}

func TestNewOrderWithNoItem(t *testing.T) {
	_, err := domain.NewOrder("vendor", "location")
	if assert.Error(t, err) {
		assert.Equal(t, "items must not be empty", err.Error())
	}
}
func TestNewOrder(t *testing.T) {
	item1 := &domain.OrderLineItem{ID: "1", Quantity: 1, Note: ""}
	item2 := &domain.OrderLineItem{ID: "2", Quantity: 2, Note: "note"}

	order, err := domain.NewOrder("vendor", "location", item1, item2)
	if assert.NoError(t, err) {
		assert.Equal(t, "vendor", order.Vendor)
		assert.Equal(t, "location", order.Location)
		assert.Equal(t, "init", order.State)
	}
}
