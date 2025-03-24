package defaultEmployeeUseCases

import (
	"errors"
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

// EmployeeCheckoutUseCase defines the interface for employee checkout.
type EmployeeCheckoutUseCase interface {
	Checkout(input dto.CheckoutInput) (dto.CheckoutOutput, error)
}

// DefaultEmployeeCheckoutUseCase is the default implementation of EmployeeCheckoutUseCase.
type DefaultEmployeeCheckoutUseCase struct {
	stayService    ports.StayService    // Service to update or end a stay
	paymentService ports.PaymentService // Service to process payment
}

// NewEmployeeCheckoutUseCase constructs a new instance of DefaultEmployeeCheckoutUseCase.
func NewEmployeeCheckoutUseCase(stayService ports.StayService, paymentService ports.PaymentService) EmployeeCheckoutUseCase {
	return &DefaultEmployeeCheckoutUseCase{
		stayService:    stayService,
		paymentService: paymentService,
	}
}

// Checkout processes the payment for a stay and finalizes the checkout process.
func (uc *DefaultEmployeeCheckoutUseCase) Checkout(input dto.CheckoutInput) (dto.CheckoutOutput, error) {
	// Validate inputs.
	if input.StayID <= 0 {
		return dto.CheckoutOutput{}, errors.New("invalid stay ID")
	}
	if input.FinalPrice < 0 {
		return dto.CheckoutOutput{}, errors.New("final price cannot be negative")
	}
	if input.PaymentMethod == "" {
		return dto.CheckoutOutput{}, errors.New("payment method cannot be empty")
	}

	// Process payment using the PaymentService.
	if err := uc.paymentService.ProcessPayment(input.StayID, input.FinalPrice, input.PaymentMethod); err != nil {
		return dto.CheckoutOutput{}, err
	}

	// Finalize the stay checkout. Here we call EndStay to "end" the stay.
	if err := uc.stayService.EndStay(input.StayID); err != nil {
		return dto.CheckoutOutput{}, err
	}

	return dto.CheckoutOutput{
		StayID:  input.StayID,
		Message: "Checkout successful",
	}, nil
}
