package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sql-project-backend/internal/adapters/application/jwtimpl"
	defaultAdminUseCases "github.com/sql-project-backend/internal/adapters/application/usecases/adminUseCases/defaultAdminUseCases"
	defaultAnonymousUseCases "github.com/sql-project-backend/internal/adapters/application/usecases/anonymousUseCases/defaultAnonymousUseCases"
	defaultClientUseCases "github.com/sql-project-backend/internal/adapters/application/usecases/clientUseCases/defaultClientUseCases"
	defaultEmployeeUseCases "github.com/sql-project-backend/internal/adapters/application/usecases/employeeUseCases/defaultEmployeeUseCases"
	defaultServices "github.com/sql-project-backend/internal/adapters/domain/defaultServices"
	"github.com/sql-project-backend/internal/adapters/domain/mockServices"
	myPostgreImpl "github.com/sql-project-backend/internal/adapters/framework/driven/db/sql"
	"github.com/sql-project-backend/internal/adapters/framework/driving/rest"
	"github.com/sql-project-backend/internal/models"
)

func main() {

	// Get the JWT secret key from environment variables
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("JWT_SECRET_KEY is not set")
	}
	sql_url := os.Getenv("POSTGRES_CONNECTION_URI")
	// Open database connection
	db, err := sql.Open("postgres", sql_url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close() // Important!

	// Verify connection works
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	// Instantiate a robust JWT token service.
	tokenService := jwtimpl.NewJwtTokenService(secretKey, 24*time.Hour)

	// Instantiate mock repositories.
	clientRepo, err := myPostgreImpl.NewPostgresClientRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize client repo: %v", err)
	}
	employeeRepo, err := myPostgreImpl.NewPostgresEmployeeRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize employee repo: %v", err)
	}
	hotelRepo, err := myPostgreImpl.NewPostgresHotelRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize hotel repo: %v", err)
	}
	hotelChainRepo, err := myPostgreImpl.NewPostgresHotelChainRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize hotel chain repo: %v", err)
	}
	roomRepo, err := myPostgreImpl.NewPostgresRoomRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize room repo: %v", err)
	}
	reservationRepo, err := myPostgreImpl.NewPostgresReservationRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize reservation repo: %v", err)
	}
	stayRepo, err := myPostgreImpl.NewPostgresStayRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize stay repo: %v", err)
	}
	queryRepo, err := myPostgreImpl.NewPostgresQueryRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize query repo: %v", err)
	}

	// Adding multiple rooms to your mock repo for testing.
	for i := 1; i <= 5; i++ {
		room := &models.Room{
			ID:        i,
			HotelID:   1,
			Capacity:  2,
			Price:     100.0 + float64(i*10),
			Telephone: "123-456-789" + strconv.Itoa(i),
		}
		if _, err := roomRepo.Save(room); err != nil {
			log.Fatalf("Failed to add room %d: %v", i, err)
		}
	}

	// Instantiate domain services using the repositories.
	clientService := defaultServices.NewClientService(clientRepo)
	employeeService := defaultServices.NewEmployeeService(employeeRepo)
	hotelService := defaultServices.NewHotelService(hotelRepo)
	hotelChainService := defaultServices.NewHotelChainService(hotelChainRepo)
	roomService := defaultServices.NewRoomService(roomRepo)
	reservationService := defaultServices.NewReservationService(reservationRepo)
	stayService := defaultServices.NewStayService(stayRepo)
	paymentService := mockServices.NewPaymentService()

	// Instantiate application use cases.
	registrationUseCase := defaultClientUseCases.NewClientRegistrationUseCase(clientService)
	loginUseCase := defaultClientUseCases.NewClientLoginUseCase(clientRepo, tokenService)
	profileUseCase := defaultClientUseCases.NewClientProfileManagementUseCase(clientService, clientRepo)
	makeReservationUseCase := defaultClientUseCases.NewClientMakeReservationUseCase(reservationService)
	resManagementUseCase := defaultClientUseCases.NewClientReservationsManagementUseCase(reservationService)
	searchRoomsUseCase := defaultAnonymousUseCases.NewSearchRoomsUseCase(roomRepo, queryRepo)

	employeeLoginUseCase := defaultEmployeeUseCases.NewEmployeeLoginUseCase(employeeRepo, tokenService)
	checkInUseCase := defaultEmployeeUseCases.NewEmployeeCheckInUseCase(stayService, reservationRepo, roomRepo)
	createNewStayUseCase := defaultEmployeeUseCases.NewEmployeeCreateNewStayUseCase(stayService)
	checkoutUseCase := defaultEmployeeUseCases.NewEmployeeCheckoutUseCase(stayService, paymentService)

	adminHotelManagementUseCase := defaultAdminUseCases.NewAdminHotelManagementUseCase(hotelService)
	adminHotelChainUseCase := defaultAdminUseCases.NewAdminHotelChainManagementUseCase(hotelChainService)
	adminRoomManagementUseCase := defaultAdminUseCases.NewAdminRoomManagementUseCase(roomService, roomRepo)
	adminAccountManagementUseCase := defaultAdminUseCases.NewAdminAccountManagementUseCase(clientRepo, employeeRepo, clientService, employeeService)

	// Instantiate REST handlers.
	clientHandler := rest.NewClientHandler(registrationUseCase, loginUseCase, profileUseCase, makeReservationUseCase, resManagementUseCase)
	employeeHandler := rest.NewEmployeeHandler(employeeLoginUseCase, checkInUseCase, createNewStayUseCase, checkoutUseCase)
	adminHandler := rest.NewAdminHandler(adminHotelManagementUseCase, adminHotelChainUseCase, adminRoomManagementUseCase, adminAccountManagementUseCase)
	anonymousHandler := rest.NewAnonymousHandler(searchRoomsUseCase)

	// Set up Gorilla Mux router.
	router := mux.NewRouter()

	// Client routes.
	router.HandleFunc("/clients/register", clientHandler.RegisterClient).Methods("POST")
	router.HandleFunc("/clients/login", clientHandler.LoginClient).Methods("POST")
	router.HandleFunc("/clients/{id:[0-9]+}/profile", clientHandler.GetProfile).Methods("GET")
	router.HandleFunc("/clients/{id:[0-9]+}/profile", clientHandler.UpdateProfile).Methods("PUT", "PATCH")
	router.HandleFunc("/clients/reservations", clientHandler.MakeReservation).Methods("POST")
	router.HandleFunc("/clients/{id:[0-9]+}/reservations", clientHandler.ViewReservations).Methods("GET")
	router.HandleFunc("/clients/reservations/{reservationID:[0-9]+}", clientHandler.CancelReservation).Methods("DELETE")

	// Employee routes.
	router.HandleFunc("/employees/login", employeeHandler.LoginEmployee).Methods("POST")
	router.HandleFunc("/employees/checkin", employeeHandler.CheckIn).Methods("POST")
	router.HandleFunc("/employees/stay", employeeHandler.CreateNewStay).Methods("POST")
	// New checkout route for employees.
	router.HandleFunc("/employees/checkout", employeeHandler.Checkout).Methods("POST")

	// Admin routes.
	router.HandleFunc("/admin/hotels", adminHandler.AddHotel).Methods("POST")
	router.HandleFunc("/admin/hotels/{hotelID:[0-9]+}", adminHandler.UpdateHotel).Methods("PUT", "PATCH")
	router.HandleFunc("/admin/hotels/{hotelID:[0-9]+}", adminHandler.DeleteHotel).Methods("DELETE")

	router.HandleFunc("/admin/hotelchains", adminHandler.AddHotelChain).Methods("POST")
	router.HandleFunc("/admin/hotelchains/{chainID:[0-9]+}", adminHandler.UpdateHotelChain).Methods("PUT", "PATCH")
	router.HandleFunc("/admin/hotelchains/{chainID:[0-9]+}", adminHandler.DeleteHotelChain).Methods("DELETE")

	router.HandleFunc("/admin/rooms", adminHandler.AddRoom).Methods("POST")
	router.HandleFunc("/admin/rooms/{roomID:[0-9]+}", adminHandler.UpdateRoom).Methods("PUT", "PATCH")
	router.HandleFunc("/admin/rooms/{roomID:[0-9]+}", adminHandler.DeleteRoom).Methods("DELETE")

	router.HandleFunc("/admin/accounts/{accountID:[0-9]+}", adminHandler.GetAccount).Methods("GET")
	router.HandleFunc("/admin/accounts/clients", adminHandler.ListClientAccounts).Methods("GET")
	router.HandleFunc("/admin/accounts/clients", adminHandler.CreateClientAccount).Methods("POST")
	router.HandleFunc("/admin/accounts/clients/{accountID:[0-9]+}", adminHandler.UpdateClientAccount).Methods("PUT", "PATCH")
	router.HandleFunc("/admin/accounts/clients/{accountID:[0-9]+}", adminHandler.DeleteClientAccount).Methods("DELETE")
	router.HandleFunc("/admin/accounts/employees", adminHandler.ListEmployeeAccounts).Methods("GET")
	router.HandleFunc("/admin/accounts/employees", adminHandler.CreateEmployeeAccount).Methods("POST")
	router.HandleFunc("/admin/accounts/employees/{accountID:[0-9]+}", adminHandler.UpdateEmployeeAccount).Methods("PUT", "PATCH")
	router.HandleFunc("/admin/accounts/employees/{accountID:[0-9]+}", adminHandler.DeleteEmployeeAccount).Methods("DELETE")

	// Anonymous route.
	router.HandleFunc("/search/rooms", anonymousHandler.SearchRooms).Methods("GET")
	router.HandleFunc("/search/hotels/{hotelID:[0-9]+}/room-count", anonymousHandler.CountRoomsInHotel).Methods("GET")
	router.HandleFunc("/search/zones/rooms", anonymousHandler.GetRoomsByZone).Methods("GET")

	log.Println("Server is running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
