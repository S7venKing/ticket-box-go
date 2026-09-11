# Identity Service API

This documentation describes the public HTTP endpoints exposed by the API gateway for the identity service. These endpoints are intended to be consumed by clients such as Postman, frontend apps, or mobile apps.

Base URL:

- Local development: http://localhost:8081
- Docker service port: 8081 -> 8080 (host mapping)

## Authentication flow

The gateway issues JWT access tokens after successful login. Token is sent in the Authorization header:

```http
Authorization: Bearer <access_token>
```

## Endpoints

### 1) Health check

GET /health

Description: Returns health status of the API gateway.

Request:

```http
GET /health
```

Response example:

```http
200 OK
ok
```

---

### 2) Register a new user

POST /register

Description: Creates a new user account.

Headers:

```http
Content-Type: application/json
```

Request body:

```json
{
  "email": "user@example.com",
  "password": "12345678",
  "full_name": "John Doe",
  "phone": "0900000001"
}
```

Success response:

```http
201 Created
```

```json
{
  "user": {
    "id": "535a21e2-c71d-4e8f-b8d4-c9fc3dc339af",
    "email": "user@example.com",
    "full_name": "John Doe",
    "phone": "0900000001",
    "is_active": true,
    "created_at": {
      "seconds": 1789058481,
      "nanos": 402516447
    },
    "updated_at": {
      "seconds": 1789058481,
      "nanos": 402516447
    }
  }
}
```

Possible errors:

- 400 Bad Request: invalid JSON or validation failure
- 409 or 400 depending on underlying identity service response: duplicate email

---

### 3) Login

POST /login

Description: Authenticates a user and returns a JWT token plus the identity token.

Headers:

```http
Content-Type: application/json
```

Request body:

```json
{
  "email": "user@example.com",
  "password": "12345678"
}
```

Success response:

```http
200 OK
```

```json
{
  "user": {
    "id": "535a21e2-c71d-4e8f-b8d4-c9fc3dc339af",
    "email": "user@example.com",
    "full_name": "John Doe",
    "phone": "0900000001",
    "is_active": true,
    "created_at": {
      "seconds": 1789058481,
      "nanos": 402516000
    },
    "updated_at": {
      "seconds": 1789058481,
      "nanos": 402516000
    }
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "identity_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Possible errors:

- 401 Unauthorized: invalid credentials

---

### 4) Get current user profile

GET /me

Description: Returns the authenticated user based on the JWT token in the Authorization header.

Headers:

```http
Authorization: Bearer <access_token>
```

Success response:

```http
200 OK
```

```json
{
  "user_id": "535a21e2-c71d-4e8f-b8d4-c9fc3dc339af",
  "user_email": "user@example.com",
  "user": {
    "id": "535a21e2-c71d-4e8f-b8d4-c9fc3dc339af",
    "email": "user@example.com",
    "full_name": "John Doe",
    "phone": "0900000001",
    "is_active": true,
    "created_at": {
      "seconds": 1789058481,
      "nanos": 402516000
    },
    "updated_at": {
      "seconds": 1789058481,
      "nanos": 402516000
    }
  }
}
```

Possible errors:

- 401 Unauthorized: missing or invalid token
- 404 Not Found: user not found

---

### 5) Logout

POST /logout

Description: Logs the current user out from the gateway. This is a client-side logout pattern for JWT-based authentication.

Headers:

```http
Authorization: Bearer <access_token>
```

Success response:

```http
200 OK
```

```json
{
  "message": "logged out"
}
```

Possible errors:

- 401 Unauthorized: invalid or expired token

---

## Postman import guide

1. Open Postman.
2. Click Import.
3. Choose the file `docs/postman/identity-service-api.postman_collection.json`.
4. The collection will be loaded with:
   - baseUrl variable
   - preconfigured auth token variable
   - examples for register/login/me/logout
5. After login, Postman script automatically stores the token in the collection variable `accessToken`.

## Postman collection variables

- `baseUrl`: default `http://localhost:8081`
- `accessToken`: receives the JWT from the login call

## Example login script for Postman

Add this in the Tests tab of the Login request:

```javascript
var jsonData = pm.response.json();
if (jsonData.access_token) {
  pm.collectionVariables.set("accessToken", jsonData.access_token);
  pm.collectionVariables.set("userId", jsonData.user.id);
}
```

Then use the token in later requests:

```http
Authorization: Bearer {{accessToken}}
```

## Notes

- The gateway is the public API layer.
- The identity service itself is a gRPC service and is not directly called from Postman.
- Clients should call the gateway endpoints instead of the gRPC server directly.
