package create_ticket

import (
	"services.kitchen/internal/appservice/port"
	"services.kitchen/internal/domain"
	"services.shared/apperror"
)

func NewCreateTicketService(r port.Repo) *CreateTicketService {
	return &CreateTicketService{r}
}

type CreateTicketService struct {
	repo port.Repo
}

type CreateTicketInput struct {
	OrderID   int64
	CommandID string
	Vendor    string
	LineItems []domain.LineItem
}

func (s *CreateTicketService) CreateTicket(input CreateTicketInput) error {
	if input.CommandID == "" {
		return apperror.New(apperror.InvalidCommand, "command ID must not be empty")
	}
	if len(input.LineItems) == 0 {
		return apperror.New(apperror.InvalidCommand, "line_items must not be empty")
	}

	if ticket := s.repo.FindTicketByOrderID(input.OrderID); ticket != nil {
		return nil
	}

	ticket := domain.Ticket{
		OrderID:   input.OrderID,
		CommandID: input.CommandID,
		Vendor:    input.Vendor,
		Status:    domain.TicketStatusPending,
		LineItems: input.LineItems,
	}

	err := s.repo.CreateTicket(&ticket)
	if err != nil {
		return apperror.WithLog(err, "create ticket in db")
	}

	return nil
}
