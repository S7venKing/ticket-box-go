# Ticket Box Go

Ticket Box Go is a Go microservice project for managing user accounts,
organizers, events, and event approval workflows.

The project currently contains:

- an identity service for users, passwords, authentication, and roles;
- a ticket service for organizer and event business rules;
- an API gateway for the public HTTP API;
- MySQL persistence shared by the local Docker Compose stack;
- protobuf contracts and generated Go clients/servers.

## 1. Architecture

```text
                    HTTP client
                        |
                        v
              +----------------------+
              |    api-gateway       |
              |      HTTP :8080      |
              +----------+-----------+
                         |
            +------------+------------+
            |                         |
            v                         v
 +----------------------+   +----------------------+
 |   identity-service   |   |    ticket-service    |
 |      gRPC :50051     |   |      gRPC :50052     |
 +----------+-----------+   +----------+-----------+
            |                          |
            +------------+-------------+
                         v
                    MySQL :3306
              +----------------------+
              | identity_db          |
              | ticket_db            |
              +----------------------+
```

The public client should call the API gateway. The two backend services are
internal gRPC services and should not be called directly by browser or mobile
clients.

Each service follows the same layered structure:

```text
domain
  -> application
  -> infrastructure
  -> interfaces
```

Responsibilities:

- `domain`: entities, value rules, state transitions, and repository contracts.
- `application`: commands, queries, handlers, and DTOs.
- `infrastructure`: MySQL repositories and external adapters.
- `interfaces`: gRPC handlers and protocol mappers.
- `api-gateway`: HTTP routes, JWT validation, role checks, and gRPC clients.

## 2. Services

### 2.1 Identity service

Directory: [`services/identity-service/`](./services/identity-service/)

The identity service owns:

- user registration;
- password hashing and password verification;
- login through gRPC;
- user profile updates;
- account activation/deactivation;
- user roles;
- user persistence in `identity_db.users`.

Supported roles:

| Role | Meaning |
|---|---|
| `user` | Normal account created by public registration |
| `organizer` | Account linked to an organizer profile; can create and submit events |
| `admin` | Can create organizer accounts and approve events |

The identity protobuf contract is in
[`proto/identity/v1/identity.proto`](./proto/identity/v1/identity.proto).

### 2.2 Ticket service

Directory: [`services/ticket-service/`](./services/ticket-service/)

The ticket service owns:

- organizer profiles;
- organizer validation and lifecycle;
- event creation and validation;
- event listing and lookup;
- event status transitions;
- organizer/event persistence in `ticket_db`.

The event approval lifecycle is:

```text
draft -> pending_approval -> published
```

Additional terminal states are `cancelled` and `closed`.

The ticket protobuf contract is in
[`proto/ticket/v1/ticket.proto`](./proto/ticket/v1/ticket.proto).

### 2.3 API gateway

Directory: [`services/api-gateway/`](./services/api-gateway/)

The gateway:

- exposes the public HTTP API;
- forwards identity operations to identity-service;
- forwards organizer/event operations to ticket-service;
- validates `Authorization: Bearer <token>`;
- places user ID, email, and role into the request context;
- enforces admin and organizer permissions;
- signs its own access JWT after a successful login.

## 3. Roles and business rules

### Public registration

`POST /register` always creates a `user` role. A client cannot select
`admin` or `organizer` through public registration.

### Admin creates an organizer

`POST /admin/organizers` performs two operations:

1. creates an identity account with role `organizer`;
2. creates the matching organizer profile in ticket-service.

The organizer account email and organizer profile email must match so the
gateway can verify ownership later.

### Organizer creates an event

An organizer can create an event only when the supplied `organizer_id`
belongs to an organizer profile with the same email as the authenticated JWT.
An admin can create an event for any organizer.

New events start in `draft`.

### Organizer submits an event

`POST /events/{id}/submit` changes an event from `draft` to
`pending_approval`. Only the organizer who owns the event can submit it.

### Admin approves an event

`POST /events/{id}/approve` changes an event from `pending_approval` to
`published`. Only an admin can approve an event.

## 4. HTTP API

The gateway is exposed at `http://localhost:8081` by the default Docker Compose
configuration.

### 4.1 Authentication endpoints

| Method | Endpoint | Authorization | Description |
|---|---|---|---|
| `GET` | `/health` | Public | Health response |
| `POST` | `/register` | Public | Create a normal user |
| `POST` | `/login` | Public | Authenticate and receive access token |
| `POST` | `/logout` | Bearer token optional | Validate token and end client session |
| `GET` | `/me` | Bearer token | Get current user |
| `GET` | `/profile` | Bearer token | Alias for current user |
| `GET` | `/users/me` | Bearer token | Alias for current user |

#### `GET /health`

Response:

```text
ok
```

#### `POST /register`

Request:

```json
{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "User One",
  "phone": "0900000001"
}
```

Response `201 Created`:

```json
{
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "full_name": "User One",
    "phone": "0900000001",
    "is_active": true,
    "role": "user"
  }
}
```

#### `POST /login`

Request:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response `200 OK`:

```json
{
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "full_name": "User One",
    "is_active": true,
    "role": "user"
  },
  "access_token": "gateway-jwt",
  "identity_token": "identity-jwt"
}
```

The gateway JWT contains the user ID, email, role, issued-at time, expiry, and
token type `access`.

#### `GET /me`

Request header:

```text
Authorization: Bearer <access_token>
```

Response `200 OK`:

```json
{
  "user_id": "user-uuid",
  "user_email": "user@example.com",
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "full_name": "User One",
    "is_active": true,
    "role": "user"
  }
}
```

### 4.2 Organizer endpoints

| Method | Endpoint | Authorization | Description |
|---|---|---|---|
| `POST` | `/admin/organizers` | Admin | Create organizer account and profile |
| `GET` | `/organizers` | Public | List organizers |
| `GET` | `/organizers/{id}` | Public | Get organizer |
| `POST` | `/organizers` | Admin | Create organizer profile |
| `PUT` | `/organizers/{id}` | Admin | Update organizer profile |

#### `POST /admin/organizers`

Request:

```json
{
  "name": "Tech Organizer",
  "email": "organizer@example.com",
  "password": "password123",
  "phone": "0900000000",
  "slug": "tech-organizer"
}
```

Response `201 Created`:

```json
{
  "account": {
    "id": "user-uuid",
    "email": "organizer@example.com",
    "full_name": "Tech Organizer",
    "role": "organizer",
    "is_active": true
  },
  "organizer": {
    "id": "organizer-uuid",
    "name": "Tech Organizer",
    "email": "organizer@example.com",
    "phone": "0900000000",
    "slug": "tech-organizer",
    "is_active": true
  }
}
```

#### `GET /organizers?offset=0&limit=20`

Response `200 OK`:

```json
{
  "organizers": [
    {
      "id": "organizer-uuid",
      "name": "Tech Organizer",
      "email": "organizer@example.com",
      "slug": "tech-organizer",
      "is_active": true
    }
  ]
}
```

### 4.3 Event endpoints

| Method | Endpoint | Authorization | Description |
|---|---|---|---|
| `GET` | `/events` | Public | List events |
| `GET` | `/events/{id}` | Public | Get event |
| `POST` | `/events` | Organizer/Admin | Create draft event |
| `POST` | `/events/{id}/submit` | Owning organizer | Submit for admin review |
| `POST` | `/events/{id}/approve` | Admin | Approve and publish event |
| `POST` | `/events/{id}/publish` | Admin | Legacy direct publish route |

#### `POST /events`

Request:

```json
{
  "organizer_id": "organizer-uuid",
  "title": "Technology Conference",
  "description": "A technology conference",
  "venue": "Ho Chi Minh City",
  "start_at": "2026-10-10T09:00:00Z",
  "end_at": "2026-10-10T17:00:00Z",
  "capacity": 500
}
```

Response `201 Created`:

```json
{
  "event": {
    "id": "event-uuid",
    "organizer_id": "organizer-uuid",
    "title": "Technology Conference",
    "description": "A technology conference",
    "venue": "Ho Chi Minh City",
    "start_at": "2026-10-10T09:00:00Z",
    "end_at": "2026-10-10T17:00:00Z",
    "capacity": 500,
    "status": "EVENT_STATUS_DRAFT"
  }
}
```

#### `POST /events/{id}/submit`

No request body is required.

Response `200 OK`:

```json
{
  "event": {
    "id": "event-uuid",
    "status": "EVENT_STATUS_PENDING_APPROVAL"
  }
}
```

#### `POST /events/{id}/approve`

No request body is required.

Response `200 OK`:

```json
{
  "event": {
    "id": "event-uuid",
    "status": "EVENT_STATUS_PUBLISHED"
  }
}
```

### 4.4 Common HTTP status codes

| Status | Meaning |
|---:|---|
| `200` | Request succeeded |
| `201` | Resource created |
| `400` | Invalid request or domain validation error |
| `401` | Missing or invalid JWT |
| `403` | Valid JWT but insufficient role/ownership |
| `404` | Resource not found |
| `502` | Gateway could not call a backend service |

## 5. End-to-end local flow

The following flow demonstrates public registration, admin promotion,
organizer creation, event submission, and admin approval.

### Step 1: Start the stack

```bash
docker compose up --build -d
docker compose ps
```

Expected services:

```text
mysql
identity-service
ticket-service
api-gateway
```

### Step 2: Create a normal account

```bash
curl -X POST http://localhost:8081/register ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"admin@example.com\",\"password\":\"password123\",\"full_name\":\"Local Admin\"}"
```

### Step 3: Promote the account to admin for local development

```sql
UPDATE identity_db.users
SET role = 'admin'
WHERE email = 'admin@example.com';
```

Login again after this update:

```bash
curl -X POST http://localhost:8081/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"admin@example.com\",\"password\":\"password123\"}"
```

Copy the returned `access_token` into `ADMIN_TOKEN`.

### Step 4: Admin creates an organizer

```bash
curl -X POST http://localhost:8081/admin/organizers ^
  -H "Authorization: Bearer ADMIN_TOKEN" ^
  -H "Content-Type: application/json" ^
  -d "{\"name\":\"Tech Organizer\",\"email\":\"organizer@example.com\",\"password\":\"password123\",\"phone\":\"0900000000\",\"slug\":\"tech-organizer\"}"
```

Copy the returned organizer account token after logging in and copy the
organizer ID into `ORGANIZER_ID`:

```bash
curl -X POST http://localhost:8081/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"organizer@example.com\",\"password\":\"password123\"}"
```

### Step 5: Organizer creates an event

```bash
curl -X POST http://localhost:8081/events ^
  -H "Authorization: Bearer ORGANIZER_TOKEN" ^
  -H "Content-Type: application/json" ^
  -d "{\"organizer_id\":\"ORGANIZER_ID\",\"title\":\"Technology Conference\",\"description\":\"Conference\",\"venue\":\"HCM\",\"start_at\":\"2026-10-10T09:00:00Z\",\"end_at\":\"2026-10-10T17:00:00Z\",\"capacity\":500}"
```

Copy the returned event ID into `EVENT_ID`.

### Step 6: Organizer submits the event

```bash
curl -X POST http://localhost:8081/events/EVENT_ID/submit ^
  -H "Authorization: Bearer ORGANIZER_TOKEN"
```

### Step 7: Admin approves the event

```bash
curl -X POST http://localhost:8081/events/EVENT_ID/approve ^
  -H "Authorization: Bearer ADMIN_TOKEN"
```

The final event status should be `published`.

## 6. Local development

### Prerequisites

- Go 1.26 or newer;
- Docker Desktop/Docker Engine with Compose;
- `buf` CLI when protobuf source files are changed;
- an available host port `8081` for the gateway, `50051` for identity,
  `50052` for ticket-service, and `3306` for MySQL.

### Start and stop Docker Compose

```bash
docker compose up --build
docker compose down
```

Run in detached mode:

```bash
docker compose up --build -d
docker compose ps
docker compose logs -f api-gateway
```

Reset the local database and rerun initialization migrations:

```bash
docker compose down -v
docker compose up --build
```

The `-v` option removes the local MySQL volume. Do not use it if the local
database contains data that must be preserved.

### Run tests

```bash
go test ./services/identity-service/... ./services/ticket-service/... ./services/api-gateway/...
```

Run a specific service:

```bash
go test ./services/identity-service/...
go test ./services/ticket-service/...
go test ./services/api-gateway/...
```

### Regenerate protobuf files

Source contracts:

- [`proto/identity/v1/identity.proto`](./proto/identity/v1/identity.proto)
- [`proto/ticket/v1/ticket.proto`](./proto/ticket/v1/ticket.proto)

Regenerate generated Go code from the repository root:

```bash
buf generate
```

Generated code is stored under [`gen/`](./gen/). Do not manually edit generated
protobuf files.

## 7. Docker and database configuration

The default Compose stack contains:

| Service | Internal port | Host port | Database |
|---|---:|---:|---|
| `mysql` | `3306` | `3306` | `identity_db`, `ticket_db` |
| `identity-service` | `50051` | `50051` | `identity_db` |
| `ticket-service` | `50052` | `50052` | `ticket_db` |
| `api-gateway` | `8080` | `8081` | None |

Identity migration:

- [`services/identity-service/migrations/001_create_users.up.sql`](./services/identity-service/migrations/001_create_users.up.sql)
- [`services/identity-service/migrations/002_add_user_role.up.sql`](./services/identity-service/migrations/002_add_user_role.up.sql)

Ticket migration:

- [`services/ticket-service/migrations/001_create_ticket_schema.up.sql`](./services/ticket-service/migrations/001_create_ticket_schema.up.sql)

Important: MySQL initialization scripts run only when the data directory is
empty. If the container has already initialized its persistent volume, adding
a new migration file will not execute it automatically. Use
`docker compose down -v` for a clean local database or apply the SQL manually.

## 8. Environment variables

### API gateway

| Variable | Default/local value | Purpose |
|---|---|---|
| `PORT` | `8080` | HTTP port inside the container |
| `JWT_SECRET` | Compose development secret | Gateway JWT signing key |
| `IDENTITY_SERVICE_ADDR` | `identity-service:50051` | Identity gRPC address |
| `TICKET_SERVICE_ADDR` | `ticket-service:50052` | Ticket gRPC address |

### Identity service

| Variable | Default/local value | Purpose |
|---|---|---|
| `APP_ENV` | `development` | Runtime environment |
| `APP_PORT` | `50051` | gRPC port |
| `MYSQL_HOST` | `mysql` in Compose | MySQL host |
| `MYSQL_PORT` | `3306` | MySQL port |
| `MYSQL_DATABASE` | `identity_db` | Identity schema |
| `MYSQL_USER` | `root` locally | MySQL user |
| `MYSQL_PASSWORD` | `root` locally | MySQL password |
| `JWT_SECRET` | Compose development secret | Identity token signing key |

### Ticket service

| Variable | Default/local value | Purpose |
|---|---|---|
| `APP_ENV` | `development` | Runtime environment |
| `APP_PORT` | `50052` | gRPC port |
| `MYSQL_HOST` | `mysql` in Compose | MySQL host |
| `MYSQL_PORT` | `3306` | MySQL port |
| `MYSQL_DATABASE` | `ticket_db` | Ticket schema |
| `MYSQL_USER` | `root` locally | MySQL user |
| `MYSQL_PASSWORD` | `root` locally | MySQL password |

The values in Compose are for local development only. Production deployments
must inject secrets securely and must not use the example JWT secret or root
database credentials.

## 9. Repository layout

```text
.
├── proto/
│   ├── identity/v1/identity.proto
│   └── ticket/v1/ticket.proto
├── gen/
│   ├── identity/v1/
│   └── ticket/v1/
├── services/
│   ├── identity-service/
│   │   ├── cmd/server/
│   │   ├── internal/domain/
│   │   ├── internal/application/
│   │   ├── internal/infrastructure/
│   │   ├── internal/interfaces/grpc/
│   │   └── migrations/
│   ├── ticket-service/
│   │   ├── cmd/server/
│   │   ├── internal/domain/
│   │   ├── internal/application/
│   │   ├── internal/infrastructure/
│   │   ├── internal/interfaces/grpc/
│   │   └── migrations/
│   └── api-gateway/
│       ├── cmd/
│       └── internal/
├── docker-compose.yaml
├── buf.gen.yaml
└── go.work
```

## 10. API documentation and Postman

Existing API documentation:

- [`docs/identity-service-api.md`](./docs/identity-service-api.md)
- [`docs/postman/identity-service-api.postman_collection.json`](./docs/postman/identity-service-api.postman_collection.json)

The Postman collection uses the gateway base URL
`http://localhost:8081`. For organizer/event testing, first obtain an
admin/organizer access token through `/login`, then set the
`Authorization: Bearer <token>` header on protected requests.

## 11. Troubleshooting

### Gateway port is already in use

The default host port is `8081`. Check the process using it or change the host
mapping in `docker-compose.yaml`, for example:

```yaml
ports:
  - "8082:8080"
```

### Ticket service cannot connect to MySQL

Check service status and logs:

```bash
docker compose ps
docker compose logs ticket-service
docker compose logs mysql
```

The ticket service must use:

```text
MYSQL_HOST=mysql
MYSQL_DATABASE=ticket_db
MYSQL_USER=root
MYSQL_PASSWORD=root
```

### Existing MySQL volume does not contain new tables

The initialization scripts run only on the first MySQL startup. Recreate the
local volume:

```bash
docker compose down -v
docker compose up --build
```

### JWT role is not updated

JWT role claims are created at login time. After changing a user's role in the
database, log in again and use the newly returned `access_token`.

## 12. Security notes

- Password hashes are never returned in DTOs or protobuf responses.
- JWT secrets and database credentials in Compose are development defaults.
- Role checks happen at the gateway; ownership checks compare the JWT email
  with the organizer profile email.
- The ticket and identity gRPC ports should be kept private in production.
- Production should use TLS, managed secrets, a non-root database user, and
  proper refresh-token/revocation handling.
