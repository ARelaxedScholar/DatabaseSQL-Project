package defaultServices_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/sql-project-backend/internal/adapters/domain/defaultServices"
	"github.com/sql-project-backend/internal/adapters/framework/driven/db/mocks"
	"github.com/sql-project-backend/internal/models"
)

// helper to create a valid Problem
func validProblem(description string) models.Problem {
	return models.Problem{
		ID:             0,
		Severity:       models.Moderate,
		Description:    description,
		SignaledWhen:   time.Now(),
		IsResolved:     false,
		ResolutionDate: time.Time{}, // unresolved, so zero time is fine
	}
}

func TestAddRoom_Success(t *testing.T) {
	mockRepo := mocks.NewMockRoomRepository()
	service := defaultServices.NewRoomService(mockRepo)

	id := 0
	hotelId := 1
	capacity := 2
	floor := 3
	price := 100.0
	telephone := "555-0101"
	viewTypes := map[models.ViewType]struct{}{
		models.Sea: {},
	}
	roomType := models.Simple
	isExtensible := false
	amenities := map[models.Amenity]struct{}{
		models.WIFI: {},
	}
	// Create a slice with one valid problem.
	problems := []models.Problem{
		validProblem("Broken window"),
	}

	room, err := service.AddRoom(id, hotelId, capacity, floor, price, telephone, viewTypes, roomType, isExtensible, amenities, problems)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if room.ID == 0 {
		t.Errorf("expected non-zero room ID, got %d", room.ID)
	}
	if room.HotelID != hotelId {
		t.Errorf("expected hotelID %d, got %d", hotelId, room.HotelID)
	}
}

func TestUpdateRoom_Success(t *testing.T) {
	mockRepo := mocks.NewMockRoomRepository()
	service := defaultServices.NewRoomService(mockRepo)

	// First add a room.
	viewTypes := map[models.ViewType]struct{}{
		models.Sea: {},
	}
	amenities := map[models.Amenity]struct{}{
		models.WIFI: {},
	}
	problems := []models.Problem{
		validProblem("Broken window"),
	}
	room, err := service.AddRoom(0, 1, 2, 3, 100.0, "555-0101", viewTypes, models.Simple, false, amenities, problems)
	if err != nil {
		t.Fatalf("failed to add room: %v", err)
	}

	// Now update the room.
	newPrice := 150.0
	newTelephone := "555-0202"
	updatedRoom, err := service.UpdateRoom(room.ID, 1, 2, 3, newPrice, newTelephone, viewTypes, models.Simple, false, amenities, problems)
	if err != nil {
		t.Fatalf("expected update to succeed, got error: %v", err)
	}
	if updatedRoom.Price != newPrice {
		t.Errorf("expected updated price %f, got %f", newPrice, updatedRoom.Price)
	}
	if updatedRoom.Telephone != newTelephone {
		t.Errorf("expected updated telephone %s, got %s", newTelephone, updatedRoom.Telephone)
	}
}

func TestUpdateRoom_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockRoomRepository()
	service := defaultServices.NewRoomService(mockRepo)

	viewTypes := map[models.ViewType]struct{}{
		models.Sea: {},
	}
	amenities := map[models.Amenity]struct{}{
		models.WIFI: {},
	}
	problems := []models.Problem{}

	// Attempt to update a room with an ID that doesn't exist.
	_, err := service.UpdateRoom(999, 1, 2, 3, 100.0, "555-0101", viewTypes, models.Simple, false, amenities, problems)
	if err == nil {
		t.Fatal("expected error for non-existent room, got nil")
	}
}

func TestDeleteRoom_Success(t *testing.T) {
	mockRepo := mocks.NewMockRoomRepository()
	service := defaultServices.NewRoomService(mockRepo)

	viewTypes := map[models.ViewType]struct{}{
		models.Sea: {},
	}
	amenities := map[models.Amenity]struct{}{
		models.WIFI: {},
	}
	problems := []models.Problem{}

	// Add a room.
	room, err := service.AddRoom(0, 1, 2, 3, 100.0, "555-0101", viewTypes, models.Simple, false, amenities, problems)
	if err != nil {
		t.Fatalf("failed to add room: %v", err)
	}

	// Delete the room.
	err = service.DeleteRoom(room.ID)
	if err != nil {
		t.Fatalf("expected delete to succeed, got error: %v", err)
	}

	// Verify deletion: FindByID should return an error or nil.
	r, err := mockRepo.FindByID(room.ID)
	if err == nil && r != nil {
		t.Errorf("expected room to be deleted, but found one")
	}
}

func TestDeleteRoom_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockRoomRepository()
	service := defaultServices.NewRoomService(mockRepo)

	err := service.DeleteRoom(999)
	if err == nil {
		t.Fatal("expected error for non-existent room, got nil")
	}
}

func TestFindAvailableRooms_Success(t *testing.T) {
	mockRepo := mocks.NewMockRoomRepository()
	service := defaultServices.NewRoomService(mockRepo)
	now := time.Now()

	viewTypes := map[models.ViewType]struct{}{
		models.Sea: {},
	}
	amenities := map[models.Amenity]struct{}{
		models.WIFI: {},
	}
	problems := []models.Problem{}

	// Add 3 rooms for hotelID=1.
	for i := 0; i < 3; i++ {
		_, err := service.AddRoom(0, 1, 2, 3, 100.0, "555-010"+strconv.Itoa(i), viewTypes, models.Simple, false, amenities, problems)
		if err != nil {
			t.Fatalf("failed to add room %d: %v", i, err)
		}
	}
	// Add 2 rooms for hotelID=2.
	for i := 0; i < 2; i++ {
		_, err := service.AddRoom(0, 2, 2, 3, 150.0, "555-020"+strconv.Itoa(i), viewTypes, models.Simple, false, amenities, problems)
		if err != nil {
			t.Fatalf("failed to add room %d for hotel 2: %v", i, err)
		}
	}
	// Act: Retrieve available rooms for hotelID=1.
	rooms, err := service.FindAvailableRooms(1, now, now.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(rooms) != 3 {
		t.Errorf("expected 3 rooms for hotel 1, got %d", len(rooms))
	}
}

func TestAssignRoomForReservation_Success(t *testing.T) {
	mockRepo := mocks.NewMockRoomRepository()
	service := defaultServices.NewRoomService(mockRepo)
	now := time.Now()

	viewTypes := map[models.ViewType]struct{}{
		models.Sea: {},
	}
	amenities := map[models.Amenity]struct{}{
		models.WIFI: {},
	}
	problems := []models.Problem{}

	// Add a room for hotelID=1.
	room, err := service.AddRoom(0, 1, 2, 3, 100.0, "555-0101", viewTypes, models.Simple, false, amenities, problems)
	if err != nil {
		t.Fatalf("failed to add room: %v", err)
	}

	// Create a dummy reservation for hotelID=1.
	reservation := &models.Reservation{
		ClientID:        1,
		HotelID:         1,
		RoomID:          0, // not yet assigned
		StartDate:       now.Add(24 * time.Hour),
		EndDate:         now.Add(48 * time.Hour),
		ReservationDate: now,
		TotalPrice:      100.0,
		Status:          models.Confirmed,
	}

	assignedRoomID, err := service.AssignRoomForReservation(reservation)
	if err != nil {
		t.Fatalf("expected room assignment to succeed, got error: %v", err)
	}
	if assignedRoomID != room.ID {
		t.Errorf("expected assigned room ID %d, got %d", room.ID, assignedRoomID)
	}
}

func TestAssignRoomForReservation_NoAvailableRooms(t *testing.T) {
	mockRepo := mocks.NewMockRoomRepository()
	service := defaultServices.NewRoomService(mockRepo)
	now := time.Now()

	// No rooms added for hotelID=1.
	reservation := &models.Reservation{
		ClientID:        1,
		HotelID:         1,
		RoomID:          0,
		StartDate:       now.Add(24 * time.Hour),
		EndDate:         now.Add(48 * time.Hour),
		ReservationDate: now,
		TotalPrice:      100.0,
		Status:          models.Confirmed,
	}

	_, err := service.AssignRoomForReservation(reservation)
	if err == nil {
		t.Fatal("expected error when no available rooms, got nil")
	}
	if err.Error() != "no available rooms" {
		t.Errorf("expected 'no available rooms' error, got: %v", err)
	}
}
