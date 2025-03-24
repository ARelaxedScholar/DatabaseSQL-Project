package defaultServices_test

import (
	"strings"
	"testing"
	"time"

	"github.com/sql-project-backend/internal/adapters/domain/defaultServices"
	"github.com/sql-project-backend/internal/adapters/framework/driven/db/mocks"
)

func TestRegisterStay_Success(t *testing.T) {
	mockRepo := mocks.NewMockStayRepository()
	service := defaultServices.NewStayService(mockRepo)
	now := time.Now()

	// We'll use nil for optional pointers.
	var reservationID *int = nil
	checkInEmployeeID := new(int)
	*checkInEmployeeID = 1
	var checkOutEmployeeID *int = nil
	comments := "Test stay registration"

	stay, err := service.RegisterStay(0, 1, 101, reservationID, now.Add(24*time.Hour), now.Add(48*time.Hour), 200.0, "Credit Card", checkInEmployeeID, checkOutEmployeeID, comments)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if stay.ID == 0 {
		t.Errorf("expected non-zero stay ID, got %d", stay.ID)
	}
	if stay.ClientID != 1 {
		t.Errorf("expected clientID 1, got %d", stay.ClientID)
	}
	if stay.RoomID != 101 {
		t.Errorf("expected roomID 101, got %d", stay.RoomID)
	}
	if stay.PaymentMethod != "Credit Card" {
		t.Errorf("expected payment method 'Credit Card', got %s", stay.PaymentMethod)
	}
	if stay.Comments != comments {
		t.Errorf("expected comments %q, got %q", comments, stay.Comments)
	}
}

func TestRegisterStay_InvalidInput(t *testing.T) {
	mockRepo := mocks.NewMockStayRepository()
	service := defaultServices.NewStayService(mockRepo)
	now := time.Now()

	// Provide an invalid date range: arrivalDate is after departureDate.
	_, err := service.RegisterStay(0, 1, 101, nil, now.Add(48*time.Hour), now.Add(24*time.Hour), 200.0, "Credit Card", nil, nil, "Invalid date range")
	if err == nil {
		t.Fatal("expected error due to invalid date range, got nil")
	}
	// Optionally, check if error message contains a hint about dates.
	if !strings.Contains(strings.ToLower(err.Error()), "date") {
		t.Errorf("expected error message to mention date, got: %v", err)
	}
}

func TestUpdateStay_Success(t *testing.T) {
	mockRepo := mocks.NewMockStayRepository()
	service := defaultServices.NewStayService(mockRepo)
	now := time.Now()

	var reservationID *int = nil
	checkInEmployeeID := new(int)
	*checkInEmployeeID = 1
	var checkOutEmployeeID *int = nil
	comments := "Initial stay"

	// First register a stay.
	stay, err := service.RegisterStay(0, 1, 101, reservationID, now.Add(24*time.Hour), now.Add(48*time.Hour), 200.0, "Credit Card", checkInEmployeeID, checkOutEmployeeID, comments)
	if err != nil {
		t.Fatalf("failed to register stay: %v", err)
	}

	// Update with new final price and comments.
	newFinalPrice := 250.0
	newComments := "Updated stay"
	updatedStay, err := service.UpdateStay(stay.ID, 1, 101, reservationID, now.Add(24*time.Hour), now.Add(48*time.Hour), newFinalPrice, "Credit Card", checkInEmployeeID, checkOutEmployeeID, newComments)
	if err != nil {
		t.Fatalf("expected update to succeed, got error: %v", err)
	}
	if updatedStay.FinalPrice != newFinalPrice {
		t.Errorf("expected final price %.2f, got %.2f", newFinalPrice, updatedStay.FinalPrice)
	}
	if updatedStay.Comments != newComments {
		t.Errorf("expected comments %q, got %q", newComments, updatedStay.Comments)
	}
}

func TestUpdateStay_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockStayRepository()
	service := defaultServices.NewStayService(mockRepo)
	now := time.Now()

	// Attempt to update a stay that doesn't exist.
	_, err := service.UpdateStay(999, 1, 101, nil, now.Add(24*time.Hour), now.Add(48*time.Hour), 200.0, "Credit Card", nil, nil, "Test")
	if err == nil {
		t.Fatal("expected error for non-existent stay, got nil")
	}
}

func TestEndStay_Success(t *testing.T) {
	mockRepo := mocks.NewMockStayRepository()
	service := defaultServices.NewStayService(mockRepo)
	now := time.Now()

	var reservationID *int = nil
	checkInEmployeeID := new(int)
	*checkInEmployeeID = 1
	var checkOutEmployeeID *int = nil
	comments := "Test end stay"

	// Create a stay.
	stay, err := service.RegisterStay(0, 1, 101, reservationID, now.Add(24*time.Hour), now.Add(48*time.Hour), 200.0, "Credit Card", checkInEmployeeID, checkOutEmployeeID, comments)
	if err != nil {
		t.Fatalf("failed to register stay: %v", err)
	}

	// End the stay.
	err = service.EndStay(stay.ID)
	if err != nil {
		t.Fatalf("expected end stay to succeed, got error: %v", err)
	}

	// Verify deletion: FindByID should return error or nil.
	endedStay, err := mockRepo.FindByID(stay.ID)
	if err == nil && endedStay != nil {
		t.Errorf("expected stay to be deleted, but found one")
	}
}

func TestEndStay_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockStayRepository()
	service := defaultServices.NewStayService(mockRepo)

	err := service.EndStay(999)
	if err == nil {
		t.Fatal("expected error for non-existent stay, got nil")
	}
}
