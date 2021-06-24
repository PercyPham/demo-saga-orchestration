package authorize_payment

import (
	"services.payment/internal/appservice/port"
)

func NewAuthorizePaymentService(repo port.Repo) *AuthorizePaymentService {
	return &AuthorizePaymentService{repo}
}

type AuthorizePaymentService struct {
	repo port.Repo
}

type AuthorizePaymentInput struct {
	OrderID int64
	Total   int64
}

func (s *AuthorizePaymentService) AuthorizePayment(input AuthorizePaymentInput) error {
	panic("implement") // TODO: implement
}
