# Ticket Box API Reference

Base URL: `http://localhost:8081`

All protected endpoints use:

```http
Authorization: Bearer <access_token>
```

## Authentication

### `GET /health`

Public health check. Returns `ok`.

### `POST /register`

Creates a normal `user` account.

```json
{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "User One",
  "phone": "0900000001"
}
```

Returns `201` with `{ "user": { ... } }`.

### `POST /login`

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Returns the user, `access_token`, and `identity_token`.

### `POST /logout`

Requires a valid bearer token when the header is supplied. Returns
`{ "message": "logged out" }`.

### `GET /me`

Requires a bearer token and returns the current user.

## Admin and organizer accounts

### `POST /admin/organizers`

Admin only. Creates an identity account with role `organizer` and a matching
organizer profile.

```json
{
  "name": "Tech Organizer",
  "email": "organizer@example.com",
  "password": "password123",
  "phone": "0900000000",
  "slug": "tech-organizer"
}
```

### `GET /organizers?offset=0&limit=20`

Public paginated organizer list.

### `GET /organizers/{id}`

Public organizer detail.

### `POST /organizers`

Admin only. Creates an organizer profile.

### `PUT /organizers/{id}`

Admin only. Updates name, phone, and slug.

## Events

### `GET /events?offset=0&limit=20`

Public event list.

### `GET /events/{id}`

Public event detail.

### `POST /events`

Organizer or admin. Organizers may only use their own organizer profile.

```json
{
  "organizer_id": "organizer-uuid",
  "title": "Technology Conference",
  "description": "A technology conference",
  "venue": "HCM",
  "start_at": "2026-10-10T09:00:00Z",
  "end_at": "2026-10-10T17:00:00Z",
  "capacity": 500
}
```

New events have status `EVENT_STATUS_DRAFT`.

### `POST /events/{id}/submit`

Owning organizer only. Changes status:

```text
draft -> pending_approval
```

### `POST /events/{id}/approve`

Admin only. Changes status:

```text
pending_approval -> published
```

### `POST /events/{id}/publish`

Admin-only legacy direct publish endpoint. New clients should use
`/submit` followed by `/approve`.

## Status codes

- `200`: success
- `201`: resource created
- `400`: invalid body or business validation failure
- `401`: missing/invalid authentication
- `403`: insufficient role or ownership
- `404`: resource not found
- `502`: backend service unavailable

## Business workflow

```text
register -> login
admin creates organizer account
organizer login
organizer creates event (draft)
organizer submits event (pending_approval)
admin approves event (published)
```
