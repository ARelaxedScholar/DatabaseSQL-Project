package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type ClientHandler struct {
	RegistrationUseCase           ports.ClientRegistrationUseCase
	LoginUseCase                  ports.ClientLoginUseCase
	ProfileUseCase                ports.ClientProfileManagementUseCase
	MakeReservationUseCase        ports.ClientMakeReservationUseCase
	ReservationsManagementUseCase ports.ClientReservationsManagementUseCase
}

func NewClientHandler(
	regUseCase ports.ClientRegistrationUseCase,
	loginUseCase ports.ClientLoginUseCase,
	profileUseCase ports.ClientProfileManagementUseCase,
	makeResUseCase ports.ClientMakeReservationUseCase,
	resManagementUseCase ports.ClientReservationsManagementUseCase,
) *ClientHandler {
	return &ClientHandler{
		RegistrationUseCase:           regUseCase,
		LoginUseCase:                  loginUseCase,
		ProfileUseCase:                profileUseCase,
		MakeReservationUseCase:        makeResUseCase,
		ReservationsManagementUseCase: resManagementUseCase,
	}
}

func (h *ClientHandler) RegisterClient(w http.ResponseWriter, r *http.Request) {
	var input dto.ClientRegistrationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid registration input: "+err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.RegistrationUseCase.RegisterClient(input)
	if err != nil {
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *ClientHandler) LoginClient(w http.ResponseWriter, r *http.Request) {
	var input dto.ClientLoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid login input: "+err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.LoginUseCase.Login(input)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *ClientHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing client id in URL", http.StatusBadRequest)
		return
	}
	clientID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid client id", http.StatusBadRequest)
		return
	}
	output, err := h.ProfileUseCase.GetProfile(clientID)
	if err != nil {
		http.Error(w, "Profile not found: "+err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *ClientHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing client id in URL", http.StatusBadRequest)
		return
	}
	clientID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid client id", http.StatusBadRequest)
		return
	}
	var input dto.ClientProfileUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	input.ClientID = clientID
	output, err := h.ProfileUseCase.UpdateProfile(input)
	if err != nil {
		http.Error(w, "Update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *ClientHandler) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var input dto.ReservationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid reservation input: "+err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.MakeReservationUseCase.MakeReservation(input)
	if err != nil {
		http.Error(w, "Reservation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *ClientHandler) ViewReservations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing client id in URL", http.StatusBadRequest)
		return
	}
	clientID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid client id", http.StatusBadRequest)
		return
	}
	outputs, err := h.ReservationsManagementUseCase.ViewReservations(clientID)
	if err != nil {
		http.Error(w, "Error retrieving reservations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(outputs)
}

func (h *ClientHandler) CancelReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservationIDStr, ok := vars["reservationID"]
	if !ok {
		http.Error(w, "Missing reservation id in URL", http.StatusBadRequest)
		return
	}
	reservationID, err := strconv.Atoi(reservationIDStr)
	if err != nil {
		http.Error(w, "Invalid reservation id", http.StatusBadRequest)
		return
	}
	if err := h.ReservationsManagementUseCase.CancelReservation(reservationID); err != nil {
		http.Error(w, "Cancellation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
