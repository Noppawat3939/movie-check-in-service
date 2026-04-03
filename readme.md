# Movie Check-in Service

A backend service for managing movie showtimes and seat reservations.

This system is designed to support movie check-in and seat booking operations, ensuring that the same seat cannot be reserved more than once for the same movie showtime.

---

### Overview

The purpose of this project is to build a backend system for managing movie check-ins.

Users can:

- View available movies
- Select showtimes for each movie
- Choose available seats
- Reserve a seat for a selected showtime

The system must strictly prevent duplicate seat reservations in the same movie showtime.

For example:

- Movie: Avengers
- Showtime: 19:00
- Seat: A1

Once seat `A1` is reserved for this showtime, no other reservation can be made for the same seat in the same round.

---

### Core

##### 1. Movie Management

The system should support movie information management, such as:

- movie title
- description
- duration
- release date

---

##### 2. Showtime Management

Each movie can have multiple showtimes.

Example:

- 10:00
- 13:00
- 16:00
- 19:00

Each showtime contains:

- movie reference
- show date and time
- total number of seats
- available seats

---

##### 3. Seat Reservation

Users should be able to reserve seats for a selected showtime.

Each reservation includes:

- movie
- showtime
- seat number
- customer information (optional)

---

##### 4. Duplicate Seat Protection (Important)

The most important requirement is:

> The same seat must not be reserved more than once for the same movie and showtime.

This means:

- allowed same movie
- allowed same showtime
- not allowed same seat

Example:

Seat `A1` for showtime `19:00` is already reserved.

Another request trying to reserve:

`Movie A + 19:00 + A1`

must be rejected.

This rule must be guaranteed at the database level to prevent race conditions in real-world concurrent requests.

---

##### 5. Payment

Payment is currently out of scope for this version.

The focus of this project is:

> consistency and prevention of duplicate seat reservations

---

### Stacks

- **Language:** Go
- **Framework:** Gin
- **Database:** PostgreSQL
- **Containerization:** Docker
- **Architecture:** Clean Architecture

### Project structure

movie-checkin-service/
├── cmd/
│ └── api/
│ └── main.go # entrypoint, wire everything together
│
├── internal/
│ │
│ │ ├── movie_usecase.go # interface + logic
│ │ ├── showtime_usecase.go
│ │ └── reservation_usecase.go # core: lock → tx → insert
│ │
│ ├── infrastructure/
│ │ ├── postgres/
│ │ │ ├── db.go # db connection
│ │ │ ├── movie_repo.go
│ │ │ ├── showtime_repo.go
│ │ │ ├── seat_repo.go
│ │ │ └── reservation_repo.go
│ │ └── redis/
│ │ ├── client.go # redis connection
│ │ └── lock.go # SETNX / DEL distributed lock
│ │
│ └── delivery/
│ └── http/
│ ├── router.go # register all routes
│ ├── middleware/
│ │ └── error_handler.go # central error → HTTP response
│ └── handler/
│ ├── movie_handler.go
│ ├── showtime_handler.go
│ └── reservation_handler.go
│
├── migrations/
│ ├── 000001_create_movies.up.sql
│ ├── 000001_create_movies.down.sql
│ ├── 000002_create_showtimes.up.sql
│ ├── 000002_create_showtimes.down.sql
│ ├── 000003_create_seats.up.sql
│ ├── 000003_create_seats.down.sql
│ ├── 000004_create_reservations.up.sql
│ └── 000004_create_reservations.down.sql
│
├── docker-compose.yml # app + postgres + redis
├── Dockerfile
├── .env.example
├── Makefile
└── go.mod
