package orderdomain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"services.order/internal/orderdomain"
)

func TestNewOrderWithEmptyVendor(t *testing.T) {
	item := &orderdomain.OrderLineItem{"id", 1, "note"}
	_, err := orderdomain.NewOrder("", "location", item)
	if assert.Error(t, err) {
		assert.Equal(t, "vendor must not be empty", err.Error())
	}
}

func TestNewOrderWithEmptyLocation(t *testing.T) {
	item := &orderdomain.OrderLineItem{"id", 1, "note"}
	_, err := orderdomain.NewOrder("vendor", "", item)
	if assert.Error(t, err) {
		assert.Equal(t, "location must not be empty", err.Error())
	}
}

func TestNewOrderWithDuplicateItems(t *testing.T) {
	item1 := &orderdomain.OrderLineItem{ID: "1", Quantity: 1, Note: ""}
	item2 := &orderdomain.OrderLineItem{ID: "1", Quantity: 1, Note: "note2"}
	_, err := orderdomain.NewOrder("vendor", "location", item1, item2)
	if assert.Error(t, err) {
		assert.Equal(t, "duplicate items 1", err.Error())
	}
}

func TestNewOrderWithEmptyItemID(t *testing.T) {
	item1 := &orderdomain.OrderLineItem{ID: "1", Quantity: 1, Note: "note1"}
	item2 := &orderdomain.OrderLineItem{ID: "", Quantity: 1, Note: "note2"}
	item3 := &orderdomain.OrderLineItem{ID: "3", Quantity: 1, Note: "note2"}
	_, err := orderdomain.NewOrder("vendor", "location", item1, item2, item3)
	if assert.Error(t, err) {
		assert.Equal(t, "empty item id at index 1", err.Error())
	}
}

func TestNewOrderWithZeroQuantityItem(t *testing.T) {
	item := &orderdomain.OrderLineItem{"id", 0, "note"}
	_, err := orderdomain.NewOrder("vendor", "location", item)
	if assert.Error(t, err) {
		assert.Equal(t, "item quantity must be greater than zero", err.Error())
	}
}

func TestNewOrderWithNegativeQuantityItem(t *testing.T) {
	item := &orderdomain.OrderLineItem{"id", -1, "note"}
	_, err := orderdomain.NewOrder("vendor", "location", item)
	if assert.Error(t, err) {
		assert.Equal(t, "item quantity must be greater than zero", err.Error())
	}
}

func TestNewOrderWithNoItem(t *testing.T) {
	_, err := orderdomain.NewOrder("vendor", "location")
	if assert.Error(t, err) {
		assert.Equal(t, "items must not be empty", err.Error())
	}
}
func TestNewOrder(t *testing.T) {
	item1 := &orderdomain.OrderLineItem{ID: "1", Quantity: 1, Note: ""}
	item2 := &orderdomain.OrderLineItem{ID: "2", Quantity: 2, Note: "note"}

	order, err := orderdomain.NewOrder("vendor", "location", item1, item2)
	if assert.NoError(t, err) {
		assert.Equal(t, "vendor", order.Vendor)
		assert.Equal(t, "location", order.Location)
		assert.Equal(t, "init", order.State)
	}
}
