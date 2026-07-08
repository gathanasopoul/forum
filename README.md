# Forum

A web forum built with Go and SQLite for the Zone01 Athens project.

## Team

- Christos Paloglou
- Giorgos Salaounis
- Giannis Athanasopoulos


## Requirements

- Go 1.25.3 or compatible
- CGO enabled (required for SQLite)
- SQLite CLI (`sqlite3`) — optional, useful for inspecting the database
- Docker — optional, for running the app in a container

## Setup

Clone the repository and install dependencies:

```bash
git clone https://platform.zone01.gr/git/cpaloglo/forum.git
cd forum
go mod download
```

## Authentication

The forum supports **email/password** registration and login out of the box. No environment file is required for this flow.

**GitHub OAuth** is optional. To enable it, copy `.env.example` to `.env` and add your GitHub OAuth app credentials (see [Environment Variables](#environment-variables)).

## Environment Variables

Only needed for **GitHub OAuth**. Email/password auth works without a `.env` file.

1. Copy the template: `cp .env.example .env`
2. Set your GitHub OAuth app values in `.env`:

```env
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
```

If `.env` is missing, the server still starts and email/password authentication works normally.

## Run

From the project root:

```bash
go run cmd/server/main.go
```

The server starts at [http://localhost:8080](http://localhost:8080).

## Docker

### Docker Compose (recommended)

Build and run from the project root:

```bash
docker compose up --build
```

Open [http://localhost:8080](http://localhost:8080).

Stop the container:

```bash
docker compose down
```

### Manual build and run

```bash
docker build -t forum .
docker run -d -p 8080:8080 --name forum forum
docker ps -a
```

Stop and remove:

```bash
docker stop forum && docker rm forum
```

The SQLite database is created inside the container at `/app/forum.db`. Without a volume, data is reset when the container is removed.

Inspect the database inside the container:

```bash
docker exec -it forum sqlite3 /app/forum.db "SELECT * FROM users;"
```

## Routes

| Method | Path | Description | Auth required |
|--------|------|-------------|---------------|
| GET | `/` | Home page with post list and filters | No |
| GET | `/?category=Go` | Filter posts by category | No |
| GET | `/?filter=mine` | Filter posts created by logged-in user | Yes |
| GET | `/?filter=liked` | Filter posts liked by logged-in user | Yes |
| GET/POST | `/register` | Registration form / create user | No |
| GET/POST | `/login` | Login form / authenticate with email and password | No |
| GET | `/logout` | End session and logout | No |
| GET | `/auth/github/login` | Redirect to GitHub OAuth (requires `.env`) | No |
| GET | `/auth/github/callback` | GitHub OAuth callback (requires `.env`) | No |
| GET/POST | `/create-post` | Create a new post | Yes |
| GET | `/post?id=1` | View a single post with comments | No |
| POST | `/comment/create` | Add a comment to a post | Yes |
| POST | `/like` | Like a post | Yes |
| POST | `/dislike` | Dislike a post | Yes |
| POST | `/comment/like` | Like a comment | Yes |
| POST | `/comment/dislike` | Dislike a comment | Yes |
| GET | `/static/*` | CSS and static assets | No |

## Database

SQLite database file: `forum.db` (created automatically on first run, local only).

Schema is defined in `internal/database/schema.sql`. On startup, `RunMigrations()` in `internal/database/migrations.go` reads this file and creates the tables.

Useful commands:

```bash
sqlite3 forum.db ".tables"
sqlite3 forum.db "SELECT * FROM users;"
sqlite3 forum.db "SELECT * FROM posts;"
sqlite3 forum.db "SELECT * FROM comments;"
sqlite3 forum.db "SELECT p.title, c.name FROM posts p JOIN post_categories pc ON p.id = pc.post_id JOIN categories c ON c.id = pc.category_id;"
```

## Project Structure

```
forum/
├── cmd/server/          # Application entry point
├── Dockerfile           # Container image build
├── docker-compose.yml   # Run forum with Docker Compose
├── .env.example         # GitHub OAuth template (copy to .env locally)
├── internal/
│   ├── auth/            # Cookie and session helpers
│   ├── database/        # SQLite connection and migrations
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   └── services/        # Database business logic
├── templates/           # HTML templates
├── static/              # CSS and JavaScript assets
└── README.md
```

## Dependencies

- [go-sqlite3](https://github.com/mattn/go-sqlite3) — SQLite driver
- [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — Password hashing
- [google/uuid](https://github.com/google/uuid) — Session tokens
- [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) — GitHub OAuth (optional)
- [godotenv](https://github.com/joho/godotenv) — Load `.env` file locally


