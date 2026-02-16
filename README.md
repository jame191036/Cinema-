# Cinema Ticket Booking System

A real-time cinema ticket booking system with seat locking, concurrency protection, and admin dashboard.

## Architecture

```
┌──────────────────────────────────────────────────────┐
│                    Nginx (port 3000)                  │
│              Reverse Proxy & Load Balancer            │
├──────────┬────────────────────┬──────────────────────┤
│  /api/*  │       /ws/*        │        /*            │
│          │                    │                       │
▼          ▼                    ▼                       │
┌──────────────────┐    ┌──────────────────┐           │
│  Go Backend      │    │  Nuxt Frontend   │           │
│  (Gin) :8080     │    │  (SPA) :3000     │           │
├──────────────────┤    └──────────────────┘           │
│ • REST API       │                                    │
│ • WebSocket Hub  │                                    │
│ • Booking Service│                                    │
│ • Timeout Worker │                                    │
│ • MQ Producer    │                                    │
└──┬───┬───┬───────┘                                    │
   │   │   │                                            │
   ▼   ▼   ▼                                            │
┌─────┐ ┌─────┐ ┌──────────┐                           │
│Mongo│ │Redis│ │ RabbitMQ │                            │
│ :27k│ │:6379│ │ :5672    │                            │
└─────┘ └─────┘ └────┬─────┘                            │
                     │                                  │
                     ▼                                  │
               ┌──────────┐                             │
               │ Consumer │ (writes audit logs)         │
               └──────────┘                             │
└──────────────────────────────────────────────────────┘
```

### Tech Stack

| Layer      | Technology                                |
| ---------- | ----------------------------------------- |
| Frontend   | Nuxt 3 + Vue 3 + shadcn-vue + Tailwind v4 |
| Backend    | Go + Gin + gorilla/websocket              |
| Database   | MongoDB 7                                 |
| Cache/Lock | Redis 7                                   |
| Queue      | RabbitMQ 3                                |
| Proxy      | Nginx                                     |
| Container  | Docker Compose                            |

## Booking Flow

```
User selects seats → Lock Request
    │
    ▼
┌─ Redis SET NX EX 300 ─────────────────────────┐
│  Try acquire lock for each seat                │
│  Key: lock:showtime:{id}:seat:{code}           │
│  Value: userId                                 │
│  TTL: 300 seconds (5 minutes)                  │
└────────────────────────────────────────────────┘
    │ success                    │ fail
    ▼                            ▼
┌─ MongoDB ──────────┐    Return 409 Conflict
│ Update seat_reserv  │    (seat already locked)
│ → state: LOCKED     │
│ Create booking      │
│ → status: LOCKED    │
└─────────────────────┘
    │
    ▼
WebSocket broadcast SEAT_LOCKED to all clients
    │
    ▼
User clicks "Pay & Confirm"
    │
    ▼
┌─ Mock Payment ─────────────────────────────────┐
│ Create payment record (status: SUCCESS)         │
└─────────────────────────────────────────────────┘
    │
    ▼
┌─ Confirm Booking ──────────────────────────────┐
│ 1. Verify: booking LOCKED + payment SUCCESS     │
│ 2. Verify: lock not expired                     │
│ 3. Atomic update: seat_reservations → BOOKED    │
│    (WHERE state=LOCKED AND locked_by=userId)    │
│ 4. Update booking → BOOKED                      │
│ 5. Release Redis locks                          │
│ 6. Broadcast SEAT_BOOKED via WebSocket          │
│ 7. Publish BookingConfirmed to RabbitMQ         │
└─────────────────────────────────────────────────┘
```

### Timeout Flow

A background worker runs every 5 seconds:

1. Queries `seat_reservations` where `state=LOCKED` AND `lock_expires_at < now`
2. Resets seats to `AVAILABLE` in MongoDB
3. Updates booking status to `EXPIRED`
4. Force-releases Redis locks
5. Broadcasts `SEAT_RELEASED` via WebSocket
6. Creates audit log entries

## Redis Lock Strategy

### Key Design

```
lock:showtime:{showtimeId}:seat:{seatCode} → {userId}
TTL: 300 seconds (5 minutes)
```

### Acquire Lock (Atomic)

```
SET lock:showtime:abc:seat:A1 user123 NX EX 300
```

- `NX` = only set if not exists (prevents double-lock)
- `EX 300` = auto-expire in 5 minutes

### Release Lock (Lua Script)

```lua
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
end
return 0
```

Only the lock owner can release it — prevents accidental unlocking.

### Correctness Guarantees

| Concern           | Solution                                                    |
| ----------------- | ----------------------------------------------------------- |
| Double booking    | MongoDB unique index on `(showtime_id, seat_code)`          |
| Race condition    | Redis `SET NX` atomic lock                                  |
| Redis failure     | MongoDB unique constraint as final safety net               |
| Stale locks       | TTL auto-expire + background worker cleanup                 |
| Lock owner check  | Lua script for atomic compare-and-delete                    |
| Concurrent confirm| MongoDB atomic update with state + owner conditions         |

## Message Queue (RabbitMQ)

### Use Case: Async Audit Logging

When a booking is confirmed:

1. **Producer** (in confirm endpoint) publishes `BookingConfirmed` event to `booking.events` queue
2. **Consumer** (background goroutine) receives the event and writes an audit log to MongoDB

### Event Contract

**Queue**: `booking.events`

```json
{
  "eventId": "uuid",
  "eventType": "BookingConfirmed",
  "occurredAt": "2024-01-01T00:00:00Z",
  "bookingId": "mongo-object-id",
  "userId": "mongo-object-id",
  "showtimeId": "mongo-object-id",
  "seats": ["A1", "A2"]
}
```

## WebSocket (Real-time Seat Map)

### Connection

```
ws://localhost:3000/ws/showtimes/{showtimeId}
```

### Events

| Event           | Direction     | Payload                                    |
| --------------- | ------------- | ------------------------------------------ |
| SYNC_SNAPSHOT   | Server→Client | Full seat state on connect                 |
| SEAT_LOCKED     | Server→Client | `{seatCode, lockedByUserId, lockExpiresAt}` |
| SEAT_RELEASED   | Server→Client | `{seatCode}`                               |
| SEAT_BOOKED     | Server→Client | `{seatCode}`                               |

All clients in the same showtime room see real-time updates.

## API Endpoints

### Authentication

| Method | Path              | Description        | Auth   |
| ------ | ----------------- | ------------------ | ------ |
| POST   | /api/auth/login   | Demo login         | Public |

### Movies & Showtimes

| Method | Path                         | Description           | Auth |
| ------ | ---------------------------- | --------------------- | ---- |
| GET    | /api/movies                  | List all movies       | JWT  |
| GET    | /api/movies/:id              | Get movie details     | JWT  |
| GET    | /api/showtimes?movie_id=     | List showtimes        | JWT  |
| GET    | /api/showtimes/:id/seats     | Get seat map          | JWT  |

### Booking

| Method | Path                             | Description        | Auth |
| ------ | -------------------------------- | ------------------ | ---- |
| POST   | /api/showtimes/:id/seats/lock    | Lock seats         | JWT  |
| POST   | /api/bookings/:id/pay            | Mock payment       | JWT  |
| POST   | /api/bookings/:id/confirm        | Confirm booking    | JWT  |
| POST   | /api/bookings/:id/cancel         | Cancel booking     | JWT  |
| GET    | /api/bookings/:id                | Get booking        | JWT  |

### Admin

| Method | Path                    | Description        | Auth  |
| ------ | ----------------------- | ------------------ | ----- |
| GET    | /api/admin/bookings     | List bookings      | Admin |
| GET    | /api/admin/audit-logs   | List audit logs    | Admin |

## Database Schema

### Key Collections

- **users** — auth provider, email, name, role (USER/ADMIN)
- **movies** — title, duration, rating
- **showtimes** — movie reference, start time, auditorium
- **seatmaps** — seat layout (rows × seats)
- **bookings** — user, showtime, seats, status (LOCKED/BOOKED/CANCELLED/EXPIRED)
- **seat_reservations** — per-seat state (**unique index on showtime_id + seat_code**)
- **payments** — mock payment records
- **audit_logs** — event trail for all booking activities

### Critical Index

```
seat_reservations: unique(showtime_id, seat_code)
```

This is the final line of defense against double booking.

## How to Run

### Prerequisites

- Docker & Docker Compose

### Quick Start

```bash
# Clone the repo
git clone <repo-url>

# Start everything
docker compose up --build
```

### Access Points

| Service           | URL                          |
| ----------------- | ---------------------------- |
| Frontend (UI)     | http://localhost:3000         |
| Backend API       | http://localhost:8080/api     |
| RabbitMQ Console  | http://localhost:15672        |
| MongoDB           | mongodb://localhost:27017     |

### Demo Login

1. Open http://localhost:3000
2. Click **Demo User** or **Demo Admin** to log in
3. Browse movies → select showtime → pick seats → lock → pay → confirm

### Development (without Docker)

```bash
# Terminal 1: Start infrastructure
docker compose up mongo redis rabbitmq

# Terminal 2: Go backend
cd backend
MONGO_URI=mongodb://localhost:27017/cinema \
REDIS_ADDR=localhost:6379 \
RABBITMQ_URL=amqp://guest:guest@localhost:5672/ \
go run ./cmd/server

# Terminal 3: Nuxt frontend
cd frontend
npm install
npm run dev
```

## Project Structure

```
cinema/
├── backend/
│   ├── cmd/server/main.go          # Entry point
│   ├── internal/
│   │   ├── config/config.go        # Environment config
│   │   ├── handlers/               # HTTP handlers
│   │   │   ├── auth.go             # Login
│   │   │   ├── booking.go          # Lock/Pay/Confirm/Cancel
│   │   │   ├── movie.go            # Movie CRUD
│   │   │   ├── showtime.go         # Showtime queries
│   │   │   ├── admin.go            # Admin endpoints
│   │   │   └── seed.go             # Seed data
│   │   ├── middleware/auth.go      # JWT + RBAC
│   │   ├── models/                 # MongoDB models
│   │   ├── mq/rabbitmq.go         # RabbitMQ producer/consumer
│   │   ├── services/               # MongoDB + Redis services
│   │   ├── worker/worker.go        # Timeout cleanup worker
│   │   └── ws/hub.go              # WebSocket hub
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── app/
│   │   ├── pages/                  # Nuxt file-based routing
│   │   │   ├── index.vue           # Movie list
│   │   │   ├── login.vue           # Demo login
│   │   │   ├── showtimes/[movieId].vue
│   │   │   ├── seats/[showtimeId].vue  # Real-time seat map
│   │   │   ├── checkout/[bookingId].vue
│   │   │   └── admin/              # Admin pages
│   │   ├── components/ui/          # shadcn-vue components
│   │   ├── composables/useApi.ts   # Axios with auth
│   │   ├── stores/auth.ts          # Pinia auth store
│   │   ├── layouts/default.vue     # App layout
│   │   └── middleware/             # Route guards
│   ├── nuxt.config.ts
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml
├── .env
└── README.md
```
