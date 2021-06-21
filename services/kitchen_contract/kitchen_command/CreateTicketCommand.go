package kitchen_command

type CreateTicketCommand struct {
	OrderID   int64           `json:"order_id"`
	Vendor    string          `json:"vendor"`
	LineItems []OrderLineItem `json:"line_items"`
}

type OrderLineItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
	Note     string `json:"note"`
}
