# Forum

A web forum built with Go and SQLite for the Zone01 Athens project.

## Team

- Christos Paloglou
- Giorgos Salaounis
- Giannis Athanasopoulos

## Current Status

### Done

- HTTP server with SQLite connection and database initialization
- User registration (email, username, password) with duplicate validation
- Password hashing with bcrypt and UUID session tokens
- User login, logout, and session cookies (24h expiration)
- Posts with categories (create, list, single post view)
- Comments (create and display)
- Likes and dislikes on posts and comments
- Post filtering by category, created posts, and liked posts
- Responsive UI with light/dark theme and static assets
- Docker setup (Dockerfile, docker-compose.yml)
- Database schema for all forum entities
- Unit tests for auth, posts, and filters

### In Progress / TODO

- Unit tests (broader coverage)

## Requirements

- Go 1.25.3 or compatible
- CGO enabled (required for SQLite)
- SQLite CLI (`sqlite3`) — optional, useful for inspecting the database
- Docker Desktop — optional, required only for `docker compose up --build`

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

## Docker

Build and run with Docker Compose (from project root):

```bash
docker compose up --build
```

Open [http://localhost:8080](http://localhost:8080).

Stop the container:

```bash
docker compose down
```

The SQLite database is created inside the container at `/app/forum.db`.

## Routes

| Method | Path | Description | Auth required |
|--------|------|-------------|---------------|
| GET | `/` | Home page with post list and filters | No |
| GET | `/?category=Go` | Filter posts by category | No |
| GET | `/?filter=mine` | Filter posts created by logged-in user | Yes |
| GET | `/?filter=liked` | Filter posts liked by logged-in user | Yes |
| GET/POST | `/register` | Registration form / create user | No |
| GET/POST | `/login` | Login form / authenticate | No |
| GET | `/logout` | End session and logout | No |
| GET/POST | `/create-post` | Create a new post | Yes |
| GET | `/post?id=1` | View a single post with comments | No |
| POST | `/comment/create` | Add a comment to a post | Yes |
| POST | `/like` | Like a post | Yes |
| POST | `/dislike` | Dislike a post | Yes |
| POST | `/comment/like` | Like a comment | Yes |
| POST | `/comment/dislike` | Dislike a comment | Yes |

## Database

SQLite database file: `forum.db` (created automatically on first run, local only).

Schema is defined in `internal/database/schema.sql`. On startup, `RunMigrations()` in `internal/database/migrations.go` reads this file and creates the tables.

Useful commands:

```bash
sqlite3 forum.db ".tables"
sqlite3 forum.db "SELECT username, email FROM users;"
sqlite3 forum.db "SELECT id, title FROM posts;"
sqlite3 forum.db "SELECT p.title, c.name FROM posts p JOIN post_categories pc ON p.id = pc.post_id JOIN categories c ON c.id = pc.category_id;"
```

## Project Structure

```
forum/
├── cmd/server/          # Application entry point
├── Dockerfile         # Container image build
├── docker-compose.yml # Run forum with Docker Compose
├── internal/
│   ├── auth/            # Cookie and session helpers
│   ├── database/        # SQLite connection and migrations
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   └── services/        # Database business logic
├── templates/           # HTML templates
├── static/              # CSS and JavaScript assets
├── WORKFLOW.md          # Team git workflow and development phases
└── README.md
```

## Dependencies

- [go-sqlite3](https://github.com/mattn/go-sqlite3) — SQLite driver
- [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — Password hashing
- [google/uuid](https://github.com/google/uuid) — Session tokens

## Git Workflow

See [WORKFLOW.md](WORKFLOW.md) for branch naming, pull request process, and phase ownership.

**Do not commit:** `forum.db` (local database), build binaries, or personal notes.
