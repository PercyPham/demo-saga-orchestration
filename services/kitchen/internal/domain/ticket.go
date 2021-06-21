package domain

const (
	TicketStatusPending   = "PENDING"
	TicketStatusAccepted  = "ACCEPTED"
	TicketStatusRejected  = "REJECTED"
	TicketStatusApproved  = "APPROVED"
	TicketStatusFulfilled = "FULFILLED"
)

type Ticket struct {
	CommandID string     `json:"command_id"`
	OrderID   int64      `json:"order_id"`
	Vendor    string     `json:"vendor"`
	Status    string     `json:"status"`
	LineItems []LineItem `json:"line_items"`
}

type LineItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
	Note     string `json:"note"`
}
