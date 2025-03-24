package defaultServices

import (
	"strings"
	"testing"
	"time"

	"github.com/sql-project-backend/internal/adapters/framework/driven/db/mocks"
)

// TestRegisterClient tests the RegisterClient method of DefaultClientService.
func TestRegisterClient(t *testing.T) {
	// Create a mock repository.
	mockRepo := mocks.NewMockClientRepository()
	// Instantiate the client service with the mock repository.
	service := NewClientService(mockRepo)

	// Test data for a valid client.
	now := time.Now()
	inputID := 0 // When registering a new client, id is typically zero.
	sin := "123456789"
	firstName := "John"
	lastName := "Doe"
	address := "123 Main St"
	phone := "555-1234"
	email := "john@example.com"
	joinDate := now

	// Call RegisterClient.
	client, err := service.RegisterClient(inputID, sin, firstName, lastName, address, phone, email, joinDate)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if client.ID == 0 {
		t.Errorf("Expected a non-zero client ID, got %v", client.ID)
	}
	if client.FirstName != firstName {
		t.Errorf("Expected FirstName %s, got %s", firstName, client.FirstName)
	}
}

func TestRegisterClient_DuplicateEmail(t *testing.T) {
	mockRepo := mocks.NewMockClientRepository()
	service := NewClientService(mockRepo)

	now := time.Now()
	email := "testasas@example.com"

	// First registration should succeed.
	_, err := service.RegisterClient(0, "123456789", "John", "Doe", "123 Main St", "555-1234", email, now)
	if err != nil {
		t.Fatalf("expected no error on first registration, got: %v", err)
	}

	// Second registration with the same email should fail.
	_, err = service.RegisterClient(0, "987654321", "Jane", "Smith", "456 Other St", "555-5678", email, now)
	if err == nil {
		t.Fatalf("expected error for duplicate email, got nil")
	}
	// Optionally, check that the error message indicates a duplicate email.
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("expected error to mention duplicate email, got: %v", err)
	}
}

// TestRegisterClient_InvalidInput tests RegisterClient with invalid input.
func TestRegisterClient_InvalidInput(t *testing.T) {
	mockRepo := mocks.NewMockClientRepository()
	service := NewClientService(mockRepo)

	// For example, an empty firstName may be invalid.
	now := time.Now()
	_, err := service.RegisterClient(0, "123456789", "", "Doe", "123 Main St", "555-1234", "john@example.com", now)
	if err == nil {
		t.Fatalf("Expected error for invalid input, got nil")
	}
}

// TestUpdateClient tests UpdateClient for a client that exists.
func TestUpdateClient(t *testing.T) {
	mockRepo := mocks.NewMockClientRepository()
	service := NewClientService(mockRepo)

	// First register a valid client.
	now := time.Now()
	registeredClient, err := service.RegisterClient(0, "123456789", "John", "Doe", "123 Main St", "555-1234", "john@example.com", now)
	if err != nil {
		t.Fatalf("RegisterClient failed: %v", err)
	}

	// Now update the client.
	updatedFirstName := "Jane"
	updatedLastName := "Smith"
	updatedAddress := "456 New Ave"
	updatedPhone := "555-5678"
	updatedEmail := "jane@example.com"

	updatedClient, err := service.UpdateClient(registeredClient.ID, updatedFirstName, updatedLastName, updatedAddress, updatedPhone, updatedEmail)
	if err != nil {
		t.Fatalf("UpdateClient failed: %v", err)
	}
	if updatedClient.FirstName != updatedFirstName || updatedClient.LastName != updatedLastName {
		t.Errorf("UpdateClient did not update names correctly")
	}
	if updatedClient.Address != updatedAddress || updatedClient.Phone != updatedPhone || updatedClient.Email != updatedEmail {
		t.Errorf("UpdateClient did not update contact info correctly")
	}
}

// TestUpdateClient_NotFound tests UpdateClient when the client does not exist.
func TestUpdateClient_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockClientRepository()
	service := NewClientService(mockRepo)

	// Attempt to update a client with an ID that doesn't exist.
	_, err := service.UpdateClient(999, "Jane", "Smith", "456 New Ave", "555-5678", "jane@example.com")
	if err == nil {
		t.Fatalf("Expected error for non-existent client, got nil")
	}
}

// TestRemoveClient tests RemoveClient for a client that exists.
func TestRemoveClient(t *testing.T) {
	mockRepo := mocks.NewMockClientRepository()
	service := NewClientService(mockRepo)

	// Register a client.
	now := time.Now()
	client, err := service.RegisterClient(0, "123456789", "John", "Doe", "123 Main St", "555-1234", "john@example.com", now)
	if err != nil {
		t.Fatalf("RegisterClient failed: %v", err)
	}

	// Now remove the client.
	err = service.RemoveClient(client.ID)
	if err != nil {
		t.Fatalf("RemoveClient failed: %v", err)
	}

	// Verify the client is removed.
	removedClient, err := mockRepo.FindByID(client.ID)
	if err == nil && removedClient != nil {
		t.Errorf("Expected client to be removed, but it still exists")
	}
}

// TestRemoveClient_NotFound tests RemoveClient when the client is not found.
func TestRemoveClient_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockClientRepository()
	service := NewClientService(mockRepo)

	// Attempt to remove a client that doesn't exist.
	err := service.RemoveClient(999)
	if err == nil {
		t.Fatalf("Expected error for removing non-existent client, got nil")
	}
}
