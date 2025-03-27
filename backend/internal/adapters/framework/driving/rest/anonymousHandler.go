package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type AnonymousHandler struct {
	SearchRoomsUseCase ports.SearchRoomsUseCase
}

func NewAnonymousHandler(searchRoomsUseCase ports.SearchRoomsUseCase) *AnonymousHandler {
	return &AnonymousHandler{
		SearchRoomsUseCase: searchRoomsUseCase,
	}
}

func (h *AnonymousHandler) SearchRooms(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	startDate, err := parseTimeParam(q.Get("startDate"))
	if err != nil {
		http.Error(w, "invalid startDate: "+err.Error(), http.StatusBadRequest)
		return
	}
	endDate, err := parseTimeParam(q.Get("endDate"))
	if err != nil {
		http.Error(w, "invalid endDate: "+err.Error(), http.StatusBadRequest)
		return
	}
	capacity, err := parseIntParam(q.Get("capacity"))
	if err != nil {
		http.Error(w, "invalid capacity: "+err.Error(), http.StatusBadRequest)
		return
	}
	priceMin, err := parseFloatParam(q.Get("priceMin"))
	if err != nil {
		http.Error(w, "invalid priceMin: "+err.Error(), http.StatusBadRequest)
		return
	}
	priceMax, err := parseFloatParam(q.Get("priceMax"))
	if err != nil {
		http.Error(w, "invalid priceMax: "+err.Error(), http.StatusBadRequest)
		return
	}
	hotelChainID, err := parseIntParam(q.Get("hotelChainID"))
	if err != nil {
		http.Error(w, "invalid hotelChainID: "+err.Error(), http.StatusBadRequest)
		return
	}
	roomType := q.Get("roomType")

	input := dto.RoomSearchInput{
		StartDate:    startDate,
		EndDate:      endDate,
		Capacity:     capacity,
		PriceMin:     priceMin,
		PriceMax:     priceMax,
		HotelChainID: hotelChainID,
		RoomType:     roomType,
	}

	output, err := h.SearchRoomsUseCase.SearchRooms(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func parseTimeParam(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	// Assume the timestamp is in seconds.
	seconds, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(seconds, 0), nil
}

func parseIntParam(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func parseFloatParam(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}
