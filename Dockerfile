# FROM golang:1.21-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN go build -o main ./cmd/api

# FROM alpine:latest

# RUN apk --no-cache add ca-certificates

# WORKDIR /app

# COPY --from=builder /app/main .
# COPY --from=builder /app/.env .
# COPY scripts/init.sql ./scripts/

# EXPOSE 8080

# CMD ["./main"]


FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install PostgreSQL client untuk migration
RUN apk add --no-cache postgresql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o portfolio-app ./cmd/api

FROM alpine:latest

# Install dependencies untuk migration
RUN apk --no-cache add \
    ca-certificates \
    postgresql-client \
    bash

WORKDIR /app

# Copy application dan scripts
COPY --from=builder /app/portfolio-app .
COPY --from=builder /app/.env.local .
COPY --from=builder /app/scripts ./scripts

# Make migration script executable
RUN chmod +x ./scripts/migrate-supabase.sh

EXPOSE 8080

# Run migration lalu start aplikasi
CMD ["sh", "-c", "./scripts/migrate-supabase.sh && ./portfolio-app"]