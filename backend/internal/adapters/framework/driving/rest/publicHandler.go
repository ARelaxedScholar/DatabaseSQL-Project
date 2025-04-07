// internal/rest/public_handler.go
package rest

import (
	"encoding/json"
	"net/http"

	"github.com/sql-project-backend/internal/ports"
)

type PublicHandler struct {
	HotelChainRepo ports.HotelChainRepository
	HotelRepo      ports.HotelRepository
	RoomTypeRepo   ports.RoomTypeRepository
}

// GET /hotelchains
func (h *PublicHandler) GetHotelChains(w http.ResponseWriter, r *http.Request) {
	chains, err := h.HotelChainRepo.ListHotelChains(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chains)
}

// GET /hotels
func (h *PublicHandler) GetHotels(w http.ResponseWriter, r *http.Request) {
	hotels, err := h.HotelRepo.ListHotels(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotels)
}

func (h *PublicHandler) GetRoomTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.RoomTypeRepo.ListRoomTypes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}
