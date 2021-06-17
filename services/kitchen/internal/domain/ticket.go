package domain

const (
	TicketStatusPending   = "PENDING"
	TicketStatusRejected  = "REJECTED"
	TicketStatusApproved  = "APPROVED"
	TicketStatusFulfilled = "FULFILLED"
)

type Ticket struct {
	OrderID   int64      `json:"order_id" gorm:"primaryKey"`
	Vendor    string     `json:"vendor"`
	SagaID    string     `json:"saga_id"`
	CommandID string     `json:"command_id"`
	Status    string     `json:"status"`
	LineItems []LineItem `json:"line_items"`
}

type LineItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
	Note     string `json:"note"`
}
