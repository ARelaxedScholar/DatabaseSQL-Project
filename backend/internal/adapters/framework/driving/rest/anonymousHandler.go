package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

// New additions from the Query repo.
func (h *AnonymousHandler) CountRoomsInHotel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["hotelID"]
	if !ok {
		http.Error(w, "hotelID missing", http.StatusBadRequest)
		return
	}
	hotelID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid hotelID", http.StatusBadRequest)
		return
	}

	log.Printf("CountRoomsInHotel called with hotelID=%d\n", hotelID)

	count, err := h.SearchRoomsUseCase.GetNumberOfRoomsForHotel(hotelID)
	if err != nil {
		log.Printf("ERROR CountRoomsInHotel SQL: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]int{"total_capacity": count}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("ERROR encoding JSON: %v\n", err)
	}
}

func (h *AnonymousHandler) GetRoomsByZone(w http.ResponseWriter, r *http.Request) {
	output, err := h.SearchRoomsUseCase.GetNumberOfRoomsPerZone()
	if err != nil {
		log.Printf("Error fetching rooms by zone: %v", err) // Log internally
		http.Error(w, "Could not fetch rooms by zone", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(output); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *AnonymousHandler) SearchRooms(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// For each parameter, if the query param is not empty, parse it and take its address.
	var startDate *time.Time
	if s := q.Get("startDate"); s != "" {
		t, err := parseTimeParam(s)
		if err != nil {
			http.Error(w, "invalid startDate: "+err.Error(), http.StatusBadRequest)
			return
		}
		startDate = &t
	}

	var endDate *time.Time
	if s := q.Get("endDate"); s != "" {
		t, err := parseTimeParam(s)
		if err != nil {
			http.Error(w, "invalid endDate: "+err.Error(), http.StatusBadRequest)
			return
		}
		endDate = &t
	}

	var capacity *int
	if s := q.Get("capacity"); s != "" {
		c, err := parseIntParam(s)
		if err != nil {
			http.Error(w, "invalid capacity: "+err.Error(), http.StatusBadRequest)
			return
		}
		capacity = &c
	}

	var priceMin *float64
	if s := q.Get("priceMin"); s != "" {
		p, err := parseFloatParam(s)
		if err != nil {
			http.Error(w, "invalid priceMin: "+err.Error(), http.StatusBadRequest)
			return
		}
		priceMin = &p
	}

	var priceMax *float64
	if s := q.Get("priceMax"); s != "" {
		p, err := parseFloatParam(s)
		if err != nil {
			http.Error(w, "invalid priceMax: "+err.Error(), http.StatusBadRequest)
			return
		}
		priceMax = &p
	}

	var hotelChainID *int
	if s := q.Get("hotelChainID"); s != "" {
		id, err := parseIntParam(s)
		if err != nil {
			http.Error(w, "invalid hotelChainID: "+err.Error(), http.StatusBadRequest)
			return
		}
		hotelChainID = &id
	}

	var roomType *string
	if s := q.Get("roomType"); s != "" {
		roomType = &s
	}

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
	// Parse the date in the "MM-DD-YYYY" format.
	t, err := time.Parse("01-02-2006", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
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
