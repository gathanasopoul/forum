# Forum Project Workflow

## Team Rules

* Never push directly to `main`.
* Every feature must be developed in a separate branch.
* Open a Pull Request before merging into `main`.
* Pull latest changes before starting work.
* Keep commits small and descriptive.
* Write comments in English.
* Test features locally before pushing.

---

## Git Workflow

### Update local repository

```bash
git checkout main
git pull origin main
```

### Create feature branch

```bash
git checkout -b feature/feature-name
```

Examples:

```bash
git checkout -b feature/database
git checkout -b feature/authentication
git checkout -b feature/posts
git checkout -b feature/comments
git checkout -b feature/reactions
git checkout -b feature/filtering
git checkout -b feature/docker
```

### Push branch

```bash
git push -u origin feature/feature-name
```

### Merge process

1. Push feature branch.
2. Create Pull Request.
3. Review code.
4. Merge into main.
5. Delete branch.

---

# Development Phases

## Phase 1 - Foundation

Goal:

* HTTP Server
* SQLite Connection
* Database Initialization
* Error Handling

Files:

* cmd/server/main.go
* internal/database/db.go
* internal/database/migrations.go
* migrations/001_init.sql


---

## Phase 2 - Authentication

Goal:

* Register
* Login
* Logout
* Session Cookies

Files:

* internal/auth/*
* internal/models/user.go
* internal/models/session.go


---

## Phase 3 - Posts

Goal:

* Create Post
* View Posts
* View Single Post

Files:

* internal/handlers/post.go
* internal/models/post.go
* internal/services/post_service.go

---

## Phase 4 - Comments

Goal:

* Create Comment
* Display Comments

Files:

* internal/handlers/comment.go
* internal/models/comment.go
* internal/services/comment_service.go


---

## Phase 5 - Reactions

Goal:

* Like Post
* Dislike Post
* Like Comment
* Dislike Comment

Files:

* internal/handlers/like.go
* internal/models/reaction.go
* internal/services/reaction_service.go


---

## Phase 6 - Filters

Goal:

* Filter by Category
* Filter by Created Posts
* Filter by Liked Posts

Files:

* internal/handlers/filter.go


---

## Phase 7 - Frontend

Goal:

* Templates
* Styling
* Navigation

Files:

* templates/*
* static/*

Shared

---

## Phase 8 - Docker

Goal:

* Dockerfile
* docker-compose.yml


---

# Definition of Done

A feature is complete when:

* Code builds successfully.
* No compilation errors.
* Feature works locally.
* Code is pushed to feature branch.
* Pull Request approved.
* Merged into main.
