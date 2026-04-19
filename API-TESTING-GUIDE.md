# API Testing Guide - JWT Bearer Authentication

**Last Updated:** 2026-04-18  
**Status:** ✅ All APIs migrated from CSRF + Cookie to JWT Bearer tokens

---

## 📋 Quick Start

### 1. **Import Postman Collection**
```
File: Portfolio-API-Postman-Collection.json
```
- Open Postman
- Click "Import"
- Select the JSON file
- All 27 endpoints ready to test

### 2. **Set Environment Variables**
```
base_url: http://localhost:8080
access_token: (will be set after login)
refresh_token: (will be set after login)
```

### 3. **Test Login First**
```
POST /api/auth/login
Body: { "username": "admin", "password": "password123" }
```

---

## 🔐 Authentication Flow

### Step 1: Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

**Response:**
```json
{
  "message": "Login successful",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

**Action:** Save `access_token` and `refresh_token` for next requests.

### Step 2: Use Token in Protected Endpoints
```bash
curl -X GET http://localhost:8080/api/admin/visitors/stats \
  -H "Authorization: Bearer <access_token>"
```

**Token expiry:**
- Access Token: 15 minutes
- Refresh Token: 7 days

### Step 3: Refresh When Expired
```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'
```

---

## 📊 Complete API Endpoint List

### ✅ **Authentication (No Auth Required)**
| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/api/auth/login` | Get JWT tokens |
| POST | `/api/auth/logout` | Logout |
| POST | `/api/auth/refresh` | Refresh access token |
| GET | `/api/auth/me` | Get current user (requires Bearer) |
| POST | `/api/login` | Login (backward compat) |

### ✅ **About**
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/api/about` | ❌ | Get about (public) |
| POST | `/api/admin/about` | ✅ | Create about |
| PUT | `/api/admin/about/:id` | ✅ | Update about |

### ✅ **Portfolio**
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/api/portfolio` | ❌ | Get all (public) |
| GET | `/api/portfolio/:id` | ❌ | Get one (public) |
| POST | `/api/admin/portfolio` | ✅ | Create |
| PUT | `/api/admin/portfolio/:id` | ✅ | Update |
| DELETE | `/api/admin/portfolio/:id` | ✅ | Delete |
| DELETE | `/api/admin/portfolio-media/:portfolio_id/:media_id` | ✅ | Delete media |

### ✅ **Skills**
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/api/skills` | ❌ | Get all (public) |
| GET | `/api/skills/:id` | ❌ | Get one (public) |
| GET | `/api/skills/category` | ❌ | Get by category (public) |
| POST | `/api/admin/skills` | ✅ | Create |
| PUT | `/api/admin/skills/:id` | ✅ | Update |
| DELETE | `/api/admin/skills/:id` | ✅ | Delete |

### ✅ **Qualifications**
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/api/qualifications` | ❌ | Get all (public) |
| GET | `/api/qualifications/:id` | ❌ | Get one (public) |
| POST | `/api/admin/qualifications` | ✅ | Create |
| PUT | `/api/admin/qualifications/:id` | ✅ | Update |
| DELETE | `/api/admin/qualifications/:id` | ✅ | Delete |

### ✅ **Users**
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| POST | `/api/admin/create-user` | ✅ | Create user |
| POST | `/api/admin/reset-admin` | ✅ | Reset password |

### ✅ **Visitor Analytics**
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| POST | `/api/visitors/track` | ❌ | Track visit |
| GET | `/api/visitors/stats` | ❌ | Get public stats |
| GET | `/api/admin/visitors/stats` | ✅ | Get admin stats |

### ✅ **Health & Info**
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/health` | ❌ | Health check |
| GET | `/` | ❌ | API info |

---

## 🧪 Testing Scenarios

### Scenario 1: Public Read-Only
```bash
# No auth required
curl http://localhost:8080/api/about
curl http://localhost:8080/api/portfolio
curl http://localhost:8080/api/skills
curl http://localhost:8080/api/qualifications
curl http://localhost:8080/api/health
```

### Scenario 2: Authenticated Admin Operations
```bash
# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}' \
  | jq -r '.access_token')

# 2. Use token for admin operations
curl -X POST http://localhost:8080/api/admin/skills \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Go","level":"Advanced","score":85,"category":"Backend","icon":"go-icon"}'
```

### Scenario 3: Token Refresh
```bash
# 1. Get refresh token from login
REFRESH=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}' \
  | jq -r '.refresh_token')

# 2. Wait for access token to expire (15 min)
# 3. Refresh tokens
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH\"}"
```

---

## ✅ Verification Checklist

### Authentication
- [ ] POST `/api/auth/login` returns tokens in body ✓
- [ ] Tokens NOT in cookies (checked headers) ✓
- [ ] POST `/api/auth/refresh` accepts JSON body ✓
- [ ] Bearer token required for admin routes ✓
- [ ] No X-XSRF-TOKEN header needed ✓

### Public Endpoints
- [ ] GET `/api/about` works without auth ✓
- [ ] GET `/api/portfolio` works without auth ✓
- [ ] GET `/api/skills` works without auth ✓
- [ ] GET `/api/qualifications` works without auth ✓
- [ ] GET `/api/visitors/stats` works without auth ✓
- [ ] GET `/health` works without auth ✓

### Admin Endpoints (Protected)
- [ ] POST `/api/admin/about` requires Bearer token ✓
- [ ] PUT `/api/admin/about/:id` requires Bearer token ✓
- [ ] POST `/api/admin/portfolio` requires Bearer token ✓
- [ ] POST `/api/admin/skills` requires Bearer token ✓
- [ ] POST `/api/admin/qualifications` requires Bearer token ✓
- [ ] POST `/api/admin/create-user` requires Bearer token ✓
- [ ] GET `/api/admin/visitors/stats` requires Bearer token ✓

### Error Cases
- [ ] Missing Bearer token returns 401 ✓
- [ ] Invalid token returns 401 ✓
- [ ] Non-admin user gets 403 ✓
- [ ] Missing required fields returns 400 ✓
- [ ] Expired token returns 401 ✓

### Headers & CORS
- [ ] Authorization header set correctly ✓
- [ ] No Cookie headers for auth ✓
- [ ] CSP headers present ✓
- [ ] X-Frame-Options: DENY present ✓
- [ ] X-Content-Type-Options: nosniff present ✓
- [ ] CORS allows requests from allowed origins ✓

---

## 🚀 curl Command Templates

### Create Skill
```bash
curl -X POST http://localhost:8080/api/admin/skills \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "React",
    "level": "Advanced",
    "score": 90,
    "category": "Frontend",
    "icon": "react-icon"
  }'
```

### Create Qualification
```bash
curl -X POST http://localhost:8080/api/admin/qualifications \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "education",
    "institution": "University",
    "title": "Bachelor CS",
    "description": "Computer Science",
    "start_date": "2020-01-01T00:00:00Z",
    "end_date": "2024-06-01T00:00:00Z",
    "current": false
  }'
```

### Create Portfolio
```bash
curl -X POST http://localhost:8080/api/admin/portfolio \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Project",
    "description": "Project description",
    "category": "Web",
    "tags": "React, Node.js",
    "project_url": "https://example.com"
  }'
```

### Get Admin Stats
```bash
curl -X GET http://localhost:8080/api/admin/visitors/stats \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🐛 Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| **401 Unauthorized** | Check Bearer token format: `Authorization: Bearer <token>` |
| **Missing token error** | Login first with `/api/auth/login` |
| **Token expired** | Use `/api/auth/refresh` with refresh_token |
| **403 Forbidden** | Verify user has admin role |
| **CORS error** | Add frontend origin to `ALLOWED_ORIGINS` env var |
| **400 Bad Request** | Check JSON body format and required fields |
| **No cookies in response** | ✓ Expected! Tokens now in JSON body |
| **X-XSRF-TOKEN header not needed** | ✓ Correct! Removed after migration |

---

## 📈 Response Examples

### Successful Login
```json
{
  "message": "Login successful",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

### Successful Skill Creation
```json
{
  "message": "Skill created successfully",
  "skill": {
    "id": 1,
    "name": "React",
    "level": "Advanced",
    "score": 90,
    "category": "Frontend",
    "icon": "react-icon",
    "created_at": "2026-04-18T08:11:40Z"
  }
}
```

### Error Response
```json
{
  "error": "Invalid credentials"
}
```

---

## 🎯 Testing Checklist Summary

```
✅ All 27 endpoints working with JWT Bearer tokens
✅ No cookies used for authentication
✅ No CSRF tokens or headers needed
✅ CSP headers added for XSS protection
✅ All admin routes require Bearer token
✅ Public routes accessible without auth
✅ Token refresh working correctly
✅ Error handling proper (401, 403, 400)
✅ CORS properly configured
✅ Docker build successful
✅ Documentation updated
```

---

**Status:** Ready for production testing! 🚀
