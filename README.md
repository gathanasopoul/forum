# Forum

A web forum built with Go and SQLite for the Zone01 Athens project.

## Team

- Christos Paloglou
- Giorgos Salaounis
- Giannis Athanasopoulos

## Current Status

### Done

- HTTP server with SQLite connection and database initialization
- User registration (email, username, password)
- Duplicate email validation
- Password hashing with bcrypt
- User login and logout
- Session management with cookies (UUID token, 24h expiration)
- Database schema for users, sessions, posts, comments, categories, and reactions

### In Progress / TODO

- Posts and categories
- Comments
- Likes and dislikes
- Post filtering (by category, created posts, liked posts)
- Frontend templates and styling
- Docker setup
- Unit tests

## Requirements

- Go 1.25.3 or compatible
- CGO enabled (required for SQLite)

## Setup

Clone the repository and install dependencies:

```bash
git clone https://platform.zone01.gr/git/cpaloglo/forum.git
cd forum
go mod download
```

## Run

From the project root:

```bash
go run cmd/server/main.go
```

The server starts at [http://localhost:8080](http://localhost:8080).

## Routes

| Method | Path       | Description              | Auth required |
|--------|------------|--------------------------|---------------|
| GET    | `/`        | Home page (placeholder)  | No            |
| GET    | `/register`| Registration form        | No            |
| POST   | `/register`| Create a new user        | No            |
| GET    | `/login`   | Login form               | No            |
| POST   | `/login`   | Authenticate user        | No            |
| GET    | `/logout`  | End session and logout   | No            |

## Database

SQLite database file: `forum.db` (created automatically on first run).

Schema is defined in `internal/database/schema.sql`.

Useful commands:

```bash
sqlite3 forum.db ".tables"
sqlite3 forum.db "SELECT username, email FROM users;"
sqlite3 forum.db "SELECT user_id, token, expires_at FROM sessions;"
```

## Project Structure

```
forum/
├── cmd/server/          # Application entry point
├── internal/
│   ├── auth/            # Cookie and session helpers
│   ├── database/        # SQLite connection and migrations
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   └── services/        # Database business logic
├── templates/           # HTML templates
└── WORKFLOW.md          # Team git workflow and development phases
```

## Dependencies

- [go-sqlite3](https://github.com/mattn/go-sqlite3) — SQLite driver
- [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — Password hashing
- [google/uuid](https://github.com/google/uuid) — Session tokens

## Git Workflow

See [WORKFLOW.md](WORKFLOW.md) for branch naming, pull request process, and phase ownership.
