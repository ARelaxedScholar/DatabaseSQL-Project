package rest

import (
	"encoding/json"
	"net/http"

	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type EmployeeHandler struct {
	LoginUseCase         ports.EmployeeLoginUseCase
	CheckInUseCase       ports.EmployeeCheckInUseCase
	CreateNewStayUseCase ports.EmployeeCreateNewStayUseCase
}

func NewEmployeeHandler(loginUseCase ports.EmployeeLoginUseCase, checkInUseCase ports.EmployeeCheckInUseCase, createNewStayUseCase ports.EmployeeCreateNewStayUseCase) *EmployeeHandler {
	return &EmployeeHandler{
		LoginUseCase:         loginUseCase,
		CheckInUseCase:       checkInUseCase,
		CreateNewStayUseCase: createNewStayUseCase,
	}
}

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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *EmployeeHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
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
	input.EmployeeID = employeeID
	output, err := h.CheckInUseCase.CheckIn(input)
	if err != nil {
		http.Error(w, "Check-in failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *EmployeeHandler) CreateNewStay(w http.ResponseWriter, r *http.Request) {
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
	input.CheckInEmployeeID = employeeID
	output, err := h.CreateNewStayUseCase.CreateNewStay(input)
	if err != nil {
		http.Error(w, "Create new stay failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
