package domain_test

import (
	"strings"
	"testing"

	"github.com/sql-project-backend/internal/adapters/domain" // Adjust the import path as necessary.
	"github.com/sql-project-backend/internal/ports"
)

func TestProcessPayment_Success(t *testing.T) {
	// Create an instance of the payment service.
	var paymentService ports.PaymentService = domain.NewPaymentService()

	// Call ProcessPayment with valid parameters.
	err := paymentService.ProcessPayment(1, 100.0, "Credit Card")
	if err != nil {
		t.Fatalf("expected success, but got error: %v", err)
	}
}

func TestProcessPayment_InvalidStayID(t *testing.T) {
	var paymentService ports.PaymentService = domain.NewPaymentService()

	// Use a stay ID that is zero (or negative) to trigger the error.
	err := paymentService.ProcessPayment(0, 100.0, "Credit Card")
	if err == nil {
		t.Fatal("expected error for invalid stay ID, got nil")
	}
	if !strings.Contains(err.Error(), "Stay ID") {
		t.Errorf("expected error to mention 'Stay ID', got: %v", err)
	}
}

func TestProcessPayment_NegativeAmount(t *testing.T) {
	var paymentService ports.PaymentService = domain.NewPaymentService()

	// Use a negative amount.
	err := paymentService.ProcessPayment(1, -50.0, "Credit Card")
	if err == nil {
		t.Fatal("expected error for negative amount, got nil")
	}
	if !strings.Contains(err.Error(), "Amount") {
		t.Errorf("expected error to mention 'Amount', got: %v", err)
	}
}

func TestProcessPayment_EmptyPaymentMethod(t *testing.T) {
	var paymentService ports.PaymentService = domain.NewPaymentService()

	// Use an empty payment method.
	err := paymentService.ProcessPayment(1, 100.0, "")
	if err == nil {
		t.Fatal("expected error for empty payment method, got nil")
	}
	if !strings.Contains(err.Error(), "Payment method") {
		t.Errorf("expected error to mention 'Payment method', got: %v", err)
	}
}
