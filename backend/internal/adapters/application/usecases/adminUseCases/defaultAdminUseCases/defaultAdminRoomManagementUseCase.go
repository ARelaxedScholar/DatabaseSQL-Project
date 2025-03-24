package defaultAdminUseCases

import (
	"errors"
	"strings"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultAdminRoomManagementUseCase struct {
	roomService ports.RoomService
}

func NewAdminRoomManagementUseCase(roomService ports.RoomService) ports.AdminRoomManagementUseCase {
	return &DefaultAdminRoomManagementUseCase{
		roomService: roomService,
	}
}

func (uc *DefaultAdminRoomManagementUseCase) AddRoom(input dto.RoomInput) (dto.RoomOutput, error) {
	vtMap, err := convertViewTypes(input.ViewTypes)
	if err != nil {
		return dto.RoomOutput{}, err
	}
	amenitiesMap, err := convertAmenities(input.Amenities)
	if err != nil {
		return dto.RoomOutput{}, err
	}
	problems, err := convertProblems(input.Problems)
	if err != nil {
		return dto.RoomOutput{}, err
	}
	roomType, err := models.ParseRoomType(input.RoomType)
	if err != nil {
		return dto.RoomOutput{}, err
	}
	room, err := uc.roomService.AddRoom(
		input.ID,
		input.HotelID,
		input.Capacity,
		input.Floor,
		input.Price,
		input.Telephone,
		vtMap,
		roomType,
		input.IsExtensible,
		amenitiesMap,
		problems,
	)
	if err != nil {
		return dto.RoomOutput{}, err
	}
	return mapRoomToOutput(room), nil
}

func (uc *DefaultAdminRoomManagementUseCase) UpdateRoom(input dto.RoomUpdateInput) (dto.RoomOutput, error) {
	vtMap, err := convertViewTypes(input.ViewTypes)
	if err != nil {
		return dto.RoomOutput{}, err
	}
	amenitiesMap, err := convertAmenities(input.Amenities)
	if err != nil {
		return dto.RoomOutput{}, err
	}
	problems, err := convertProblems(input.Problems)
	if err != nil {
		return dto.RoomOutput{}, err
	}
	roomType, err := models.ParseRoomType(input.RoomType)
	if err != nil {
		return dto.RoomOutput{}, err
	}
	room, err := uc.roomService.UpdateRoom(
		input.ID,
		input.HotelID,
		input.Capacity,
		input.Floor,
		input.Price,
		input.Telephone,
		vtMap,
		roomType,
		input.IsExtensible,
		amenitiesMap,
		problems,
	)
	if err != nil {
		return dto.RoomOutput{}, err
	}
	return mapRoomToOutput(room), nil
}

func (uc *DefaultAdminRoomManagementUseCase) DeleteRoom(roomID int) error {
	return uc.roomService.DeleteRoom(roomID)
}

func mapRoomToOutput(room *models.Room) dto.RoomOutput {
	var viewTypes []string
	for vt := range room.ViewTypes {
		viewTypes = append(viewTypes, vt.String())
	}
	var amenities []string
	for a := range room.Amenities {
		amenities = append(amenities, a.String())
	}
	var problems []string
	for _, p := range room.Problems {
		problems = append(problems, problemToString(p))
	}
	return dto.RoomOutput{
		RoomID:       room.ID,
		HotelID:      room.HotelID,
		Capacity:     room.Capacity,
		Floor:        room.Floor,
		Price:        room.Price,
		Telephone:    room.Telephone,
		ViewTypes:    viewTypes,
		RoomType:     room.RoomType.String(),
		IsExtensible: room.IsExtensible,
		Amenities:    amenities,
		Problems:     problems,
	}
}

func convertViewTypes(input []string) (map[models.ViewType]struct{}, error) {
	vtMap := make(map[models.ViewType]struct{})
	for _, v := range input {
		vt, err := models.ParseViewType(v)
		if err != nil {
			return nil, err
		}
		vtMap[vt] = struct{}{}
	}
	return vtMap, nil
}

func convertAmenities(input []string) (map[models.Amenity]struct{}, error) {
	amenitiesMap := make(map[models.Amenity]struct{})
	for _, a := range input {
		amenity, err := models.ParseAmenity(a)
		if err != nil {
			return nil, err
		}
		amenitiesMap[amenity] = struct{}{}
	}
	return amenitiesMap, nil
}

func convertProblems(input []string) ([]models.Problem, error) {
	var problems []models.Problem
	for _, s := range input {
		p, err := ParseProblem(s)
		if err != nil {
			return nil, err
		}
		problems = append(problems, p)
	}
	return problems, nil
}

func ParseProblem(s string) (models.Problem, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return models.Problem{}, errors.New("problem description cannot be empty")
	}
	// Create a new Problem with the description.
	// Adjust the structure if your models.Problem requires additional fields.
	return models.Problem{Description: s}, nil
}

func problemToString(p models.Problem) string {
	return p.Description
}
