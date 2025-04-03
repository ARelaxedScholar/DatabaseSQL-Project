package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

// EmployeeHandler defines the REST endpoints for employee-related operations.
type EmployeeHandler struct {
	LoginUseCase         ports.EmployeeLoginUseCase
	CheckInUseCase       ports.EmployeeCheckInUseCase
	CreateNewStayUseCase ports.EmployeeCreateNewStayUseCase
	CheckoutUseCase      ports.EmployeeCheckoutUseCase // New field for checkout use case
}

// NewEmployeeHandler constructs a new EmployeeHandler.
func NewEmployeeHandler(
	loginUseCase ports.EmployeeLoginUseCase,
	checkInUseCase ports.EmployeeCheckInUseCase,
	createNewStayUseCase ports.EmployeeCreateNewStayUseCase,
	checkoutUseCase ports.EmployeeCheckoutUseCase,
) *EmployeeHandler {
	return &EmployeeHandler{
		LoginUseCase:         loginUseCase,
		CheckInUseCase:       checkInUseCase,
		CreateNewStayUseCase: createNewStayUseCase,
		CheckoutUseCase:      checkoutUseCase,
	}
}

// LoginEmployee is a public endpoint that handles employee login.
func (h *EmployeeHandler) LoginEmployee(w http.ResponseWriter, r *http.Request) {
	var input dto.EmployeeLoginInput
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

func (h *EmployeeHandler) MagicLogin(w http.ResponseWriter, r *http.Request) {
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

// CheckIn is a protected endpoint that allows an authenticated employee to check in.
func (h *EmployeeHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	// Retrieve the authenticated employee ID from the context.
	employeeID, ok := r.Context().Value("employeeID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var input dto.CheckInInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid check-in input: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Set the employee ID from the context.
	input.EmployeeID = employeeID

	output, err := h.CheckInUseCase.CheckIn(input)
	if err != nil {
		http.Error(w, "Check-in failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// CreateNewStay is a protected endpoint that allows an authenticated employee to create a new stay.
func (h *EmployeeHandler) CreateNewStay(w http.ResponseWriter, r *http.Request) {
	// Retrieve the authenticated employee ID from the context.
	employeeID, ok := r.Context().Value("employeeID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var input dto.NewStayInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid new stay input: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Set the employee ID from the context as the check-in employee.
	input.CheckInEmployeeID = employeeID

	output, err := h.CreateNewStayUseCase.CreateNewStay(input)
	if err != nil {
		http.Error(w, "Create new stay failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// Checkout is a protected endpoint that allows an authenticated employee to process a checkout.
func (h *EmployeeHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	// Retrieve the authenticated employee ID from the context.
	employeeID, ok := r.Context().Value("employeeID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var input dto.CheckoutInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid checkout input: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Log the employee ID processing the checkout for auditing/debugging.
	log.Printf("Employee %d processing checkout for stay %d", employeeID, input.StayID)

	input.EmpoyeeID = employeeID // needed to update the stays

	output, err := h.CheckoutUseCase.Checkout(input)
	if err != nil {
		http.Error(w, "Checkout failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
