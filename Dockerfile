# Build stage — CGO is required for github.com/mattn/go-sqlite3
FROM golang:1.25-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o server cmd/server/main.go

# Run stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/internal/database/schema.sql internal/database/schema.sql
COPY --from=builder /app/templates templates
COPY --from=builder /app/static static

EXPOSE 8080

CMD ["./server"]
