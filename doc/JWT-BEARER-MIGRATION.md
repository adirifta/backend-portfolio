# Frontend Integration Guide - JWT Bearer Authentication

This document explains how to integrate with the backend API after the migration from CSRF + Cookie authentication to JWT Bearer tokens.

## Overview

The backend now uses **JWT Bearer tokens** instead of HttpOnly cookies + CSRF tokens. This is simpler and more secure, especially against XSS attacks.

### Key Changes
- ✅ Tokens returned in **JSON response body** (not set as cookies)
- ✅ Tokens sent via **`Authorization: Bearer <token>` header** (not cookies)
- ✅ No CSRF token handling needed
- ✅ Content Security Policy (CSP) protects against XSS
- ✅ Stateless authentication (scales better)

## Authentication Flow

### 1. Login
**Request:**
```bash
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "your-password"
}
```

**Response (200 OK):**
```json
{
  "message": "Login successful",
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

**Frontend Action:**
```javascript
// Store tokens securely
localStorage.setItem('access_token', response.access_token);
localStorage.setItem('refresh_token', response.refresh_token);
```

### 2. Using Access Token (Making Requests)
**All authenticated requests require the `Authorization` header:**

```bash
GET /api/admin/create-user
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**JavaScript Example:**
```javascript
const accessToken = localStorage.getItem('access_token');
const response = await fetch('/api/admin/users', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${accessToken}`,
    'Content-Type': 'application/json'
  }
});
```

### 3. Token Refresh
**When access token expires (after 15 minutes):**

**Request:**
```bash
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (200 OK):**
```json
{
  "message": "Token refreshed successfully",
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

**Frontend Action:**
```javascript
// Update stored tokens
localStorage.setItem('access_token', response.access_token);
localStorage.setItem('refresh_token', response.refresh_token);
```

### 4. Logout
**Request:**
```bash
POST /api/auth/logout
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "message": "Logged out successfully"
}
```

**Frontend Action:**
```javascript
// Remove tokens from storage
localStorage.removeItem('access_token');
localStorage.removeItem('refresh_token');

// Redirect to login
window.location.href = '/login';
```

## Implementation Example (JavaScript/React)

### Fetch Wrapper with Auto-Refresh
```javascript
const API_BASE = 'http://localhost:8080';

async function apiCall(endpoint, options = {}) {
  let token = localStorage.getItem('access_token');
  
  if (!token) {
    throw new Error('Not authenticated');
  }

  let response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      ...options.headers,
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });

  // If token expired (401), try to refresh
  if (response.status === 401) {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) {
      throw new Error('Session expired, please login');
    }

    // Attempt refresh
    const refreshResponse = await fetch(`${API_BASE}/api/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken })
    });

    if (refreshResponse.ok) {
      const refreshData = await refreshResponse.json();
      localStorage.setItem('access_token', refreshData.access_token);
      localStorage.setItem('refresh_token', refreshData.refresh_token);
      
      // Retry original request with new token
      return apiCall(endpoint, options);
    } else {
      // Refresh failed, force logout
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      window.location.href = '/login';
      throw new Error('Session expired');
    }
  }

  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }

  return response.json();
}
```

### React Hook Example
```javascript
import { useState, useEffect } from 'react';

function useAuth() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (token) {
      // Fetch current user
      apiCall('/api/auth/me').then(setUser).catch(() => {
        localStorage.removeItem('access_token');
      }).finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  const login = async (username, password) => {
    const response = await fetch(`${API_BASE}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    });
    
    const data = await response.json();
    localStorage.setItem('access_token', data.access_token);
    localStorage.setItem('refresh_token', data.refresh_token);
    setUser(data.user);
    return data;
  };

  const logout = async () => {
    const token = localStorage.getItem('access_token');
    await fetch(`${API_BASE}/api/auth/logout`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` }
    });
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    setUser(null);
  };

  return { user, loading, login, logout, isAuthenticated: !!user };
}
```

## Token Storage Security Considerations

### localStorage vs sessionStorage
- **localStorage**: Persists across browser sessions (remember me functionality)
  - Risk: XSS attacks can steal tokens
  - Mitigation: CSP headers prevent inline scripts

- **sessionStorage**: Cleared when browser closes
  - Safer for sensitive sessions
  - Choose based on your security requirements

### Best Practices
1. ✅ Use **httpOnly** is NOT an option here (tokens must be in JavaScript for Bearer header)
2. ✅ Store in **sessionStorage** for high-security (auto-clears on browser close)
3. ✅ OR store in **localStorage** if user wants "remember me"
4. ✅ Backend enforces CSP to prevent XSS
5. ✅ Always use **HTTPS in production** (never HTTP)
6. ✅ Validate tokens have not been tampered with (backend does this)

## CORS Configuration

Frontend must run on an allowed origin. Default allowed origins:
- https://adirdk.cloud
- https://dashboard.adirdk.com
- http://localhost:3000 (development)
- http://localhost:8080 (development)

Add your origin to `ALLOWED_ORIGINS` env var if needed:
```bash
ALLOWED_ORIGINS=http://localhost:3000,https://myapp.com
```

## Environment Variables

### .env for Backend
```bash
JWT_SECRET=your-secret-key-min-32-chars
JWT_REFRESH_SECRET=your-refresh-secret-key-min-32-chars
ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
```

### Frontend .env
```bash
REACT_APP_API_URL=http://localhost:8080
REACT_APP_TOKEN_STORAGE=sessionStorage  # or localStorage
```

## Error Handling

### Common Error Responses

**401 Unauthorized**
```json
{
  "error": "Authentication required (Bearer token)"
}
```
- Missing or invalid Authorization header
- Token expired
- Action: Refresh token or redirect to login

**403 Forbidden**
```json
{
  "error": "Admin access required"
}
```
- User authenticated but lacks admin role
- Action: Show permission denied message

**400 Bad Request**
```json
{
  "error": "Invalid request"
}
```
- Validation error
- Check request body/parameters

## Testing with curl

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# Use access token
curl http://localhost:8080/api/admin/create-user \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Refresh token
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'

# Logout
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Migration Checklist

If migrating from CSRF + Cookie auth:

- [ ] Remove all CSRF token handling from frontend
- [ ] Stop reading `csrf_token` cookie
- [ ] Stop sending `X-XSRF-TOKEN` header
- [ ] Update login to store tokens from response body
- [ ] Update all API calls to include `Authorization` header
- [ ] Implement token refresh logic
- [ ] Update logout to clear localStorage/sessionStorage
- [ ] Test with Chrome DevTools Network tab (verify Bearer headers)
- [ ] Remove any cookie-related code
- [ ] Update tests to use Bearer tokens

## Troubleshooting

**"Authentication required (Bearer token)" error**
- Check Authorization header is being sent
- Verify token format: `Bearer <token>`
- Ensure token hasn't expired (15 min)
- Try refreshing token

**"CORS policy: No 'Access-Control-Allow-Origin' header"**
- Frontend URL not in ALLOWED_ORIGINS
- Add frontend URL to backend ALLOWED_ORIGINS env var
- Ensure credentials setting is correct in fetch

**Token refresh not working**
- Verify refresh_token is valid and not expired (7 days)
- Check refresh token value is correct
- Ensure POST body includes `refresh_token` field

## Additional Resources

- [JWT Introduction](https://jwt.io/introduction)
- [OWASP: Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [XSS Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_Prevention_Cheat_Sheet.html)
