# Backend Portfolio API

REST API backend for a personal portfolio website, built with Go (Gin framework), PostgreSQL (GORM), and secured with JWT Bearer token authentication (no cookies, no CSRF needed).

## Table of Contents

- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Authentication Flow](#authentication-flow)
- [CSRF Protection](#csrf-protection)
- [Environment Variables](#environment-variables)
- [Getting Started](#getting-started)
- [API Reference](#api-reference)
- [Postman Testing Guide](#postman-testing-guide)
- [Security Features](#security-features)
- [Database Management](#database-management)
- [JWT Bearer Authentication](#jwt-bearer-authentication)
- [Frontend Integration](#frontend-integration)

---

## Architecture

```
┌─────────────┐    HTTPS     ┌───────────────────────────────────────────────┐
│   Frontend   │◄────────────►│               Gin HTTP Server                 │
│  (React/Next)│  Cookies +   │                                               │
│              │  CSRF Header │  ┌──────────────┐  ┌───────────────────────┐  │
└─────────────┘              │  │  Middleware   │  │  Handler (DI struct)  │  │
                              │  │ ┌──────────┐ │  │ ┌──────┐ ┌────────┐  │  │
                              │  │ │  CORS    │ │  │ │ Auth │ │Portfol.│  │  │
                              │  │ ├──────────┤ │  │ ├──────┤ ├────────┤  │  │
                              │  │ │ Security │ │  │ │About │ │ Skills │  │  │
                              │  │ │ Headers  │ │  │ ├──────┤ ├────────┤  │  │
                              │  │ ├──────────┤ │  │ │Quals │ │ Health │  │  │
                              │  │ │  Auth    │ │  │ └──────┘ └────────┘  │  │
                              │  │ ├──────────┤ │  └──────────┬────────────┘  │
                              │  │ │  CSRF    │ │             │ (interfaces)  │
                              │  │ └──────────┘ │  ┌──────────▼────────────┐  │
                              │  └──────────────┘  │  Repository (GORM)    │  │
                              │                     └──────────┬────────────┘  │
                              │                     ┌──────────▼────────────┐  │
                              │                     │  PostgreSQL Database  │  │
                              │                     └───────────────────────┘  │
                              └───────────────────────────────────────────────┘
```

**Tech Stack:**
- **Language:** Go 1.21
- **HTTP Framework:** Gin v1.9.1
- **ORM:** GORM v1.25.5 with PostgreSQL driver
- **JWT:** golang-jwt/jwt/v4 (HMAC-SHA256)
- **Password Hashing:** bcrypt (golang.org/x/crypto)
- **Deployment:** Docker, Google Cloud Run

---

## Project Structure

```
backend-portfolio/
├── cmd/
│   └── api/
│       └── main.go              # Production entrypoint (Cloud Run / Docker)
├── internal/                    # Private application packages (Go convention)
│   ├── auth/
│   │   └── jwt.go               # JWTService — JWT token generation & validation (DI)
│   ├── handler/
│   │   ├── handler.go           # Handler struct with injected dependencies
│   │   ├── about.go             # About CRUD handlers
│   │   ├── auth.go              # Login, Logout, Refresh, GetMe, CreateUser
│   │   ├── health.go            # Health check & info endpoints
│   │   ├── portfolio.go         # Portfolio CRUD with file upload
│   │   ├── qualification.go     # Qualification CRUD handlers
│   │   └── skill.go             # Skill CRUD handlers (category, level, score)
│   ├── middleware/
│   │   ├── auth.go              # JWT Bearer token auth middleware (DI)
│   │   └── security.go          # Security response headers (CSP, X-Frame-Options)
│   ├── repository/
│   │   ├── interfaces.go        # Repository interfaces (DIP)
│   │   ├── about.go             # AboutRepository — GORM implementation
│   │   ├── portfolio.go         # PortfolioRepository — GORM + transactions
│   │   ├── qualification.go     # QualificationRepository — GORM
│   │   ├── skill.go             # SkillRepository — GORM + grouped category
│   │   └── user.go              # UserRepository — GORM implementation
│   └── router/
│       └── router.go            # Shared route definitions (DRY)
├── config/
│   └── config.go                # Environment-based configuration
├── database/
│   └── database.go              # Database connection, retry, migration
├── models/
│   ├── about.go                 # About model
│   ├── portfolio.go             # Portfolio + PortfolioMedia models
│   ├── qualification.go         # Qualification model
│   ├── skill.go                 # Skill model
│   └── user.go                  # User model
├── scripts/
│   ├── deploy.sh                # Cloud Run deploy script
│   ├── init.sql                 # Database initialization
│   ├── migrate-supabase.sh      # Supabase migration
│   └── migrations/              # SQL migration files
├── uploads/                     # File uploads directory
├── main.go                      # Development entrypoint
├── Dockerfile                   # Production Docker build
├── Dockerfile.prod              # Alternative production build
├── Dockerfile.migrator          # Database migration
├── docker-compose.yml           # Local development
├── docker-compose.prod.yml      # Production Docker Compose
├── docker-compose.supabase.yml  # Supabase connection
├── go.mod                       # Go module definition
└── go.sum                       # Dependency checksums
```

### Design Decisions (SOLID)

| Principle | Implementation |
|-----------|---------------|
| **S** — Single Responsibility | Each handler file owns one domain; repositories own DB access; middleware owns cross-cutting concerns |
| **O** — Open/Closed | New resources only require a new repository interface + handler file; no changes to router/middleware |
| **L** — Liskov Substitution | Repository interfaces allow swapping GORM for any implementation (e.g. mock for testing) |
| **I** — Interface Segregation | Five focused repository interfaces instead of one large "store" interface |
| **D** — Dependency Inversion | Handlers depend on repository *interfaces*, not concrete GORM types; all dependencies injected via constructors |

| Decision | Rationale |
|----------|-----------|
| `internal/` packages | Go convention for private application code; prevents external imports |
| `internal/router/` — shared router | Both main.go files call the same `SetupRouter()` — routes defined once (DRY) |
| Constructor injection | No global mutable state; all dependencies explicit in `handler.New()` |
| Repository pattern with interfaces | Decouples handlers from GORM; enables unit testing with mocks |
| Cookie-based JWT | Immune to XSS token theft (HttpOnly cookies) |
| Signed Double Submit Cookie | Stateless CSRF without database storage |

---

## Authentication Flow

### Token Pair System

The API uses a **dual token architecture**:

| Token | Type | Expiry | Storage | Purpose |
|-------|------|--------|---------|---------|
| Access Token | JWT (HS256) | 15 minutes | `access_token` cookie (HttpOnly, Secure, SameSite=Strict) | Authenticate API requests |
| Refresh Token | JWT (HS256) | 7 days | `refresh_token` cookie (HttpOnly, Secure, SameSite=Strict, Path=/api/auth/refresh) | Obtain new access tokens |

### Login Flow

```
Client                              Server
  │                                    │
  ├──POST /api/auth/login──────────────►│
  │  {"username":"admin","password":…}  │
  │                                    │
  │◄─────────Set-Cookie: access_token──┤
  │◄─────────Set-Cookie: refresh_token─┤
  │◄─────────Set-Cookie: csrf_token────┤  (readable by JS)
  │◄─────────Set-Cookie: csrf_token_sig┤  (HttpOnly signature)
  │◄─────────200 OK {user info}────────┤
  │                                    │
```

### Token Refresh Flow

```
Client                              Server
  │                                    │
  ├──POST /api/auth/refresh────────────►│  (refresh_token cookie auto-sent)
  │                                    │
  │  Server validates refresh token,    │
  │  verifies user still exists,        │
  │  generates new token pair           │
  │                                    │
  │◄─────────Set-Cookie: new tokens────┤  (token rotation)
  │◄─────────Set-Cookie: new CSRF──────┤
  │◄─────────200 OK───────────────────┤
  │                                    │
```

### Request Authentication

```
Client                              Server
  │                                    │
  ├──PUT /api/admin/about/1────────────►│
  │  Cookie: access_token=eyJ…         │
  │  X-XSRF-TOKEN: abc123…            │
  │                                    │
  │  AuthMiddleware:                    │
  │    1. Read access_token cookie      │
  │    2. Validate JWT signature        │
  │    3. Check role == "admin"         │
  │  CSRFMiddleware:                    │
  │    4. Compare header vs cookie      │
  │    5. Verify HMAC signature         │
  │                                    │
  │◄─────────200 OK───────────────────┤
  │                                    │
```

### Cookie Properties

| Cookie | HttpOnly | Secure | SameSite | Path | MaxAge |
|--------|----------|--------|----------|------|--------|
| `access_token` | ✅ Yes | ✅ Yes* | Strict | `/` | 15 min |
| `refresh_token` | ✅ Yes | ✅ Yes* | Strict | `/api/auth/refresh` | 7 days |
| `csrf_token` | ❌ No (JS reads it) | ✅ Yes* | Strict | `/` | 15 min |
| `csrf_token_sig` | ✅ Yes | ✅ Yes* | Strict | `/` | 15 min |

\* `Secure` is configurable via `COOKIE_SECURE` env var. Set to `false` for local HTTP development.

---

## CSRF Protection

Uses the **Signed Double Submit Cookie** pattern (OWASP recommended for stateless architectures).

### How It Works

1. **On login/refresh:** Server generates a random token, HMAC-signs it with `JWT_SECRET`, stores the token in a readable cookie (`csrf_token`) and the signature in an HttpOnly cookie (`csrf_token_sig`).
2. **On mutation requests (POST/PUT/DELETE):** Frontend reads `csrf_token` cookie value and sends it in the `X-XSRF-TOKEN` header.
3. **Server validates:** Compares the header value against the cookie, then verifies the HMAC signature matches — ensuring the token was server-issued.

### Why No Database?

The HMAC signature proves the token was created by this server (only this server knows `JWT_SECRET`). An attacker cannot forge a valid token+signature pair without the secret key.

### Bypasses

CSRF validation is **skipped** for:
- `GET`, `HEAD`, `OPTIONS` methods (safe/idempotent)
- Requests using `Authorization: Bearer` header (API clients, not vulnerable to CSRF)

---

## Environment Variables

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `DB_HOST` | `localhost` | ✅ | PostgreSQL host |
| `DB_PORT` | `5432` | ✅ | PostgreSQL port |
| `DB_USER` | `postgres` | ✅ | PostgreSQL username |
| `DB_PASSWORD` | *(empty)* | ✅ | PostgreSQL password |
| `DB_NAME` | `portfolio-db` | ✅ | PostgreSQL database name |
| `DB_INSTANCE_NAME` | *(empty)* | ❌ | Cloud SQL instance (Cloud Run only) |
| `JWT_SECRET` | `your-secret-key` | ✅ | HMAC key for access tokens + CSRF signing |
| `JWT_REFRESH_SECRET` | `your-refresh-secret-key` | ✅ | HMAC key for refresh tokens |
| `PORT` | `8080` | ❌ | Server listen port |
| `COOKIE_DOMAIN` | *(empty)* | ❌ | Cookie domain (e.g., `.adirdk.com`) |
| `COOKIE_SECURE` | `true` | ❌ | Set to `false` for local HTTP dev |
| `COOKIE_SAMESITE` | `Strict` | ❌ | `Strict`, `Lax`, or `None` |
| `ALLOWED_ORIGINS` | *(hardcoded list)* | ❌ | Comma-separated CORS origins override |
| `RESET_SECRET` | *(empty)* | ❌ | Secret for admin password reset endpoint |
| `GIN_MODE` | `debug` | ❌ | Set to `release` for production |
| `K_SERVICE` | *(auto)* | ❌ | Cloud Run service name (auto-set by Cloud Run) |

### Example `.env` for local development

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=portfolio_user
DB_PASSWORD=portfolio_password
DB_NAME=portfolio_db
JWT_SECRET=my-super-secret-jwt-key-at-least-32-chars
JWT_REFRESH_SECRET=my-super-secret-refresh-key-at-least-32-chars
COOKIE_SECURE=false
COOKIE_SAMESITE=Lax
COOKIE_DOMAIN=localhost
RESET_SECRET=my-reset-secret
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

> ⚠️ **IMPORTANT:** In production, use strong random secrets (min 32 characters). Never commit secrets to version control.

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Docker & Docker Compose (optional)

### Option 1: Docker Compose (Recommended)

```bash
# Clone and start
git clone <repo-url>
cd backend-portfolio
docker compose up -d

# API available at http://localhost:8080
# Database at localhost:5432
```

### Option 2: Local Go

```bash
# Start PostgreSQL, then:
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=your_pass DB_NAME=portfolio_db
export JWT_SECRET=your-secret JWT_REFRESH_SECRET=your-refresh-secret
export COOKIE_SECURE=false

go run main.go
# or
go run ./cmd/api
```

### Create First Admin User

After starting the server, create an admin user via the database or use the API:

```sql
-- Option 1: Direct SQL (hash a password with bcrypt)
INSERT INTO users (username, password, role)
VALUES ('admin', '$2a$10$...hashed_password...', 'admin');
```

```bash
# Option 2: Use the create-user endpoint (requires existing admin auth)
# First login, then:
curl -X POST http://localhost:8080/api/admin/create-user \
  -H "Content-Type: application/json" \
  -H "X-XSRF-TOKEN: <csrf_token>" \
  -b "access_token=<token>; csrf_token=<token>" \
  -d '{"username":"admin","password":"securepassword123","role":"admin"}'
```

---

## API Reference

### Public Endpoints (No Auth)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | API info + version |
| `GET` | `/health` | Health check (DB ping) |
| `GET` | `/api/about` | Get about information |
| `GET` | `/api/portfolio` | List all portfolio items |
| `GET` | `/api/portfolio/:id` | Get single portfolio item |
| `GET` | `/api/skills` | List all skills |
| `GET` | `/api/skills/:id` | Get single skill |
| `GET` | `/api/qualifications` | List all qualifications |
| `GET` | `/api/qualifications/:id` | Get single qualification |

### Auth Endpoints

| Method | Path | Auth | CSRF | Description |
|--------|------|------|------|-------------|
| `POST` | `/api/auth/login` | ❌ | ❌ | Login (sets cookies) |
| `POST` | `/api/auth/logout` | ❌ | ❌ | Logout (clears cookies) |
| `POST` | `/api/auth/refresh` | ❌ | ❌ | Refresh tokens (uses refresh cookie) |
| `GET` | `/api/auth/me` | ✅ Cookie | ❌ | Get current user info |
| `POST` | `/api/login` | ❌ | ❌ | Login (backward compatibility) |

### Admin Endpoints (Cookie Auth + CSRF Required)

All admin endpoints require:
1. Valid `access_token` cookie (from login)
2. `X-XSRF-TOKEN` header matching `csrf_token` cookie

#### User Management

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/admin/create-user` | Create new user (min 8 char password) |
| `POST` | `/api/admin/reset-admin` | Reset admin password (requires `X-Reset-Secret` header) |

#### About

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/admin/about` | Create or update about entry |
| `PUT` | `/api/admin/about/:id` | Update specific about entry |

#### Portfolio

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/admin/portfolio` | Create portfolio item (multipart/form-data for files) |
| `PUT` | `/api/admin/portfolio/:id` | Update portfolio item |
| `DELETE` | `/api/admin/portfolio/:id` | Delete portfolio item |
| `DELETE` | `/api/admin/portfolio-media/:portfolio_id/:media_id` | Delete specific media from portfolio |

#### Skills

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/admin/skills` | Create skill (score: 0-100) |
| `PUT` | `/api/admin/skills/:id` | Update skill |
| `DELETE` | `/api/admin/skills/:id` | Delete skill |

#### Qualifications

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/admin/qualifications` | Create qualification |
| `PUT` | `/api/admin/qualifications/:id` | Update qualification |
| `DELETE` | `/api/admin/qualifications/:id` | Delete qualification |

---

### Request/Response Examples

#### Login

```bash
POST /api/auth/login
Content-Type: application/json

{"username": "admin", "password": "your_password"}
```

Response:
```json
{
  "message": "Login successful",
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

Response cookies set:
- `access_token` (HttpOnly)
- `refresh_token` (HttpOnly, path=/api/auth/refresh)
- `csrf_token` (readable by JS)
- `csrf_token_sig` (HttpOnly)

#### Create Portfolio (with file upload)

```bash
POST /api/admin/portfolio
Content-Type: multipart/form-data
Cookie: access_token=eyJ...
X-XSRF-TOKEN: <value_from_csrf_token_cookie>

title=My Project
description=Project description
url=https://example.com
files=@image1.png
files=@image2.jpg
```

#### Update Skill

```bash
PUT /api/admin/skills/1
Content-Type: application/json
Cookie: access_token=eyJ...
X-XSRF-TOKEN: <value_from_csrf_token_cookie>

{
  "name": "Go",
  "category": "Backend",
  "level": "Advanced",
  "score": 90
}
```

---

## Postman Testing Guide

### Setup

1. Open Postman and create a new Collection.
2. Go to **Settings** (gear icon) and ensure **"Automatically follow redirects"** is enabled.
3. Cookies are handled automatically by Postman when the domain matches.

### Step 1: Login

```
POST http://localhost:8080/api/auth/login
Body (JSON):
{
  "username": "admin",
  "password": "your_password"
}
```

After successful login, check the **Cookies** tab in Postman — you should see:
- `access_token`
- `refresh_token`
- `csrf_token`
- `csrf_token_sig`

### Step 2: Copy CSRF Token

1. Click the **Cookies** link under the Send button (or go to **Manage Cookies**).
2. Find the `csrf_token` cookie value (e.g., `a1b2c3d4e5f6...`).
3. Copy this value.

### Step 3: Make Admin Requests

For any `POST`, `PUT`, or `DELETE` to `/api/admin/*`:

1. The `access_token` cookie is sent automatically.
2. Add header: `X-XSRF-TOKEN` = `<csrf_token_value_from_step_2>`
3. Send the request.

### Step 4: Refresh Token

```
POST http://localhost:8080/api/auth/refresh
```

No body needed — the `refresh_token` cookie is sent automatically. New cookies are set in the response.

### Step 5: Check Session

```
GET http://localhost:8080/api/auth/me
```

Returns current user info if authenticated.

### Postman Pre-request Script (Optional Automation)

Add this to the Collection's **Pre-request Script** to auto-set the CSRF header:

```javascript
const csrfCookie = pm.cookies.get("csrf_token");
if (csrfCookie) {
    pm.request.headers.add({
        key: "X-XSRF-TOKEN",
        value: csrfCookie
    });
}
```

### Troubleshooting

| Issue | Solution |
|-------|----------|
| `401 Authentication required` | Login first, or token expired — call `/api/auth/refresh` |
| `403 CSRF token missing` | Add `X-XSRF-TOKEN` header with value from `csrf_token` cookie |
| `403 Invalid CSRF token` | CSRF token expired — refresh tokens or re-login |
| Cookies not being sent | Ensure Postman cookie domain matches your request URL |

---

## Frontend Integration

### React/Next.js Example

```typescript
// lib/api.ts
const API_URL = process.env.NEXT_PUBLIC_API_URL || "https://api.adirdk.com";

// Helper function for API calls
async function fetchAPI(url: string, options: RequestInit = {}) {
  // Read CSRF token from cookie
  const csrfToken = document.cookie
    .split("; ")
    .find((row) => row.startsWith("csrf_token="))
    ?.split("=")[1];

  const response = await fetch(`${API_URL}${url}`, {
    ...options,
    credentials: "include", // ← CRITICAL: sends cookies cross-origin
    headers: {
      "Content-Type": "application/json",
      ...(csrfToken && { "X-XSRF-TOKEN": csrfToken }),
      ...options.headers,
    },
  });

  if (response.status === 401) {
    // Try to refresh token
    const refreshRes = await fetch(`${API_URL}/api/auth/refresh`, {
      method: "POST",
      credentials: "include",
    });
    if (refreshRes.ok) {
      // Retry original request with new tokens
      return fetchAPI(url, options);
    }
    // Refresh failed — redirect to login
    window.location.href = "/login";
  }

  return response;
}

// Login
async function login(username: string, password: string) {
  const res = await fetchAPI("/api/auth/login", {
    method: "POST",
    body: JSON.stringify({ username, password }),
  });
  return res.json();
}

// Get current user
async function getMe() {
  const res = await fetchAPI("/api/auth/me");
  return res.json();
}

// Create skill (admin)
async function createSkill(data: any) {
  const res = await fetchAPI("/api/admin/skills", {
    method: "POST",
    body: JSON.stringify(data),
  });
  return res.json();
}

// Upload portfolio (multipart)
async function createPortfolio(formData: FormData) {
  const csrfToken = document.cookie
    .split("; ")
    .find((row) => row.startsWith("csrf_token="))
    ?.split("=")[1];

  const res = await fetch(`${API_URL}/api/admin/portfolio`, {
    method: "POST",
    credentials: "include",
    headers: {
      ...(csrfToken && { "X-XSRF-TOKEN": csrfToken }),
      // Do NOT set Content-Type for FormData — browser sets it with boundary
    },
    body: formData,
  });
  return res.json();
}

// Logout
async function logout() {
  await fetchAPI("/api/auth/logout", { method: "POST" });
}
```

### Axios Example

```typescript
import axios from "axios";

const api = axios.create({
  baseURL: "https://api.adirdk.com",
  withCredentials: true, // ← sends cookies
});

// Axios interceptor to auto-add CSRF header
api.interceptors.request.use((config) => {
  const csrfToken = document.cookie
    .split("; ")
    .find((row) => row.startsWith("csrf_token="))
    ?.split("=")[1];
  if (csrfToken) {
    config.headers["X-XSRF-TOKEN"] = csrfToken;
  }
  return config;
});

// Auto-refresh on 401
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401 && !error.config._retry) {
      error.config._retry = true;
      try {
        await api.post("/api/auth/refresh");
        return api(error.config); // retry
      } catch {
        window.location.href = "/login";
      }
    }
    return Promise.reject(error);
  }
);
```

### Important Notes for Frontend

1. **`credentials: "include"`** (fetch) or **`withCredentials: true`** (axios) is **mandatory** — without it, cookies won't be sent cross-origin.
2. The CORS config on the backend allows credentials for specific origins (not wildcard `*`).
3. CSRF token is in a **non-HttpOnly** cookie — read it with `document.cookie` and send in `X-XSRF-TOKEN` header.
4. Implement **silent refresh** — when access token expires (401), call `/api/auth/refresh` automatically.

---

## Deployment

### Docker

```bash
# Development
docker compose up -d

# Production
docker compose -f docker-compose.prod.yml up -d
```

### Google Cloud Run

The project includes Cloud Run support with automatic Cloud SQL socket detection:

```bash
# Build and push
docker build -f Dockerfile -t gcr.io/<PROJECT>/backend-portfolio .
docker push gcr.io/<PROJECT>/backend-portfolio

# Deploy
gcloud run deploy backend-portfolio \
  --image gcr.io/<PROJECT>/backend-portfolio \
  --add-cloudsql-instances <INSTANCE> \
  --set-env-vars "DB_HOST=/cloudsql/<INSTANCE>,DB_USER=...,DB_PASSWORD=...,DB_NAME=...,JWT_SECRET=...,JWT_REFRESH_SECRET=...,COOKIE_DOMAIN=.adirdk.com,COOKIE_SECURE=true,COOKIE_SAMESITE=Strict,RESET_SECRET=..."
```

The server automatically detects Cloud Run via the `K_SERVICE` environment variable and adjusts the database connection to use Unix sockets.

### Production Checklist

- [ ] Set strong `JWT_SECRET` and `JWT_REFRESH_SECRET` (min 32 random chars each)
- [ ] Set `COOKIE_SECURE=true`
- [ ] Set `COOKIE_SAMESITE=Strict`
- [ ] Set `COOKIE_DOMAIN` to your domain (e.g., `.adirdk.com`)
- [ ] Set `GIN_MODE=release`
- [ ] Set `RESET_SECRET` for password reset endpoint
- [ ] Configure `ALLOWED_ORIGINS` with only your domains
- [ ] Use HTTPS (required for Secure cookies)
- [ ] Set up PostgreSQL with strong credentials
- [ ] Ensure `uploads/` directory has proper permissions

---

## Security Features

### Implemented

| Feature | Implementation |
|---------|----------------|
| **XSS Token Theft Prevention** | JWT stored in HttpOnly cookies (inaccessible to JavaScript) |
| **CSRF Protection** | Signed Double Submit Cookie with HMAC verification |
| **Token Rotation** | Both access and refresh tokens are rotated on refresh |
| **Restricted Refresh Cookie** | Refresh token cookie path limited to `/api/auth/refresh` |
| **Short-lived Access Token** | 15-minute expiry, auto-refreshed by frontend |
| **Secure Cookie Attributes** | Secure, SameSite=Strict, HttpOnly flags |
| **HMAC Token Signing** | Separate secrets for access and refresh tokens |
| **Password Hashing** | bcrypt with default cost |
| **Security Headers** | X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, Referrer-Policy |
| **CORS Whitelisting** | Only specified origins allowed with credentials |
| **Algorithm Pinning** | JWT validation enforces HMAC signing method |
| **Role-based Access** | Admin-only middleware for mutation endpoints |

### Middleware Pipeline (Admin Requests)

```
Request → SecurityHeaders → CORS → AuthMiddleware → CSRFMiddleware → Handler
```

1. **SecurityHeaders:** Sets X-Frame-Options, etc.
2. **CORS:** Validates origin, handles preflight OPTIONS
3. **AuthMiddleware:** Reads `access_token` cookie → validates JWT → checks admin role
4. **CSRFMiddleware:** Compares `X-XSRF-TOKEN` header vs `csrf_token` cookie, verifies HMAC signature

---

## Database Management

### Access PostgreSQL Container

```bash
docker exec -it backend-portfolio-db-1 psql -U portfolio_user -d portfolio_db
```

### Useful Commands

```sql
-- List all tables
\dt

-- View users
SELECT * FROM users;

-- View about data
SELECT * FROM abouts;

-- Exit psql
\q
```

### Running Migrations

```bash
# Copy migration file to container
docker cp scripts/migrations/001_add_new_columns_to_abouts.sql backend-portfolio-db-1:/tmp/

# Execute migration
docker exec -it backend-portfolio-db-1 psql -U portfolio_user -d portfolio_db -f /tmp/001_add_new_columns_to_abouts.sql
```

### Drop Columns

```sql
ALTER TABLE table_name DROP COLUMN column_name;

-- Multiple columns at once
ALTER TABLE table_name
DROP COLUMN col1,
DROP COLUMN col2;
```

> **Note:** GORM auto-migrates model changes on startup. Manual migrations are only needed for complex schema changes (renaming columns, data migrations, etc.).

---

## License

MIT