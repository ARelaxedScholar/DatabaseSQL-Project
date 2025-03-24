package mockServices

import (
	"errors"
	"fmt"

	"github.com/sql-project-backend/internal/ports"
)

type MockPaymentService struct{}

func NewPaymentService() ports.PaymentService {
	return &MockPaymentService{}
}

func (s *MockPaymentService) ProcessPayment(stayId int, amount float64, paymentMethod string) error {
	if stayId <= 0 {
		return errors.New("Stay ID cannot be negative.")
	}
	if amount < 0 {
		return errors.New("Amount cannot be negative.")
	}
	if paymentMethod == "" {
		return errors.New("Payment method cannot be empty.")
	}

	// Log the payment processing (this is just a mock, in the future could be enhanced with Stripe or smth).
	fmt.Printf("Mock processing payment for stay %d: amount %.2f via %s.\n", stayId, amount, paymentMethod)
	return nil
}
