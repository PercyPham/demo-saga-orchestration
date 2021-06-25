package authorize_payment

import (
	"services.payment/internal/appservice/port"
	"services.payment/internal/domain"
	"services.payment_contract/payment_reply"
	"services.shared/apperror"
	"services.shared/saga"
	"time"
)

func NewAuthorizePaymentService(repo port.Repo, sagaManager saga.Manager) *AuthorizePaymentService {
	return &AuthorizePaymentService{repo, sagaManager}
}

type AuthorizePaymentService struct {
	repo        port.Repo
	sagaManager saga.Manager
}

type AuthorizePaymentInput struct {
	OrderID   int64
	Total     int64
	CommandID string
}

func (s *AuthorizePaymentService) AuthorizePayment(input AuthorizePaymentInput) error {
	payment := s.repo.FindPaymentByOrderID(input.OrderID)
	if payment != nil {
		return nil
	}

	payment = &domain.Payment{
		OrderID:   input.OrderID,
		Total:     input.Total,
		Status:    domain.PaymentStatusPending,
		CommandID: input.CommandID,
	}

	err := s.repo.CreatePayment(payment)
	if err != nil {
		return apperror.WithLog(err, "cannot create payment in db")
	}

	go s.mockAuthorizePayment(payment)

	return nil
}

func (s *AuthorizePaymentService) mockAuthorizePayment(payment *domain.Payment) {
	time.Sleep(10 * time.Second)
	payment.Status = domain.PaymentStatusPaid
	err := s.repo.UpdatePayment(payment)
	if err != nil {
		panic("cannot update payment to paid: " + err.Error())
	}
	paymentAuthorizedReply := payment_reply.NewPaymentAuthorizedReply()
	err = s.sagaManager.ReplySuccess(payment.CommandID, paymentAuthorizedReply)
	if err != nil {
		panic("cannot reply that payment authorized: " + err.Error())
	}
}
