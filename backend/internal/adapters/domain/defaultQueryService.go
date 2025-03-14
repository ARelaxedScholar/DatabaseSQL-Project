package domain

import (
	"github.com/sql-project-backend/internal/ports"
)

type DefaultQueryService struct {
	queryRepo ports.QueryRepository
}

func NewQueryService(queryRepo ports.QueryRepository) ports.QueryService {
	return &DefaultQueryService{
		queryRepo: queryRepo,
	}
}

func (s *DefaultQueryService) GetAvailableRoomsByZone() (map[string]int, error) {
	return s.queryRepo.GetAvailableRoomsByZone()
}

func (s *DefaultQueryService) GetHotelRoomCapacity(hotelId int) (int, error) {
	return s.queryRepo.GetHotelRoomCapacity(hotelId)
}
