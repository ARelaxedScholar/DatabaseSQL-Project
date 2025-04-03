package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type AdminHandler struct {
	HotelManagementUseCase   ports.AdminHotelManagementUseCase
	HotelChainUseCase        ports.AdminHotelChainManagementUseCase
	RoomManagementUseCase    ports.AdminRoomManagementUseCase
	AccountManagementUseCase ports.AdminAccountManagementUseCase
}

func NewAdminHandler(
	hotelMgmtUseCase ports.AdminHotelManagementUseCase,
	hotelChainUseCase ports.AdminHotelChainManagementUseCase,
	roomMgmtUseCase ports.AdminRoomManagementUseCase,
	accountMgmtUseCase ports.AdminAccountManagementUseCase,
) *AdminHandler {
	return &AdminHandler{
		HotelManagementUseCase:   hotelMgmtUseCase,
		HotelChainUseCase:        hotelChainUseCase,
		RoomManagementUseCase:    roomMgmtUseCase,
		AccountManagementUseCase: accountMgmtUseCase,
	}
}

func (h *AdminHandler) requireAdmin(w http.ResponseWriter, r *http.Request) (int, bool) {
    userID, ok := r.Context().Value("userID").(int)
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return 0, false
    }
    role, ok := r.Context().Value("role").(string)
    if !ok || role != "admin" {
        http.Error(w, "forbidden", http.StatusForbidden)
        return 0, false
    }
    return userID, true
}


func (h *AdminHandler) AddHotel(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	var input dto.HotelInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.HotelManagementUseCase.AddHotel(input)
	if err != nil {
		http.Error(w, "AddHotel failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *AdminHandler) UpdateHotel(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	vars := mux.Vars(r)
	hotelIDStr, ok := vars["hotelID"]
	if !ok {
		http.Error(w, "Missing hotelID in URL", http.StatusBadRequest)
		return
	}
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		http.Error(w, "Invalid hotelID", http.StatusBadRequest)
		return
	}
	var input dto.HotelInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	input.ID = hotelID
	output, err := h.HotelManagementUseCase.UpdateHotel(input)
	if err != nil {
		http.Error(w, "UpdateHotel failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *AdminHandler) DeleteHotel(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	vars := mux.Vars(r)
	hotelIDStr, ok := vars["hotelID"]
	if !ok {
		http.Error(w, "Missing hotelID in URL", http.StatusBadRequest)
		return
	}
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		http.Error(w, "Invalid hotelID", http.StatusBadRequest)
		return
	}
	if err := h.HotelManagementUseCase.DeleteHotel(hotelID); err != nil {
		http.Error(w, "DeleteHotel failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) AddHotelChain(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	var input dto.HotelChainInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.HotelChainUseCase.AddHotelChain(input)
	if err != nil {
		http.Error(w, "AddHotelChain failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *AdminHandler) UpdateHotelChain(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	vars := mux.Vars(r)
	chainIDStr, ok := vars["chainID"]
	if !ok {
		http.Error(w, "Missing chainID in URL", http.StatusBadRequest)
		return
	}
	chainID, err := strconv.Atoi(chainIDStr)
	if err != nil {
		http.Error(w, "Invalid chainID", http.StatusBadRequest)
		return
	}
	var input dto.HotelChainInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	input.ID = chainID
	output, err := h.HotelChainUseCase.UpdateHotelChain(input)
	if err != nil {
		http.Error(w, "UpdateHotelChain failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *AdminHandler) DeleteHotelChain(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	vars := mux.Vars(r)
	chainIDStr, ok := vars["chainID"]
	if !ok {
		http.Error(w, "Missing chainID in URL", http.StatusBadRequest)
		return
	}
	chainID, err := strconv.Atoi(chainIDStr)
	if err != nil {
		http.Error(w, "Invalid chainID", http.StatusBadRequest)
		return
	}
	if err := h.HotelChainUseCase.DeleteHotelChain(chainID); err != nil {
		http.Error(w, "DeleteHotelChain failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) AddRoom(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	var input dto.RoomInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.RoomManagementUseCase.AddRoom(input)
	if err != nil {
		http.Error(w, "AddRoom failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *AdminHandler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	vars := mux.Vars(r)
	roomIDStr, ok := vars["roomID"]
	if !ok {
		http.Error(w, "Missing roomID in URL", http.StatusBadRequest)
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid roomID", http.StatusBadRequest)
		return
	}
	var input dto.RoomUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	input.ID = roomID
	output, err := h.RoomManagementUseCase.UpdateRoom(input)
	if err != nil {
		http.Error(w, "UpdateRoom failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *AdminHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	vars := mux.Vars(r)
	roomIDStr, ok := vars["roomID"]
	if !ok {
		http.Error(w, "Missing roomID in URL", http.StatusBadRequest)
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid roomID", http.StatusBadRequest)
		return
	}
	if err := h.RoomManagementUseCase.DeleteRoom(roomID); err != nil {
		http.Error(w, "DeleteRoom failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	vars := mux.Vars(r)
	accountIDStr, ok := vars["accountID"]
	if !ok {
		http.Error(w, "Missing accountID in URL", http.StatusBadRequest)
		return
	}
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid accountID", http.StatusBadRequest)
		return
	}
	output, err := h.AccountManagementUseCase.GetAccount(accountID)
	if err != nil {
		http.Error(w, "GetAccount failed: "+err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *AdminHandler) ListClientAccounts(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	outputs, err := h.AccountManagementUseCase.ListClientAccounts()
	if err != nil {
		http.Error(w, "ListClientAccounts failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(outputs)
}

func (h *AdminHandler) CreateClientAccount(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	var input dto.ClientAccountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.AccountManagementUseCase.CreateClientAccount(input)
	if err != nil {
		http.Error(w, "CreateClientAccount failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *AdminHandler) UpdateClientAccount(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	vars := mux.Vars(r)
	accountIDStr, ok := vars["accountID"]
	if !ok {
		http.Error(w, "Missing client account id in URL", http.StatusBadRequest)
		return
	}
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid client account id", http.StatusBadRequest)
		return
	}
	var input dto.ClientAccountUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.AccountManagementUseCase.UpdateClientAccount(accountID, input)
	if err != nil {
		http.Error(w, "UpdateClientAccount failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *AdminHandler) DeleteClientAccount(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	vars := mux.Vars(r)
	accountIDStr, ok := vars["accountID"]
	if !ok {
		http.Error(w, "Missing client account id in URL", http.StatusBadRequest)
		return
	}
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid client account id", http.StatusBadRequest)
		return
	}
	if err := h.AccountManagementUseCase.DeleteClientAccount(accountID); err != nil {
		http.Error(w, "DeleteClientAccount failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) ListEmployeeAccounts(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	outputs, err := h.AccountManagementUseCase.ListEmployeeAccounts()
	if err != nil {
		http.Error(w, "ListEmployeeAccounts failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(outputs)
}

func (h *AdminHandler) CreateEmployeeAccount(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	var input dto.EmployeeAccountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.AccountManagementUseCase.CreateEmployeeAccount(input)
	if err != nil {
		http.Error(w, "CreateEmployeeAccount failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *AdminHandler) UpdateEmployeeAccount(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	vars := mux.Vars(r)
	accountIDStr, ok := vars["accountID"]
	if !ok {
		http.Error(w, "Missing employee account id in URL", http.StatusBadRequest)
		return
	}
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid employee account id", http.StatusBadRequest)
		return
	}
	var input dto.EmployeeAccountUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.AccountManagementUseCase.UpdateEmployeeAccount(accountID, input)
	if err != nil {
		http.Error(w, "UpdateEmployeeAccount failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *AdminHandler) DeleteEmployeeAccount(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	vars := mux.Vars(r)
	accountIDStr, ok := vars["accountID"]
	if !ok {
		http.Error(w, "Missing employee account id in URL", http.StatusBadRequest)
		return
	}
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid employee account id", http.StatusBadRequest)
		return
	}
	if err := h.AccountManagementUseCase.DeleteEmployeeAccount(accountID); err != nil {
		http.Error(w, "DeleteEmployeeAccount failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
