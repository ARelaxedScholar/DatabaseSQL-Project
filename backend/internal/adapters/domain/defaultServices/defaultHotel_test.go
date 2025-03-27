package defaultServices_test

import (
	"strings"
	"testing"

	"github.com/sql-project-backend/internal/adapters/domain/defaultServices"
	"github.com/sql-project-backend/internal/adapters/framework/driven/db/mocks"
)

func TestAddHotel_Success(t *testing.T) {
	mockRepo := mocks.NewMockHotelRepository()
	service := defaultServices.NewHotelService(mockRepo)

	// Call AddHotel with valid input.
	hotel, err := service.AddHotel(0, 1, 4, 100, "The Unique Hotel", "123 Main St", "Foobar City", "info@uniquehotel.com", "555-0101")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if hotel.ID == 0 {
		t.Errorf("expected a non-zero hotel ID, got %d", hotel.ID)
	}
	if hotel.Name != "The Unique Hotel" {
		t.Errorf("expected hotel name 'The Unique Hotel', got %s", hotel.Name)
	}
}

func TestAddHotel_Duplicate(t *testing.T) {
	mockRepo := mocks.NewMockHotelRepository()
	service := defaultServices.NewHotelService(mockRepo)

	// First hotel should be saved without error.
	_, err := service.AddHotel(0, 1, 4, 100, "Duplicate Hotel", "123 Main St", "Foobar City", "info@dup.com", "555-0101")
	if err != nil {
		t.Fatalf("expected no error on first add, got: %v", err)
	}

	// Second attempt with the same hotel name should fail.
	_, err = service.AddHotel(0, 1, 5, 150, "Duplicate Hotel", "456 Other Ave", "Foobar City", "contact@dup.com", "555-0202")
	if err == nil {
		t.Fatal("expected error for duplicate hotel name, got nil")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		t.Errorf("expected error message to mention duplicate, got: %v", err)
	}
}

func TestUpdateHotel_Success(t *testing.T) {
	mockRepo := mocks.NewMockHotelRepository()
	service := defaultServices.NewHotelService(mockRepo)

	// First add a hotel.
	origHotel, err := service.AddHotel(0, 1, 4, 100, "Original Hotel", "123 Main St", "Foobar City", "info@orig.com", "555-0101")
	if err != nil {
		t.Fatalf("failed to add hotel: %v", err)
	}

	// Update the hotel with new values.
	updatedHotel, err := service.UpdateHotel(origHotel.ID, 2, 5, 200, "Updated Hotel", "456 New Ave", "Foobar City", "contact@upd.com", "555-0202")
	if err != nil {
		t.Fatalf("expected update to succeed, got error: %v", err)
	}
	if updatedHotel.Name != "Updated Hotel" {
		t.Errorf("expected hotel name to be 'Updated Hotel', got %s", updatedHotel.Name)
	}
	if updatedHotel.ChainID != 2 || updatedHotel.Rating != 5 || updatedHotel.NumberOfRooms != 200 {
		t.Errorf("update did not modify numeric fields correctly")
	}
}

func TestUpdateHotel_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockHotelRepository()
	service := defaultServices.NewHotelService(mockRepo)

	// Attempt to update a hotel that doesn't exist.
	_, err := service.UpdateHotel(999, 1, 4, 100, "Nonexistent Hotel", "123 Main St", "Foobar City", "info@none.com", "555-0101")
	if err == nil {
		t.Fatal("expected error for non-existent hotel, got nil")
	}
}

func TestDeleteHotel_Success(t *testing.T) {
	mockRepo := mocks.NewMockHotelRepository()
	service := defaultServices.NewHotelService(mockRepo)

	// Add a hotel.
	hotel, err := service.AddHotel(0, 1, 4, 100, "Hotel to Delete", "123 Main St", "Foobar City", "delete@hotel.com", "555-0101")
	if err != nil {
		t.Fatalf("failed to add hotel: %v", err)
	}

	// Delete the hotel.
	err = service.DeleteHotel(hotel.ID)
	if err != nil {
		t.Fatalf("expected delete to succeed, got error: %v", err)
	}

	// Verify that the hotel is removed.
	deletedHotel, err := mockRepo.FindByID(hotel.ID)
	if err == nil && deletedHotel != nil {
		t.Errorf("expected hotel to be deleted, but it was found")
	}
}

func TestDeleteHotel_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockHotelRepository()
	service := defaultServices.NewHotelService(mockRepo)

	// Attempt to delete a hotel that doesn't exist.
	err := service.DeleteHotel(999)
	if err == nil {
		t.Fatal("expected error for non-existent hotel, got nil")
	}
}
