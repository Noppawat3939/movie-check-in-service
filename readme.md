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

---

### Project Structure

```
root/
├── cmd/
│   └── api/
│       └── main.go                        # entrypoint, wire everything together
│
├── internal/
│   │
│   ├── domain/                            # Entity
│   │
│   ├── usecase/                           # Business
│   │   └── domain_usecase.go              # interface + logic + test
│   │
│   ├── infra/
│   │   ├── postgresl/
│   │   │   ├── db.go                      # init connection pool
│   │   │   └── domain_repo.go
│   │   └── redis/
│   │       ├── client.go                  # init Redis client
│   │       ├── cache.go                   # generic cache
│   │       └── lock.go                    # distributed lock
│   │
│   └── delivery/
│       └── http/
│           ├── router.go                  # register all routes
│           ├── response/
│           │   └── response.go            # success / error helpers
│           └── handler/
│               └── domain_handler.go
│
├── migrations/
│   ├── seeds/
│   │   └── seed.sql                       # initial data
│   ├── 000001_do_something.up.sql
│   └── 000001_do_something.down.sql
│
├── docker-compose.yml                     # app + postgres + redis
├── Dockerfile
├── .env.example
├── Makefile
└── go.mod
```

---

### API

Base URL: `/api/v1`

| Method | Path                               | Description                             |
| ------ | ---------------------------------- | --------------------------------------- |
| GET    | /health                            | Health check                            |
| GET    | /api/v1/movies                     | Get all movies                          |
| GET    | /api/v1/movies/:id                 | Get movie by ID                         |
| GET    | /api/v1/showtime/:showtimeID/seats | Get seats availability for a showtime   |
| POST   | /api/v1/reservation                | Create reservation                      |
| GET    | /api/v1/reservation/:showtimeID    | List reservations by showtime           |
| PATCH  | /api/v1/reservation/:id/seat       | Change seat for an existing reservation |

---

### Getting Started

**1. Copy environment file**

```bash
cp .env.example .env
```

**2. Start services**

```bash
make docker-up
```

**3. Run migrations**

```bash
make migrate-up-one
```

**4. Seed initial data**

```bash
make seed
```

---

### Makefile

```bash
make docker-up        # start app + postgres + redis
make docker-down      # stop all services
make run              # run app locally (without docker)
make tidy             # go mod tidy
make migrate-up-one   # run next migration
make migrate-down-one # rollback last migration
make migrate-version  # show current migration version
make seed             # insert initial data
```
