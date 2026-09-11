# Ticket Box Go

Ticket Box Go is a Go microservice project for identity management, organizer
management, event approval, and the public HTTP API gateway.

## Architecture

```text
Client
  |
  v
api-gateway (HTTP :8081)
  |                         |
  v                         v
identity-service (gRPC)   ticket-service (gRPC)
  |                         |
  +----------- MySQL -------+
```

The repository uses a layered architecture in each service:

```text
domain -> application -> infrastructure -> interfaces
```

## Services

### `identity-service`

Owns user accounts, password hashing, authentication, profile changes, account
activation, and account roles.

Supported roles:

- `user`: normal public account
- `organizer`: can create and submit events belonging to its organizer profile
- `admin`: can create organizer accounts and approve events

The service exposes an internal gRPC API defined in
[`proto/identity/v1/identity.proto`](./proto/identity/v1/identity.proto).

### `ticket-service`

Owns organizer profiles and event business rules. It provides:

- organizer create, read, list, and update operations
- event create, read, list, and approval operations
- MySQL persistence in the `ticket_db` schema
- event lifecycle: `draft -> pending_approval -> published`

The gRPC contract is defined in
[`proto/ticket/v1/ticket.proto`](./proto/ticket/v1/ticket.proto).

### `api-gateway`

The gateway is the public HTTP entry point. It validates JWT access tokens,
enforces roles, and calls the identity/ticket gRPC services.

## Main HTTP API

The gateway is exposed on `http://localhost:8081` by Docker Compose.

### Authentication

| Method | Endpoint | Access |
|---|---|---|
| `GET` | `/health` | Public |
| `POST` | `/register` | Public; creates a `user` |
| `POST` | `/login` | Public |
| `POST` | `/logout` | Bearer token |
| `GET` | `/me` | Bearer token |

Example registration:

```bash
curl -X POST http://localhost:8081/register ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"user@example.com\",\"password\":\"password123\",\"full_name\":\"User One\"}"
```

Example login:

```bash
curl -X POST http://localhost:8081/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"user@example.com\",\"password\":\"password123\"}"
```

Use the returned `access_token` as:

```text
Authorization: Bearer <access_token>
```

### Admin and organizer management

| Method | Endpoint | Access |
|---|---|---|
| `POST` | `/admin/organizers` | Admin |
| `GET` | `/organizers` | Public |
| `GET` | `/organizers/{id}` | Public |
| `POST` | `/organizers` | Admin |
| `PUT` | `/organizers/{id}` | Admin |

`POST /admin/organizers` creates both an identity account with role
`organizer` and the matching organizer profile:

```json
{
  "name": "Tech Organizer",
  "email": "organizer@example.com",
  "password": "password123",
  "phone": "0900000000",
  "slug": "tech-organizer"
}
```

### Events

| Method | Endpoint | Access |
|---|---|---|
| `GET` | `/events` | Public |
| `GET` | `/events/{id}` | Public |
| `POST` | `/events` | Organizer or admin |
| `POST` | `/events/{id}/submit` | Owning organizer |
| `POST` | `/events/{id}/approve` | Admin |
| `POST` | `/events/{id}/publish` | Admin; legacy direct publish route |

Event workflow:

```text
POST /events
      |
      v
draft
      |
      | organizer submits
      v
pending_approval
      |
      | admin approves
      v
published
```

An organizer can only create or submit events for the organizer profile whose
email matches the authenticated organizer account.

## Local development

### Prerequisites

- Go 1.26+
- Docker Desktop or Docker Engine
- `buf` CLI if protobuf files must be regenerated

### Start the full stack

```bash
docker compose up --build
```

Services:

| Service | Container port | Host port |
|---|---:|---:|
| MySQL | `3306` | `3306` |
| identity-service | `50051` | `50051` |
| ticket-service | `50052` | `50052` |
| api-gateway | `8080` | `8081` |

Stop the stack:

```bash
docker compose down
```

The MySQL volume is persistent. To rerun the initialization migrations from a
clean local database:

```bash
docker compose down -v
docker compose up --build
```

### Run tests

```bash
go test ./services/identity-service/... ./services/ticket-service/... ./services/api-gateway/...
```

### Regenerate protobuf code

```bash
buf generate
```

Generated code is stored under [`gen/`](./gen/) and is shared by all service
modules.

## Configuration

### API gateway

- `PORT`: HTTP port inside the container, default `8080`
- `JWT_SECRET`: signing key for gateway access tokens
- `IDENTITY_SERVICE_ADDR`: identity gRPC address
- `TICKET_SERVICE_ADDR`: ticket gRPC address

### Identity service

- `APP_ENV`
- `APP_PORT`
- `MYSQL_HOST`
- `MYSQL_PORT`
- `MYSQL_DATABASE`
- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `JWT_SECRET`

### Ticket service

- `APP_ENV`
- `APP_PORT`
- `MYSQL_HOST`
- `MYSQL_PORT`
- `MYSQL_DATABASE`
- `MYSQL_USER`
- `MYSQL_PASSWORD`

Development defaults are configured in
[`docker-compose.yaml`](./docker-compose.yaml). Production deployments should
provide secrets through an external secret manager or environment injection.

## Promoting a local admin

Public registration creates a `user`. For local development, promote an
existing account after the identity migration:

```sql
UPDATE identity_db.users
SET role = 'admin'
WHERE email = 'admin@example.com';
```

Login again after changing the role so the new JWT contains the `admin` claim.

## API documentation

- Human-readable identity/gateway API notes:
  [`docs/identity-service-api.md`](./docs/identity-service-api.md)
- Importable Postman collection:
  [`docs/postman/identity-service-api.postman_collection.json`](./docs/postman/identity-service-api.postman_collection.json)

## Repository layout

```text
proto/                         Protobuf source contracts
gen/                           Generated Go protobuf modules
services/identity-service/    Identity domain, application, MySQL, gRPC
services/ticket-service/      Organizer/event domain, application, MySQL, gRPC
services/api-gateway/         HTTP routes, JWT middleware, gRPC clients
docker-compose.yaml            Local multi-service deployment
```
