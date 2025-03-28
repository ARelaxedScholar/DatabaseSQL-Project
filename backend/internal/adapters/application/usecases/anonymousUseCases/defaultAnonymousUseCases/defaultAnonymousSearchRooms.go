package defaultAnonymousUseCases

import (
	"fmt"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultSearchRoomsUseCase struct {
	roomRepo  ports.RoomRepository
	queryRepo ports.QueryRepository
}

func NewSearchRoomsUseCase(roomRepo ports.RoomRepository, queryRepo ports.QueryRepository) ports.SearchRoomsUseCase {
	return &DefaultSearchRoomsUseCase{
		roomRepo:  roomRepo,
		queryRepo: queryRepo,
	}
}

func (uc *DefaultSearchRoomsUseCase) SearchRooms(input dto.RoomSearchInput) (dto.RoomSearchOutput, error) {
	// parse the potentially empty room type
	var searchRoomType models.RoomType
	var err error
	if input.RoomType != "" { // Only parse if a category is provided
		searchRoomType, err = models.ParseRoomType(input.RoomType)
		if err != nil {
			return dto.RoomSearchOutput{}, fmt.Errorf("Invalid room category provided for search: %w", err)
		}
	}

	rooms, err := uc.roomRepo.SearchRooms(
		input.StartDate,
		input.EndDate,
		input.Capacity,
		input.PriceMin,
		input.PriceMax,
		input.HotelChainID,
		searchRoomType, // Pass the parsed RoomType enum value
	)
	if err != nil {
		return dto.RoomSearchOutput{}, err
	}

	roomOutputs := make([]dto.RoomOutput, 0, len(rooms))
	for _, room := range rooms {
		if room == nil {
			continue
		}

		viewTypes := make([]string, 0, len(room.ViewTypes))
		for vt := range room.ViewTypes {
			viewTypes = append(viewTypes, vt.String())
		}

		amenities := make([]string, 0, len(room.Amenities))
		for a := range room.Amenities {
			amenities = append(amenities, a.String())
		}

		roomOutputs = append(roomOutputs, dto.RoomOutput{
			RoomID:       room.ID,
			HotelID:      room.HotelID,
			Capacity:     room.Capacity,
			Number:       room.Number,
			Floor:        room.Floor,
			Price:        room.Price,
			Telephone:    room.Telephone,
			ViewTypes:    viewTypes,
			RoomType:     room.RoomType.String(),
			IsExtensible: room.IsExtensible,
			Amenities:    amenities,
		})
	}

	return dto.RoomSearchOutput{Rooms: roomOutputs}, nil
}

// implemented the 
func (s DefaultSearchRoomsUseCase) GetNumberOfRoomsForHotel(hotelID int) (int, error) {
	return s.queryRepo.GetHotelRoomCapacity(hotelID)
}

func (s DefaultSearchRoomsUseCase) GetNumberOfRoomsPerZone() (map[string]int, error) {
	return s.queryRepo.GetAvailableRoomsByZone()
}
