FROM golang:1.26-alpine AS builder

ENV GOWORK=off CGO_ENABLED=0
WORKDIR /src

COPY gen/go.mod gen/go.sum ./gen/
COPY services/ticket-service/go.mod ./services/ticket-service/
WORKDIR /src/services/ticket-service
RUN go mod download

WORKDIR /src
COPY gen/ ./gen/
COPY services/ticket-service/ ./services/ticket-service/
WORKDIR /src/services/ticket-service
RUN go build -trimpath -ldflags="-s -w" -o /out/ticket-service ./cmd/server

FROM alpine:3.22
RUN apk --no-cache add ca-certificates tzdata && adduser -D -u 10001 app
USER app
WORKDIR /app
COPY --from=builder /out/ticket-service ./ticket-service
EXPOSE 50052
ENTRYPOINT ["./ticket-service"]
