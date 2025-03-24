package defaultServices_test

import (
	"strings"
	"testing"
	"time"

	"github.com/sql-project-backend/internal/adapters/domain/defaultServices"
	"github.com/sql-project-backend/internal/adapters/framework/driven/db/mocks"
	"github.com/sql-project-backend/internal/models"
)

// TestCreateReservation_Success verifies that a valid reservation is created.
func TestCreateReservation_Success(t *testing.T) {
	mockRepo := mocks.NewMockReservationRepository()
	service := defaultServices.NewReservationService(mockRepo)

	now := time.Now()
	startDate := now.Add(24 * time.Hour)
	endDate := now.Add(48 * time.Hour)
	reservationDate := now
	totalPrice := 150.0
	status := models.Confirmed // assuming Confirmed is a valid ReservationStatus

	reservation, err := service.CreateReservation(0, 1, 1, 101, startDate, endDate, reservationDate, totalPrice, status)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if reservation.ID == 0 {
		t.Errorf("expected a non-zero reservation ID, got %d", reservation.ID)
	}
	if reservation.ClientID != 1 {
		t.Errorf("expected client id 1, got %d", reservation.ClientID)
	}
	if reservation.TotalPrice != totalPrice {
		t.Errorf("expected total price %.2f, got %.2f", totalPrice, reservation.TotalPrice)
	}
}

// TestCreateReservation_InvalidInput simulates invalid input (e.g. startDate after endDate)
func TestCreateReservation_InvalidInput(t *testing.T) {
	mockRepo := mocks.NewMockReservationRepository()
	service := defaultServices.NewReservationService(mockRepo)

	now := time.Now()
	// Intentionally invalid date range: start date after end date.
	startDate := now.Add(48 * time.Hour)
	endDate := now.Add(24 * time.Hour)
	reservationDate := now
	totalPrice := 150.0
	status := models.Confirmed

	_, err := service.CreateReservation(0, 1, 1, 101, startDate, endDate, reservationDate, totalPrice, status)
	if err == nil {
		t.Fatal("expected error due to invalid date range, got nil")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "date") {
		t.Errorf("expected error message to mention date, got: %v", err)
	}
}

// TestUpdateReservation_Success verifies that updating a reservation works.
func TestUpdateReservation_Success(t *testing.T) {
	mockRepo := mocks.NewMockReservationRepository()
	service := defaultServices.NewReservationService(mockRepo)

	now := time.Now()
	startDate := now.Add(24 * time.Hour)
	endDate := now.Add(48 * time.Hour)
	reservationDate := now
	totalPrice := 150.0
	status := models.Confirmed

	// Create an initial reservation.
	origRes, err := service.CreateReservation(0, 1, 1, 101, startDate, endDate, reservationDate, totalPrice, status)
	if err != nil {
		t.Fatalf("failed to create reservation: %v", err)
	}

	// Update details.
	newStartDate := now.Add(36 * time.Hour)
	newEndDate := now.Add(60 * time.Hour)
	newTotalPrice := 200.0

	updatedRes, err := service.UpdateReservation(origRes.ID, 1, 1, 101, newStartDate, newEndDate, reservationDate, newTotalPrice, status)
	if err != nil {
		t.Fatalf("expected update to succeed, got error: %v", err)
	}
	if !updatedRes.StartDate.Equal(newStartDate) {
		t.Errorf("expected start date %v, got %v", newStartDate, updatedRes.StartDate)
	}
	if updatedRes.TotalPrice != newTotalPrice {
		t.Errorf("expected total price %.2f, got %.2f", newTotalPrice, updatedRes.TotalPrice)
	}
}

// TestUpdateReservation_NotFound tests updating a non-existent reservation.
func TestUpdateReservation_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockReservationRepository()
	service := defaultServices.NewReservationService(mockRepo)

	now := time.Now()
	_, err := service.UpdateReservation(999, 1, 1, 101, now, now.Add(24*time.Hour), now, 150.0, models.Confirmed)
	if err == nil {
		t.Fatal("expected error for non-existent reservation, got nil")
	}
}

// TestCancelReservation_Success tests that a reservation can be cancelled.
func TestCancelReservation_Success(t *testing.T) {
	mockRepo := mocks.NewMockReservationRepository()
	service := defaultServices.NewReservationService(mockRepo)

	now := time.Now()
	startDate := now.Add(24 * time.Hour)
	endDate := now.Add(48 * time.Hour)
	reservationDate := now
	totalPrice := 150.0
	status := models.Confirmed

	// Create a reservation.
	res, err := service.CreateReservation(0, 1, 1, 101, startDate, endDate, reservationDate, totalPrice, status)
	if err != nil {
		t.Fatalf("failed to create reservation: %v", err)
	}

	// Cancel the reservation.
	err = service.CancelReservation(res.ID)
	if err != nil {
		t.Fatalf("expected cancel to succeed, got error: %v", err)
	}

	// Retrieve the reservation and verify that its status is Cancelled.
	updatedRes, err := mockRepo.FindByID(res.ID)
	if err != nil {
		t.Fatalf("failed to retrieve reservation after cancel: %v", err)
	}
	if updatedRes == nil {
		t.Fatalf("expected reservation to exist after cancel")
	}
	if updatedRes.Status != models.Cancelled {
		t.Errorf("expected reservation status to be Cancelled, got %v", updatedRes.Status)
	}
}

// TestCancelReservation_NotFound tests cancelling a non-existent reservation.
func TestCancelReservation_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockReservationRepository()
	service := defaultServices.NewReservationService(mockRepo)

	err := service.CancelReservation(999)
	if err == nil {
		t.Fatal("expected error for non-existent reservation, got nil")
	}
}

// TestGetReservationsByClient verifies that reservations are filtered by client ID.
func TestGetReservationsByClient(t *testing.T) {
	mockRepo := mocks.NewMockReservationRepository()
	service := defaultServices.NewReservationService(mockRepo)

	now := time.Now()
	// Create two reservations for client 1.
	for i := 0; i < 2; i++ {
		_, err := service.CreateReservation(0, 1, 1, 101+i, now.Add(24*time.Hour), now.Add(48*time.Hour), now, 150.0+float64(i*10), models.Confirmed)
		if err != nil {
			t.Fatalf("failed to create reservation %d: %v", i, err)
		}
	}
	// Create one reservation for client 2.
	_, err := service.CreateReservation(0, 2, 1, 201, now.Add(24*time.Hour), now.Add(48*time.Hour), now, 200.0, models.Confirmed)
	if err != nil {
		t.Fatalf("failed to create reservation for client 2: %v", err)
	}

	// Retrieve reservations for client 1.
	reservations, err := service.GetReservationsByClient(1)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	// Note: If your mock repository is not filtering correctly, this test may fail.
	// Expected: 2 reservations for client 1.
	if len(reservations) != 2 {
		t.Errorf("expected 2 reservations for client 1, got %d", len(reservations))
	}

	// Retrieve reservations for client 2.
	reservations, err = service.GetReservationsByClient(2)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(reservations) != 1 {
		t.Errorf("expected 1 reservation for client 2, got %d", len(reservations))
	}

	// Retrieve reservations for a client with no reservations.
	reservations, err = service.GetReservationsByClient(999)
	if err != nil {
		t.Fatalf("expected no error for client with no reservations, got: %v", err)
	}
	if len(reservations) != 0 {
		t.Errorf("expected 0 reservations for client 999, got %d", len(reservations))
	}
}
