package defaultServices_test

import (
	"strings"
	"testing"

	"github.com/sql-project-backend/internal/adapters/domain/defaultServices"
	"github.com/sql-project-backend/internal/adapters/framework/driven/db/mocks"
)

func TestCreateHotelChain_Success(t *testing.T) {
	mockRepo := mocks.NewMockHotelChainRepository()
	service := defaultServices.NewHotelChainService(mockRepo)

	chain, err := service.CreateHotelChain(0, 5, "Unique Chain", "Central Address", "info@chain.com", "555-0100")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if chain.ID == 0 {
		t.Errorf("expected a non-zero hotel chain ID, got %d", chain.ID)
	}
	if chain.Name != "Unique Chain" {
		t.Errorf("expected chain name 'Unique Chain', got %s", chain.Name)
	}
}

func TestCreateHotelChain_Duplicate(t *testing.T) {
	// Assume our mock repository enforces uniqueness by chain name.
	mockRepo := mocks.NewMockHotelChainRepository()
	service := defaultServices.NewHotelChainService(mockRepo)

	// Create the first hotel chain.
	_, err := service.CreateHotelChain(0, 5, "Duplicate Chain", "Central Address", "info@chain.com", "555-0100")
	if err != nil {
		t.Fatalf("expected no error on first creation, got: %v", err)
	}

	// Attempt to create a second chain with the same name.
	_, err = service.CreateHotelChain(0, 6, "Duplicate Chain", "Another Address", "contact@chain.com", "555-0200")
	if err == nil {
		t.Fatal("expected error for duplicate hotel chain, got nil")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		t.Errorf("expected error to mention duplicate, got: %v", err)
	}
}

func TestUpdateHotelChain_Success(t *testing.T) {
	mockRepo := mocks.NewMockHotelChainRepository()
	service := defaultServices.NewHotelChainService(mockRepo)

	// Create a hotel chain first.
	origChain, err := service.CreateHotelChain(0, 5, "Original Chain", "Central Address", "info@orig.com", "555-0100")
	if err != nil {
		t.Fatalf("failed to create hotel chain: %v", err)
	}

	// Update the hotel chain with new values.
	updatedChain, err := service.UpdateHotelChain(origChain.ID, 7, "Updated Chain", "New Central Address", "contact@upd.com", "555-0200")
	if err != nil {
		t.Fatalf("expected update to succeed, got error: %v", err)
	}
	if updatedChain.Name != "Updated Chain" {
		t.Errorf("expected chain name 'Updated Chain', got %s", updatedChain.Name)
	}
	if updatedChain.NumberOfHotel != 7 {
		t.Errorf("expected number of hotels to be 7, got %d", updatedChain.NumberOfHotel)
	}
}

func TestUpdateHotelChain_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockHotelChainRepository()
	service := defaultServices.NewHotelChainService(mockRepo)

	_, err := service.UpdateHotelChain(999, 5, "Nonexistent Chain", "Central Address", "info@none.com", "555-0100")
	if err == nil {
		t.Fatal("expected error for non-existent hotel chain, got nil")
	}
}

func TestDeleteHotelChain_Success(t *testing.T) {
	mockRepo := mocks.NewMockHotelChainRepository()
	service := defaultServices.NewHotelChainService(mockRepo)

	// Create a hotel chain.
	chain, err := service.CreateHotelChain(0, 5, "Chain to Delete", "Central Address", "delete@chain.com", "555-0100")
	if err != nil {
		t.Fatalf("failed to create hotel chain: %v", err)
	}

	// Delete the hotel chain.
	err = service.DeleteHotelChain(chain.ID)
	if err != nil {
		t.Fatalf("expected delete to succeed, got error: %v", err)
	}

	// Verify deletion.
	deletedChain, err := mockRepo.FindByID(chain.ID)
	if err == nil && deletedChain != nil {
		t.Errorf("expected hotel chain to be deleted, but it was found")
	}
}

func TestDeleteHotelChain_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockHotelChainRepository()
	service := defaultServices.NewHotelChainService(mockRepo)

	err := service.DeleteHotelChain(999)
	if err == nil {
		t.Fatal("expected error for non-existent hotel chain, got nil")
	}
}
