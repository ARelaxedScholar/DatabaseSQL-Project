package defaultAnonymousUseCases

import (
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultSearchRoomsUseCase struct {
	roomRepo ports.RoomRepository
}

func NewSearchRoomsUseCase(roomRepo ports.RoomRepository) ports.SearchRoomsUseCase {
	return &DefaultSearchRoomsUseCase{
		roomRepo: roomRepo,
	}
}

func (uc *DefaultSearchRoomsUseCase) SearchRooms(input dto.RoomSearchInput) (dto.RoomSearchOutput, error) {
	rooms, err := uc.roomRepo.SearchRooms(
		input.StartDate,
		input.EndDate,
		input.Capacity,
		input.PriceMin,
		input.PriceMax,
		input.HotelChainID,
		input.Category,
	)
	if err != nil {
		return dto.RoomSearchOutput{}, err
	}

	var roomOutputs []dto.RoomOutput
	for _, room := range rooms {
		// Convert the view types from map[ViewType]struct{} to []string.
		var viewTypes []string
		for vt := range room.ViewTypes {
			viewTypes = append(viewTypes, vt.String())
		}

		// Convert the amenities from map[Amenity]struct{} to []string.
		var amenities []string
		for a := range room.Amenities {
			amenities = append(amenities, a.String())
		}

		roomOutputs = append(roomOutputs, dto.RoomOutput{
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
		})
	}

	return dto.RoomSearchOutput{Rooms: roomOutputs}, nil
}
