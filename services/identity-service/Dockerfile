# Build from the REPO ROOT (the service depends on the sibling module gen/):
#   docker build -f services/identity-service/Dockerfile -t identity-service .
#
# GOWORK=off: the image contains only gen/ and this service, not the whole
# workspace; the `replace ../../gen` in go.mod resolves gen locally.

# =========================
# Build stage
# =========================
FROM golang:1.26-alpine AS builder

ENV GOWORK=off CGO_ENABLED=0

WORKDIR /src

# Module files first so `go mod download` is cached across source changes.
COPY gen/go.mod gen/go.sum ./gen/
COPY services/identity-service/go.mod services/identity-service/go.sum ./services/identity-service/

WORKDIR /src/services/identity-service
RUN go mod download

WORKDIR /src
COPY gen/ ./gen/
COPY services/identity-service/ ./services/identity-service/

WORKDIR /src/services/identity-service
RUN go build -trimpath -ldflags="-s -w" -o /out/identity-service ./cmd/server

# =========================
# Runtime stage
# =========================
FROM alpine:3.22

RUN apk --no-cache add ca-certificates tzdata \
    && adduser -D -u 10001 app

USER app
WORKDIR /app

COPY --from=builder /out/identity-service ./identity-service

EXPOSE 50051

ENTRYPOINT ["./identity-service"]
