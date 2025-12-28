# Todo List API

A small JSON API for user accounts and personal todos. It is written in Go using Gin,
stores data in SQLite, and protects routes with JWTs. The code is split into handlers,
services, and repositories to keep HTTP, business rules, and SQL separate.

## Features
- Register and login with bcrypt-hashed passwords
- JWT authentication for protected routes
- CRUD endpoints for todos
- Basic rate limiting per client IP

## Requirements
- Go 1.25+
- SQLite

## Setup
Create a `.env` file with a JWT key:

```
JWT_KEY=your-secret-key
```

## Run
```
go run ./cmd/api
```

The server starts on `http://localhost:8080`.

## How it works
Passwords:
- Passwords are hashed with bcrypt using cost 14.
- Plaintext passwords are never stored.

Authentication:
- JWTs are signed with HMAC SHA-256 (HS256).
- Standard claims include `sub` (user id), `iat`, and `exp`.
- Tokens expire after 24 hours.
- Protected routes require `Authorization: Bearer <token>`.

Database:
- SQLite is used as the backing store (`app.db`).
- Foreign keys are enabled on startup.
- `todos.user_id` references `users.id` and deletes are cascaded.

Rate limiting:
- In-memory limiter, 5 requests per second per client IP.

Project layout:
- `cmd/api` holds the API entry point and route wiring.
- `internal/handlers` contains HTTP handlers.
- `internal/service` contains business logic and validation.
- `internal/repository` contains SQL queries.
- `internal/models` defines request and database models.

## API
All routes are under `/api`.

Public:
- `POST /api/register`
- `POST /api/login`
- `GET /api/ping`

Protected (requires `Authorization: Bearer <token>`):
- `GET /api/users/todos`
- `POST /api/users/todos`
- `PUT /api/users/todos/:todoId`
- `DELETE /api/users/todos/:todoId`

## Notes
- Access tokens are verified on each protected request.
- The database is stored in `app.db`.

[Project link](https://roadmap.sh/projects/todo-list-api)
