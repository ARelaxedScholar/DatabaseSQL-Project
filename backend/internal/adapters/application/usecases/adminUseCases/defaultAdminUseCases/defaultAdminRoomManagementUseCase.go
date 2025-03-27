package defaultAdminUseCases

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultAdminRoomManagementUseCase struct {
	roomService ports.RoomService
	roomRepo    ports.RoomRepository
}

func NewAdminRoomManagementUseCase(roomService ports.RoomService, roomRepository ports.RoomRepository) ports.AdminRoomManagementUseCase {
	return &DefaultAdminRoomManagementUseCase{
		roomService: roomService,
		roomRepo:    roomRepository,
	}
}

// AddRoom remains the same
func (uc *DefaultAdminRoomManagementUseCase) AddRoom(input dto.RoomInput) (dto.RoomOutput, error) {
	vtMap, err := convertViewTypes(input.ViewTypes)
	if err != nil {
		return dto.RoomOutput{}, fmt.Errorf("Failed to convert view types: %w", err)
	}
	amenitiesMap, err := convertAmenities(input.Amenities)
	if err != nil {
		return dto.RoomOutput{}, fmt.Errorf("Failed to convert amenities: %w", err)
	}
	problems, err := convertProblems(input.Problems)
	if err != nil {
		return dto.RoomOutput{}, fmt.Errorf("Failed to convert problems: %w", err)
	}
	roomType, err := models.ParseRoomType(input.RoomType)
	if err != nil {
		return dto.RoomOutput{}, fmt.Errorf("Failed to parse room type: %w", err)
	}

	room, err := uc.roomService.AddRoom(
		0, input.HotelID, input.Capacity, input.Number, input.Floor, input.SurfaceArea,
		input.Price, input.Telephone, vtMap, roomType, input.IsExtensible,
		amenitiesMap, problems,
	)
	if err != nil {
		return dto.RoomOutput{}, err
	}
	return mapRoomToOutput(room), nil
}

func (uc *DefaultAdminRoomManagementUseCase) UpdateRoom(input dto.RoomUpdateInput) (dto.RoomOutput, error) {
	// 1. Fetch the existing room using the service's FindByID method.
	existingRoom, err := uc.roomRepo.FindByID(input.ID)
	if err != nil {
		// Handle errors like "not found" appropriately
		return dto.RoomOutput{}, fmt.Errorf("Failed to find room with ID %d for update: %w", input.ID, err)
	}

	// 2. Prepare variables for the service call, defaulting to existing values.
	// Use the values from the fetched existingRoom as the base.
	hotelID := existingRoom.HotelID
	capacity := existingRoom.Capacity
	number := existingRoom.Number // string
	floor := existingRoom.Floor   // string
	surfaceArea := existingRoom.SurfaceArea
	price := existingRoom.Price
	telephone := existingRoom.Telephone
	viewTypes := existingRoom.ViewTypes // Start with existing map
	roomType := existingRoom.RoomType
	isExtensible := existingRoom.IsExtensible
	amenities := existingRoom.Amenities // Start with existing map
	problems := existingRoom.Problems   // Start with existing slice

	// 3. Apply updates from the DTO if the corresponding pointer is not nil.
	if input.HotelID != nil {
		hotelID = *input.HotelID
	}
	if input.Capacity != nil {
		capacity = *input.Capacity
	}
	if input.Number != nil {
		number = *input.Number // Update string
	}
	if input.Floor != nil {
		floor = *input.Floor // Update string
	}
	if input.SurfaceArea != nil {
		surfaceArea = *input.SurfaceArea
	}
	if input.Price != nil {
		price = *input.Price
	}
	if input.Telephone != nil {
		telephone = *input.Telephone
	}
	if input.IsExtensible != nil {
		isExtensible = *input.IsExtensible
	}

	// Apply complex type updates if provided in DTO
	if input.ViewTypes != nil {
		vtMap, convErr := convertViewTypes(*input.ViewTypes)
		if convErr != nil {
			return dto.RoomOutput{}, fmt.Errorf("Failed to convert view types: %w", convErr)
		}
		viewTypes = vtMap // Replace map with new values from DTO
	}
	if input.Amenities != nil {
		amenitiesMap, convErr := convertAmenities(*input.Amenities)
		if convErr != nil {
			return dto.RoomOutput{}, fmt.Errorf("Failed to convert amenities: %w", convErr)
		}
		amenities = amenitiesMap // Replace map with new values from DTO
	}
	if input.Problems != nil {
		problemsSlice, convErr := convertProblems(*input.Problems)
		if convErr != nil {
			return dto.RoomOutput{}, fmt.Errorf("Failed to convert problems: %w", convErr)
		}
		problems = problemsSlice // Replace slice with new values from DTO
	}
	if input.RoomType != nil {
		rt, parseErr := models.ParseRoomType(*input.RoomType)
		if parseErr != nil {
			return dto.RoomOutput{}, fmt.Errorf("Failed to parse room type: %w", parseErr)
		}
		roomType = rt // Update room type enum
	}

	// 4. Call the service's original UpdateRoom method with the fully populated (merged) data.
	updatedRoom, err := uc.roomService.UpdateRoom(
		input.ID,                                                        // Use the ID from the input DTO to identify the room
		hotelID, capacity, number, floor, surfaceArea, price, telephone, // Pass merged values
		viewTypes, roomType, isExtensible, amenities, problems,
	)
	if err != nil {
		return dto.RoomOutput{}, err // Propagate service errors
	}

	// 5. Map the result from the service call back to the output DTO.
	return mapRoomToOutput(updatedRoom), nil
}

func (uc *DefaultAdminRoomManagementUseCase) DeleteRoom(roomID int) error {
	return uc.roomService.DeleteRoom(roomID)
}

// --- Helper functions remain the same ---
func mapRoomToOutput(room *models.Room) dto.RoomOutput {
	if room == nil {
		return dto.RoomOutput{}
	}
	viewTypes := make([]string, 0, len(room.ViewTypes))
	for vt := range room.ViewTypes {
		viewTypes = append(viewTypes, vt.String())
	}
	amenities := make([]string, 0, len(room.Amenities))
	for a := range room.Amenities {
		amenities = append(amenities, a.String())
	}
	problems := make([]string, 0, len(room.Problems))
	for _, p := range room.Problems {
		problems = append(problems, problemToString(p))
	}
	return dto.RoomOutput{
		RoomID: room.ID, HotelID: room.HotelID, Capacity: room.Capacity, Number: room.Number,
		Floor: room.Floor, Price: room.Price, Telephone: room.Telephone, ViewTypes: viewTypes,
		RoomType: room.RoomType.String(), IsExtensible: room.IsExtensible, Amenities: amenities,
		Problems: problems,
	}
}
func convertViewTypes(input []string) (map[models.ViewType]struct{}, error) {
	vtMap := make(map[models.ViewType]struct{})
	if input == nil {
		return vtMap, nil
	}
	for _, v := range input {
		vt, err := models.ParseViewType(v)
		if err != nil {
			return nil, fmt.Errorf("Invalid view type '%s': %w", v, err)
		}
		vtMap[vt] = struct{}{}
	}
	return vtMap, nil
}
func convertAmenities(input []string) (map[models.Amenity]struct{}, error) {
	amenitiesMap := make(map[models.Amenity]struct{})
	if input == nil {
		return amenitiesMap, nil
	}
	for _, a := range input {
		amenity, err := models.ParseAmenity(a)
		if err != nil {
			return nil, fmt.Errorf("Invalid amenity '%s': %w", a, err)
		}
		amenitiesMap[amenity] = struct{}{}
	}
	return amenitiesMap, nil
}
func convertProblems(input []string) ([]models.Problem, error) {
	problems := make([]models.Problem, 0, len(input))
	if input == nil {
		return problems, nil
	}
	for _, s := range input {
		p, err := ParseProblem(s)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse problem string '%s': %w", s, err)
		}
		problems = append(problems, p)
	}
	return problems, nil
}
func ParseProblem(s string) (models.Problem, error) {
	desc := strings.TrimSpace(s)
	if desc == "" {
		return models.Problem{}, errors.New("Problem description cannot be empty.")
	}
	severity, _ := models.ParseProblemSeverity("Moderate")
	return models.Problem{Description: desc, Severity: severity, SignaledWhen: time.Now(), IsResolved: false}, nil
}
func problemToString(p models.Problem) string { return p.Description }

var _ ports.AdminRoomManagementUseCase = (*DefaultAdminRoomManagementUseCase)(nil)
