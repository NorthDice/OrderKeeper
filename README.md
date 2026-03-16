# OrderKeeper

A REST API service for order management with JWT authentication, Redis caching, and observability stack.

## Tech Stack

- **Go 1.24** — language
- **Gin** — HTTP framework
- **PostgreSQL 15** — primary database
- **Redis 7** — caching layer
- **JWT** — authentication
- **Prometheus + Grafana + Loki + Promtail** — metrics and logging
- **Docker / Docker Compose** — containerization

## Getting Started

```bash
docker compose up --build
```

The application will be available at `http://localhost:8080`.

## API Reference

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/auth/sign-up` | Register a new user |
| `POST` | `/auth/sign-in` | Sign in and receive a JWT token |

**Sign Up:**
```json
{
  "username": "john",
  "email": "john@example.com",
  "password": "secret123"
}
```

**Sign In:**
```json
{
  "username": "john",
  "password": "secret123"
}
```

The response includes a `token` to be passed in subsequent requests:
```
Authorization: Bearer <token>
```

### Orders (require authentication)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/order/` | Create an order |
| `GET` | `/order/` | Get all orders |
| `GET` | `/order/:id` | Get order by ID |
| `PUT` | `/order/:id` | Update order |
| `DELETE` | `/order/:id` | Delete order |

**Order statuses:**
- `pending`
- `confirmed`
- `paid`
- `shipped`
- `delivered`
- `cancelled`

### Utility

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/metrics` | Prometheus metrics |

## Monitoring

| Service | URL |
|---------|-----|
| Grafana | `http://localhost:3000` (admin / admin) |
| Prometheus | `http://localhost:9090` |
| Loki | `http://localhost:3100` |

## Project Structure

```
.
├── cmd/            # Entry point
├── configs/        # Configs for Prometheus, Loki, Grafana, Promtail
├── internal/
│   ├── config/     # App configuration
│   ├── handler/    # HTTP handlers and routes
│   ├── models/     # Data models
│   ├── repository/ # Database and cache layer
│   └── service/    # Business logic
├── migrations/     # SQL migrations
├── server/         # HTTP server
├── Dockerfile
└── docker-compose.yaml
```
