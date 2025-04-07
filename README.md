# e-Hôtels — Reservation Management System

**Authors:** Group 15  
Ange Jedidiah Kouakou (300314061)  
Jepherson Fenelon (300291711)  
Mia Partington (300291725)  
Patrick Bonini (300273064)

---

## 🚀 Table of Contents

1. [Description](#description)  
2. [Features](#features)  
3. [Tech Stack](#tech-stack)  
4. [Prerequisites](#prerequisites)  
5. [Getting Started](#getting-started)  
6. [Project Structure](#project-structure)  
7. [API Reference](#api-reference)  
8. [Testing](#testing)  
9. [Deployment](#deployment)  
10. [Contributing](#contributing)  
11. [License](#license)

---

## Description

A real‑time hotel reservation system that allows clients to search for and book rooms, and employees to manage bookings, check‑in/out, and view analytics by zone and capacity.

---

## Features

### Client
- Magic‑link authentication (no passwords)  
- Search by date, capacity, price, hotel chain, room type  
- Book, view and cancel reservations  

### Employee / Admin
- Check‑in / walk‑in / checkout flows  
- CRUD management of hotels, chains, rooms, clients, employees  
- "Required Views":  
  - Rooms available per city zone  
  - Total room capacity per hotel

---

## Tech Stack

- **Database:** Supabase (PostgreSQL)  
- **Backend:** Go (Gorilla Mux, pq, JWT) on GCP Cloud Run  
- **Frontend:** HTML, CSS, vanilla JavaScript  
- **Tools:** GitHub, VS Code, Supabase CLI  

---

## Prerequisites

- **VS Code** (or any code editor)  
- **A modern browser** (Chrome, Firefox, Edge, Safari)  
- **Go** (v1.18+) if you need to run the backend locally  
- **Python 3** (optional, for serving static files)  
- **Supabase account** for database access user:sunflowerbookingtest@gmail.com pass:SunflowerBooking1234!

---

## Getting Started

1. **Clone the repository**  
   ```bash
   git clone https://github.com/ARelaxedScholar/DatabaseSQL-Project.git
   cd DatabaseSQL-Project/UI
   ```

2. **Serve the frontend**  
   You can use any static file server. For example, with Python:
   ```bash
   python3 -m http.server 8000
   ```
   Then open your browser at `http://localhost:8000/index.html`.

3. **Backend**  
   The backend is deployed on GCP Cloud Run; no local backend run is required.  
   If you do need to run locally:
   ```bash
   cd backend
   go run .main.go
   ```

---

## Project Structure

---

## API Reference

### Public Endpoints

| Method | Path                               | Description                              |
|--------|------------------------------------|------------------------------------------|
| GET    | `/hotelchains`                     | List hotel chains (id + name)            |
| GET    | `/hotels`                          | List hotels (id + name)                  |
| GET    | `/roomtypes`                       | List room types (id + name)              |
| GET    | `/search/zones/rooms`              | Rooms available per city zone            |
| GET    | `/search/hotels/{hotelID}/room-count` | Total rooms for a specific hotel      |

### Client (Authentication Required)

| Method | Path                                      | Description                          |
|--------|-------------------------------------------|--------------------------------------|
| POST   | `/clients/login`                          | Request magic login link             |
| POST   | `/clients/register`                       | Create a new client                  |
| GET    | `/clients/{id}`                           | Get client profile                   |
| PUT    | `/clients/{id}`                           | Update client profile                |
| GET    | `/clients/{id}/reservations`              | List client reservations             |
| POST   | `/clients/{id}/reservations`              | Create a reservation                 |
| DELETE | `/clients/{id}/reservations/{reservationId}` | Cancel a reservation             |

### Employee / Admin (Authentication Required)

- Similar CRUD endpoints for employees, hotels, rooms, etc.

---

## Testing

### Backend
```bash
cd backend
go test ./...
```

### Frontend
Open your browser console → perform searches & reservations → verify in Network tab.

---

## Deployment

- **Backend:** GCP Cloud Run (Docker)  
- **Database:** Supabase  
- **Frontend:** static host (eventually GitHub Pages)

---

## Contributing

1. Fork the repository  
2. Create a feature branch (`git checkout -b feat/your-feature`)  
3. Commit your changes (`git commit -m "feat: ..."`)  
4. Push to your fork (`git push origin feat/your-feature`)  
5. Open a Pull Request

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

