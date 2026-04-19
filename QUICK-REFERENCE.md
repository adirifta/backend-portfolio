# Quick Reference: JWT Bearer Auth API

## Authentication Endpoints

### POST /api/auth/login
Login and get tokens.

**Request:**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "message": "Login successful",
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

### POST /api/auth/refresh
Refresh access token.

**Request:**
```json
{
  "refresh_token": "eyJ..."
}
```

**Response:** `200 OK`
```json
{
  "message": "Token refreshed successfully",
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "user": { ... }
}
```

### GET /api/auth/me
Get current user (requires auth).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`
```json
{
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

### POST /api/auth/logout
Logout (optional - just remove token from frontend).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`
```json
{
  "message": "Logged out successfully"
}
```

## Admin Routes (All Protected)

All require: `Authorization: Bearer <access_token>` header

### Users
- `POST /api/admin/create-user` - Create new user
- `POST /api/admin/reset-admin` - Reset admin password

### About
- `POST /api/admin/about` - Create about
- `PUT /api/admin/about/:id` - Update about

### Portfolio
- `POST /api/admin/portfolio` - Create portfolio item
- `PUT /api/admin/portfolio/:id` - Update portfolio item
- `DELETE /api/admin/portfolio/:id` - Delete portfolio item
- `DELETE /api/admin/portfolio-media/:portfolio_id/:media_id` - Delete media

### Skills
- `POST /api/admin/skills` - Create skill
- `PUT /api/admin/skills/:id` - Update skill
- `DELETE /api/admin/skills/:id` - Delete skill

### Qualifications
- `POST /api/admin/qualifications` - Create qualification
- `PUT /api/admin/qualifications/:id` - Update qualification
- `DELETE /api/admin/qualifications/:id` - Delete qualification

### Visitor Stats (Admin)
- `GET /api/admin/visitors/stats` - Get visitor stats (includes IPs)

## Public Routes (No Auth Required)

- `GET /api/about` - Get about
- `GET /api/portfolio` - Get all portfolio
- `GET /api/portfolio/:id` - Get portfolio by ID
- `GET /api/skills` - Get all skills
- `GET /api/skills/:id` - Get skill by ID
- `GET /api/qualifications` - Get all qualifications
- `GET /api/qualifications/:id` - Get qualification by ID
- `POST /api/visitors/track` - Track visitor
- `GET /api/visitors/stats` - Get public visitor stats
- `GET /health` - Health check

## JavaScript Quick Start

### Setup
```javascript
const API_URL = 'http://localhost:8080';

// Login
async function login(username, password) {
  const res = await fetch(`${API_URL}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  });
  const data = await res.json();
  localStorage.setItem('access_token', data.access_token);
  localStorage.setItem('refresh_token', data.refresh_token);
  return data.user;
}

// Make protected request
async function apiCall(endpoint, options = {}) {
  const token = localStorage.getItem('access_token');
  const res = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers: {
      ...options.headers,
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  return res.json();
}

// Logout
function logout() {
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
}
```

### Usage
```javascript
// Login
await login('admin', 'password');

// Create user
await apiCall('/api/admin/create-user', {
  method: 'POST',
  body: JSON.stringify({ username: 'john', password: 'secure123', role: 'admin' })
});

// Get user
const user = await apiCall('/api/auth/me');

// Logout
logout();
```

## Token Expiry
- **Access Token:** 15 minutes
- **Refresh Token:** 7 days

## Storage Recommendation
- Use `sessionStorage` for high security (auto-clears on browser close)
- Use `localStorage` for "remember me" functionality
- Both protected by CSP headers against XSS

## CORS Note
Frontend must be on allowed origin:
```
http://localhost:3000
http://localhost:8080
https://adirdk.cloud
https://dashboard.adirdk.com
```

Contact backend team if you need to add your domain.
