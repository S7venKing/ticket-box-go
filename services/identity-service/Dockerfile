# =========================
# Build stage
# =========================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o /app/bin/app \
    ./cmd/server


# =========================
# Runtime stage
# =========================
FROM alpine:3.22

WORKDIR /app

# CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin/app ./app

EXPOSE 8080

ENTRYPOINT ["./app"]