package domain

const (
	PaymentStatusPending = "PENDING"
	PaymentStatusPaid    = "PAID"
)

type Payment struct {
	OrderID   int64  `json:"order_id"`
	Total     int64  `json:"total"`
	Status    string `json:"status"`
	CommandID string `json:"command_id"`
}
