package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

	// Returns the message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": output.Message,
	})
}

func (h *ClientHandler) MagicLogin(w http.ResponseWriter, r *http.Request) {
	// Extract the temporary token from the query parameter.
	tempToken := r.URL.Query().Get("token")
	if tempToken == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	// Call the MagicLogin use case to validate the temporary token
	// and generate a session token.
	output, err := h.LoginUseCase.MagicLogin(tempToken)
	if err != nil {
		http.Error(w, "Magic login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Securely deliver the session token by setting it in an HTTP-only cookie.
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    output.SessionToken,
		HttpOnly: true,
		Secure:   true, // Ensure HTTPS is used in production
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour), // Matches the session token's lifetime
	}
	http.SetCookie(w, cookie)

	// return a JSON response that confirms successful login.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Magic login successful.",
	})
}

func (h *ClientHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Retrieve the client ID from the request context
	clientID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
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
	clientID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	var input dto.ClientProfileUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Force the update to apply to the authenticated client.
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
	// Retrieve the authenticated client ID from the context.
	clientID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Decode the reservation input.
	var input dto.ReservationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid reservation input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set the client ID from the token context.
	input.ClientID = clientID

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
	// Extract the authenticated client ID from the request context.
	clientID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
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
	// Extract the authenticated client ID from the context.
	clientID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Retrieve the reservation ID from the URL.
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

	// Use a method that ensures the reservation belongs to the authenticated user.
	if err := h.ReservationsManagementUseCase.CancelReservation(reservationID, clientID); err != nil {
		http.Error(w, "Cancellation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
