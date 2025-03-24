package defaultServices_test

import (
	"testing"

	"github.com/sql-project-backend/internal/adapters/domain/defaultServices"
	"github.com/sql-project-backend/internal/adapters/framework/driven/db/mocks"
)

func TestGetAvailableRoomsByZone(t *testing.T) {
	// Arrange: Create the mock query repository.
	mockQueryRepo := mocks.NewMockQueryRepository()
	// Instantiate the query service.
	queryService := defaultServices.NewQueryService(mockQueryRepo)

	// Act: Get the available rooms by zone.
	result, err := queryService.GetAvailableRoomsByZone()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Assert: Check that the result matches the expected map.
	// Based on the mock implementation, we expect:
	//   "North": 10, "South": 8, "East": 12, "West": 5
	expected := map[string]int{
		"North": 10,
		"South": 8,
		"East":  12,
		"West":  5,
	}
	if len(result) != len(expected) {
		t.Fatalf("expected %d zones, got %d", len(expected), len(result))
	}
	for zone, expCap := range expected {
		if cap, ok := result[zone]; !ok {
			t.Errorf("expected zone %s not found", zone)
		} else if cap != expCap {
			t.Errorf("expected capacity for zone %s to be %d, got %d", zone, expCap, cap)
		}
	}
}

func TestGetHotelRoomCapacity_Success(t *testing.T) {
	// Arrange: Create the mock query repository.
	mockQueryRepo := mocks.NewMockQueryRepository()
	// Instantiate the query service.
	queryService := defaultServices.NewQueryService(mockQueryRepo)
	
	// Act: Request room capacity for hotel with ID 1.
	capacity, err := queryService.GetHotelRoomCapacity(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Assert: Based on our mock, hotel with ID 1 should have a capacity of 100.
	if capacity != 100 {
		t.Errorf("expected capacity 100 for hotel id 1, got %d", capacity)
	}
}

func TestGetHotelRoomCapacity_HotelNotFound(t *testing.T) {
	// Arrange: Create the mock query repository.
	mockQueryRepo := mocks.NewMockQueryRepository()
	// Instantiate the query service.
	queryService := defaultServices.NewQueryService(mockQueryRepo)
	
	// Act: Request room capacity for a hotel id that does not exist, e.g., 999.
	_, err := queryService.GetHotelRoomCapacity(999)
	
	// Assert: Expect an error indicating the hotel was not found.
	if err == nil {
		t.Fatal("expected error for non-existent hotel id, got nil")
	}
}
