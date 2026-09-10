# Build from the REPO ROOT:
#   docker build -f services/api-gateway/Dockerfile -t api-gateway .
#
# GOWORK=off: the image contains only gen/ and this service, not the whole
# workspace; the `replace ../../gen` in go.mod resolves gen locally.

# =========================
# Build stage
# =========================
FROM golang:1.26-alpine AS builder

ENV GOWORK=off CGO_ENABLED=0

WORKDIR /src

COPY gen/go.mod gen/go.sum ./gen/
COPY services/api-gateway/go.mod services/api-gateway/go.sum ./services/api-gateway/

WORKDIR /src/services/api-gateway
RUN go mod download

WORKDIR /src
COPY gen/ ./gen/
COPY services/api-gateway/ ./services/api-gateway/

WORKDIR /src/services/api-gateway
RUN go build -trimpath -ldflags="-s -w" -o /out/api-gateway ./cmd

# =========================
# Runtime stage
# =========================
FROM alpine:3.22

RUN apk --no-cache add ca-certificates tzdata \
    && adduser -D -u 10001 app

USER app
WORKDIR /app

COPY --from=builder /out/api-gateway ./api-gateway

EXPOSE 8080

ENTRYPOINT ["./api-gateway"]
