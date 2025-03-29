package defaultAnonymousUseCases

import (
	"fmt"
	"time"

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
	var searchRoomType models.RoomType
	var err error
	// Only parse RoomType if provided (non-nil and non-empty)
	if input.RoomType != nil && *input.RoomType != "" {
		searchRoomType, err = models.ParseRoomType(*input.RoomType)
		if err != nil {
			return dto.RoomSearchOutput{}, fmt.Errorf("Invalid room category provided for search: %w", err)
		}
	}

	// Convert pointer fields to concrete values (or use zero values)
	var (
		startDate    time.Time
		endDate      time.Time
		capacity     int
		priceMin     float64
		priceMax     float64
		hotelChainID int
	)

	if input.StartDate != nil {
		startDate = *input.StartDate
	} else {
		startDate = time.Time{}
	}

	if input.EndDate != nil {
		endDate = *input.EndDate
	} else {
		endDate = time.Time{}
	}

	if input.Capacity != nil {
		capacity = *input.Capacity
	} else {
		capacity = 0
	}

	if input.PriceMin != nil {
		priceMin = *input.PriceMin
	} else {
		priceMin = 0.0
	}

	if input.PriceMax != nil {
		priceMax = *input.PriceMax
	} else {
		priceMax = 0.0
	}

	if input.HotelChainID != nil {
		hotelChainID = *input.HotelChainID
	} else {
		hotelChainID = 0
	}

	// Call the repository using the concrete values.
	rooms, err := uc.roomRepo.SearchRooms(
		startDate,      // time.Time
		endDate,        // time.Time
		capacity,       // int
		priceMin,       // float64
		priceMax,       // float64
		hotelChainID,   // int
		searchRoomType, // models.RoomType (zero value if no RoomType provided)
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
