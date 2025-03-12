document.addEventListener("DOMContentLoaded", function () {
    const loginForm = document.getElementById("loginForm");
    const bookingSection = document.getElementById("booking-section");
    const employeeDashboard = document.getElementById("employeeDashboard");
    const searchForm = document.getElementById("searchForm");
    const availableRooms = document.getElementById("roomResults");
    const convertForm = document.getElementById("convertForm");
    const directRentalForm = document.getElementById("directRentalForm");
    const paymentForm = document.getElementById("paymentForm");

    loginForm.addEventListener("submit", function (event) {
        event.preventDefault();
        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;
        const role = document.querySelector("input[name='role']:checked").value;

        if (username && password) {
            if (role === "client") {
                document.getElementById("login-form").classList.remove("active");
                bookingSection.classList.add("active");
            } else if (role === "employe") {
                document.getElementById("login-form").classList.remove("active");
                employeeDashboard.classList.add("active");
            }
        }
    });

    searchForm.addEventListener("submit", function (event) {
        event.preventDefault();
        availableRooms.innerHTML = "";

        const checkIn = document.getElementById("check-in").value;
        const checkOut = document.getElementById("check-out").value;
        const capacity = document.getElementById("capacity").value;
        const price = document.getElementById("price").value;
        const hotelChain = document.getElementById("hotel-chain").value;
        const category = document.getElementById("category").value;

        if (checkIn && checkOut) {
            const mockRooms = [
                { id: 1, name: "Chambre Luxe", capacity: 2, price: 150, hotelChain: "luxury", category: "5-star" },
                { id: 2, name: "Suite Royale", capacity: 4, price: 300, hotelChain: "luxury", category: "5-star" },
                { id: 3, name: "Chambre Standard", capacity: 2, price: 100, hotelChain: "budget", category: "3-star" }
            ];

            mockRooms.forEach(room => {
                if (
                    room.capacity >= capacity &&
                    room.price <= price &&
                    (!hotelChain || room.hotelChain === hotelChain) &&
                    (!category || room.category === category)
                ) {
                    const roomDiv = document.createElement("div");
                    roomDiv.classList.add("room");
                    roomDiv.innerHTML = `<h3>${room.name}</h3><p>Capacité: ${room.capacity}</p><p>Prix: ${room.price}€ / nuit</p><button onclick="reserveRoom(${room.id})">Réserver</button>`;
                    availableRooms.appendChild(roomDiv);
                }
            });
        }
    });

    convertForm.addEventListener("submit", function (event) {
        event.preventDefault();
        const reservationId = document.getElementById("reservation-id").value;
        alert(`Réservation ID ${reservationId} convertie en location!`);
    });

    directRentalForm.addEventListener("submit", function (event) {
        event.preventDefault();
        const roomId = document.getElementById("room-id").value;
        const rentalDate = document.getElementById("rental-date").value;
        alert(`Location créée pour la chambre ${roomId} à la date ${rentalDate}`);
    });

    paymentForm.addEventListener("submit", function (event) {
        event.preventDefault();
        const rentalId = document.getElementById("rental-id").value;
        const paymentAmount = document.getElementById("payment-amount").value;
        alert(`Paiement de ${paymentAmount}€ enregistré pour la location ${rentalId}`);
    });
});

function reserveRoom(roomId) {
    alert(`Réservation de chambre ${roomId} effectuée!`);
}
