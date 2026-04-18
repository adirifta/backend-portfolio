# Frontend Setup Guide - CSRF & Cookie Authentication

Dokumentasi lengkap untuk mengintegrasikan frontend dengan backend authentication yang menggunakan CSRF tokens dan cross-site cookies.

---

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Setup by Framework](#setup-by-framework)
- [API Integration](#api-integration)
- [Authentication Flow](#authentication-flow)
- [Troubleshooting](#troubleshooting)
- [Security Best Practices](#security-best-practices)

---

## Overview

### Backend Configuration

```yaml
# Backend (https://adirdk.cloud)
COOKIE_DOMAIN=.adirdk.cloud
COOKIE_SECURE=true
COOKIE_SAMESITE=None
ALLOWED_ORIGINS=https://dashboard.adirdk.com

# Cookies set after login:
- access_token (15 min expiry)
- refresh_token (7 days expiry)
- csrf_token (readable, for JS)
- csrf_token_sig (HttpOnly, for server validation)
```

### Frontend Requirements

```
Frontend: https://dashboard.adirdk.com
Backend:  https://adirdk.cloud (cross-site)

Requirements:
1. ✅ credentials: 'include' on ALL requests
2. ✅ X-XSRF-TOKEN header on POST/PUT/DELETE
3. ✅ Handle token refresh automatically
4. ✅ Graceful logout on 401 errors
```

---

## Architecture

### Authentication Flow

```
┌─────────────────────────────────────────────────────────┐
│ Frontend (dashboard.adirdk.com)                          │
└──────────────────┬──────────────────────────────────────┘
                   │
                   │ 1. POST /api/auth/login
                   │    + credentials: 'include'
                   ▼
┌─────────────────────────────────────────────────────────┐
│ Backend (adirdk.cloud)                                   │
│ Set-Cookie:                                              │
│  - access_token (15 min)                                 │
│  - refresh_token (7 days)                                │
│  - csrf_token (readable)                                 │
│  - csrf_token_sig (HttpOnly)                             │
└──────────────────┬──────────────────────────────────────┘
                   │
                   │ 2. Response + Set-Cookie headers
                   ▼
┌─────────────────────────────────────────────────────────┐
│ Browser                                                  │
│ document.cookie: access_token, csrf_token, ...           │
└──────────────────┬──────────────────────────────────────┘
                   │
                   │ 3. POST /api/admin/portfolio
                   │    + credentials: 'include'
                   │    + X-XSRF-TOKEN header
                   ▼
┌─────────────────────────────────────────────────────────┐
│ Backend                                                  │
│ 1. Validate CSRF token                                   │
│ 2. Validate access_token cookie                          │
│ 3. Process request                                       │
└─────────────────────────────────────────────────────────┘
```

### Cookie Lifecycle

| Cookie | Expiry | HttpOnly | Purpose |
|--------|--------|----------|---------|
| `access_token` | 15 min | Yes | JWT access token for authorization |
| `refresh_token` | 7 days | Yes | Used to get new access token |
| `csrf_token` | 15 min | **No** | Readable by JS, send in X-XSRF-TOKEN header |
| `csrf_token_sig` | 15 min | Yes | HMAC signature for server validation |

---

## Setup by Framework

### 🔵 Vanilla JavaScript / Fetch API

**File: `src/api.js`**

```javascript
// Constants
const API_BASE_URL = 'https://adirdk.cloud';
const CSRF_COOKIE_NAME = 'csrf_token';
const CSRF_HEADER_NAME = 'X-XSRF-TOKEN';

// Utility: Extract cookie value
function getCookie(name) {
  const cookies = document.cookie.split('; ');
  const cookie = cookies.find(c => c.startsWith(name + '='));
  return cookie ? cookie.split('=')[1] : null;
}

// Utility: Generic API call
export async function apiCall(endpoint, options = {}) {
  const {
    method = 'GET',
    body = null,
    headers = {},
  } = options;

  // Get CSRF token for non-GET requests
  const csrfToken = ['POST', 'PUT', 'DELETE'].includes(method)
    ? getCookie(CSRF_COOKIE_NAME)
    : null;

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method,
    credentials: 'include',  // ← CRITICAL: Send/receive cookies
    headers: {
      'Content-Type': 'application/json',
      ...(csrfToken && { [CSRF_HEADER_NAME]: csrfToken }),
      ...headers,
    },
    ...(body && { body: JSON.stringify(body) }),
  });

  // Handle 401: Token expired, try refresh
  if (response.status === 401) {
    const refreshed = await refreshAccessToken();
    if (refreshed) {
      return apiCall(endpoint, options);  // Retry original request
    }
    // Redirect to login
    window.location.href = '/login';
    return null;
  }

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || `HTTP ${response.status}`);
  }

  return response.json();
}

// Auth: Login
export async function login(username, password) {
  const data = await apiCall('/api/auth/login', {
    method: 'POST',
    body: { username, password },
  });
  console.log('✅ Login successful, cookies set:', document.cookie);
  return data;
}

// Auth: Logout
export async function logout() {
  await apiCall('/api/auth/logout', { method: 'POST' });
  console.log('✅ Logged out');
  // Redirect or refresh UI
  window.location.href = '/login';
}

// Auth: Get current user
export async function getMe() {
  return apiCall('/api/auth/me');
}

// Auth: Refresh token
export async function refreshAccessToken() {
  try {
    await apiCall('/api/auth/refresh', { method: 'POST' });
    console.log('✅ Token refreshed');
    return true;
  } catch (error) {
    console.error('❌ Token refresh failed:', error);
    return false;
  }
}

// CRUD: Portfolio
export async function createPortfolio(data) {
  return apiCall('/api/admin/portfolio', {
    method: 'POST',
    body: data,
  });
}

export async function updatePortfolio(id, data) {
  return apiCall(`/api/admin/portfolio/${id}`, {
    method: 'PUT',
    body: data,
  });
}

export async function deletePortfolio(id) {
  return apiCall(`/api/admin/portfolio/${id}`, {
    method: 'DELETE',
  });
}

// Example: Usage
(async () => {
  try {
    await login('admin', 'password');
    const me = await getMe();
    console.log('Current user:', me);

    const portfolio = await createPortfolio({
      title: 'My Project',
      description: 'Project description',
    });
    console.log('Created portfolio:', portfolio);
  } catch (error) {
    console.error('Error:', error.message);
  }
})();
```

---

### ⚛️ React (Hooks + Fetch)

**File: `src/hooks/useApi.js`**

```javascript
import { useState, useCallback } from 'react';

const API_BASE_URL = 'https://adirdk.cloud';

export function useApi() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const getCookie = useCallback((name) => {
    const cookies = document.cookie.split('; ');
    const cookie = cookies.find(c => c.startsWith(name + '='));
    return cookie ? cookie.split('=')[1] : null;
  }, []);

  const call = useCallback(
    async (endpoint, options = {}) => {
      const {
        method = 'GET',
        body = null,
        headers = {},
      } = options;

      setLoading(true);
      setError(null);

      try {
        const csrfToken = ['POST', 'PUT', 'DELETE'].includes(method)
          ? getCookie('csrf_token')
          : null;

        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
          method,
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
            ...(csrfToken && { 'X-XSRF-TOKEN': csrfToken }),
            ...headers,
          },
          ...(body && { body: JSON.stringify(body) }),
        });

        if (response.status === 401) {
          // Try refresh
          const refreshResp = await fetch(`${API_BASE_URL}/api/auth/refresh`, {
            method: 'POST',
            credentials: 'include',
          });
          
          if (refreshResp.ok) {
            // Retry original request
            return call(endpoint, options);
          } else {
            window.location.href = '/login';
            return null;
          }
        }

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || `HTTP ${response.status}`);
        }

        const data = await response.json();
        setLoading(false);
        return data;
      } catch (err) {
        setError(err.message);
        setLoading(false);
        throw err;
      }
    },
    [getCookie]
  );

  return { call, loading, error };
}
```

**File: `src/context/AuthContext.js`**

```javascript
import { createContext, useContext, useState, useCallback } from 'react';
import { useApi } from '../hooks/useApi';

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const { call } = useApi();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(false);

  const login = useCallback(async (username, password) => {
    setLoading(true);
    try {
      const data = await call('/api/auth/login', {
        method: 'POST',
        body: { username, password },
      });
      setUser(data.user);
      return data;
    } finally {
      setLoading(false);
    }
  }, [call]);

  const logout = useCallback(async () => {
    setLoading(true);
    try {
      await call('/api/auth/logout', { method: 'POST' });
      setUser(null);
      window.location.href = '/login';
    } finally {
      setLoading(false);
    }
  }, [call]);

  const getMe = useCallback(async () => {
    try {
      const data = await call('/api/auth/me');
      setUser(data.user);
      return data.user;
    } catch (err) {
      setUser(null);
      throw err;
    }
  }, [call]);

  return (
    <AuthContext.Provider value={{ user, login, logout, getMe, loading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
```

**File: `src/components/LoginPage.jsx`**

```javascript
import { useState } from 'react';
import { useAuth } from '../context/AuthContext';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const { login, loading, error } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await login(username, password);
      // Redirect on success
      window.location.href = '/dashboard';
    } catch (err) {
      console.error('Login failed:', err.message);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        required
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
      />
      <button type="submit" disabled={loading}>
        {loading ? 'Loading...' : 'Login'}
      </button>
      {error && <p style={{ color: 'red' }}>{error}</p>}
    </form>
  );
}
```

---

### 🟢 Vue 3 (Composition API)

**File: `src/api.js`**

```javascript
import { ref } from 'vue';

const API_BASE_URL = 'https://adirdk.cloud';

export function useApi() {
  const loading = ref(false);
  const error = ref(null);

  const getCookie = (name) => {
    const cookies = document.cookie.split('; ');
    const cookie = cookies.find(c => c.startsWith(name + '='));
    return cookie ? cookie.split('=')[1] : null;
  };

  const call = async (endpoint, options = {}) => {
    const { method = 'GET', body = null, headers = {} } = options;

    loading.value = true;
    error.value = null;

    try {
      const csrfToken = ['POST', 'PUT', 'DELETE'].includes(method)
        ? getCookie('csrf_token')
        : null;

      const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        method,
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          ...(csrfToken && { 'X-XSRF-TOKEN': csrfToken }),
          ...headers,
        },
        ...(body && { body: JSON.stringify(body) }),
      });

      if (response.status === 401) {
        const refreshResp = await fetch(`${API_BASE_URL}/api/auth/refresh`, {
          method: 'POST',
          credentials: 'include',
        });
        if (refreshResp.ok) {
          return call(endpoint, options);
        } else {
          window.location.href = '/login';
          return null;
        }
      }

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || `HTTP ${response.status}`);
      }

      return await response.json();
    } catch (err) {
      error.value = err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return { call, loading, error, getCookie };
}
```

**File: `src/stores/auth.js`**

```javascript
import { defineStore } from 'pinia';
import { useApi } from '../api';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    loading: false,
  }),

  actions: {
    async login(username, password) {
      this.loading = true;
      try {
        const { call } = useApi();
        const data = await call('/api/auth/login', {
          method: 'POST',
          body: { username, password },
        });
        this.user = data.user;
        return data;
      } finally {
        this.loading = false;
      }
    },

    async logout() {
      this.loading = true;
      try {
        const { call } = useApi();
        await call('/api/auth/logout', { method: 'POST' });
        this.user = null;
        window.location.href = '/login';
      } finally {
        this.loading = false;
      }
    },

    async getMe() {
      try {
        const { call } = useApi();
        const data = await call('/api/auth/me');
        this.user = data.user;
        return data.user;
      } catch (err) {
        this.user = null;
        throw err;
      }
    },
  },
});
```

**File: `src/components/LoginPage.vue`**

```vue
<template>
  <form @submit.prevent="handleLogin">
    <input v-model="username" type="text" placeholder="Username" required />
    <input v-model="password" type="password" placeholder="Password" required />
    <button :disabled="loading">{{ loading ? 'Loading...' : 'Login' }}</button>
    <p v-if="error" style="color: red">{{ error }}</p>
  </form>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';

const username = ref('');
const password = ref('');
const authStore = useAuthStore();

const loading = ref(false);
const error = ref(null);

async function handleLogin() {
  loading.value = true;
  error.value = null;
  try {
    await authStore.login(username.value, password.value);
    window.location.href = '/dashboard';
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}
</script>
```

---

### 🔴 Angular (HttpClient)

**File: `src/services/api.service.ts`**

```typescript
import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError, BehaviorSubject } from 'rxjs';
import { catchError, tap, switchMap } from 'rxjs/operators';

@Injectable({ providedIn: 'root' })
export class ApiService {
  private apiUrl = 'https://adirdk.cloud';
  private isRefreshing = false;
  private refreshTokenSubject = new BehaviorSubject<string | null>(null);

  constructor(private http: HttpClient) {}

  private getCookie(name: string): string | null {
    const cookies = document.cookie.split('; ');
    const cookie = cookies.find(c => c.startsWith(name + '='));
    return cookie ? cookie.split('=')[1] : null;
  }

  private handleError(error: HttpErrorResponse) {
    if (error.status === 401) {
      // Redirect to login
      window.location.href = '/login';
    }
    return throwError(() => new Error(error.error?.error || `HTTP ${error.status}`));
  }

  login(username: string, password: string): Observable<any> {
    return this.http.post(
      `${this.apiUrl}/api/auth/login`,
      { username, password },
      { withCredentials: true }  // ← CRITICAL
    ).pipe(catchError(this.handleError));
  }

  logout(): Observable<any> {
    return this.http.post(
      `${this.apiUrl}/api/auth/logout`,
      {},
      { withCredentials: true }
    ).pipe(
      tap(() => {
        window.location.href = '/login';
      }),
      catchError(this.handleError)
    );
  }

  getMe(): Observable<any> {
    return this.http.get(
      `${this.apiUrl}/api/auth/me`,
      { withCredentials: true }
    ).pipe(catchError(this.handleError));
  }

  createPortfolio(data: any): Observable<any> {
    const csrfToken = this.getCookie('csrf_token');
    return this.http.post(
      `${this.apiUrl}/api/admin/portfolio`,
      data,
      {
        withCredentials: true,
        headers: csrfToken ? { 'X-XSRF-TOKEN': csrfToken } : undefined,
      }
    ).pipe(catchError(this.handleError));
  }

  updatePortfolio(id: number, data: any): Observable<any> {
    const csrfToken = this.getCookie('csrf_token');
    return this.http.put(
      `${this.apiUrl}/api/admin/portfolio/${id}`,
      data,
      {
        withCredentials: true,
        headers: csrfToken ? { 'X-XSRF-TOKEN': csrfToken } : undefined,
      }
    ).pipe(catchError(this.handleError));
  }

  deletePortfolio(id: number): Observable<any> {
    const csrfToken = this.getCookie('csrf_token');
    return this.http.delete(
      `${this.apiUrl}/api/admin/portfolio/${id}`,
      {
        withCredentials: true,
        headers: csrfToken ? { 'X-XSRF-TOKEN': csrfToken } : undefined,
      }
    ).pipe(catchError(this.handleError));
  }
}
```

**File: `src/services/auth.service.ts`**

```typescript
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { ApiService } from './api.service';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private userSubject = new BehaviorSubject<any>(null);
  public user$ = this.userSubject.asObservable();

  constructor(private api: ApiService) {}

  login(username: string, password: string): Observable<any> {
    return this.api.login(username, password).pipe(
      tap((data) => {
        this.userSubject.next(data.user);
      })
    );
  }

  logout(): Observable<any> {
    return this.api.logout().pipe(
      tap(() => {
        this.userSubject.next(null);
      })
    );
  }

  getMe(): Observable<any> {
    return this.api.getMe().pipe(
      tap((data) => {
        this.userSubject.next(data.user);
      })
    );
  }

  getCurrentUser(): any {
    return this.userSubject.value;
  }
}
```

---

### 📦 Axios (Global Setup)

**File: `src/api/axios.js`**

```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'https://adirdk.cloud',
  withCredentials: true,  // ← CRITICAL: Send/receive cookies
});

// Interceptor: Add CSRF token to requests
api.interceptors.request.use((config) => {
  const csrfToken = document.cookie
    .split('; ')
    .find(c => c.startsWith('csrf_token='))
    ?.split('=')[1];

  if (csrfToken && ['POST', 'PUT', 'DELETE'].includes(config.method.toUpperCase())) {
    config.headers['X-XSRF-TOKEN'] = csrfToken;
  }

  return config;
});

// Interceptor: Handle 401 (token expired)
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        await api.post('/api/auth/refresh');
        return api(originalRequest);
      } catch (refreshError) {
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default api;
```

**Usage:**

```javascript
import api from './api/axios';

// Login
api.post('/api/auth/login', { username: 'admin', password: 'xxx' })
  .then(response => console.log('Login success:', response.data))
  .catch(error => console.error('Login failed:', error));

// Create portfolio
api.post('/api/admin/portfolio', {
  title: 'My Project',
  description: 'Description'
})
  .then(response => console.log('Created:', response.data))
  .catch(error => console.error('Error:', error));
```

---

## API Integration

### Authentication Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/api/auth/login` | Login with username/password |
| POST | `/api/auth/logout` | Logout (clear cookies) |
| POST | `/api/auth/refresh` | Refresh access token |
| GET | `/api/auth/me` | Get current user info |

### Admin CRUD Endpoints

All require:
- ✅ `credentials: 'include'` (or `withCredentials: true`)
- ✅ `X-XSRF-TOKEN` header (except GET)
- ✅ Valid `access_token` cookie

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/api/admin/portfolio` | Create portfolio |
| PUT | `/api/admin/portfolio/:id` | Update portfolio |
| DELETE | `/api/admin/portfolio/:id` | Delete portfolio |
| POST | `/api/admin/skills` | Create skill |
| PUT | `/api/admin/skills/:id` | Update skill |
| DELETE | `/api/admin/skills/:id` | Delete skill |
| POST | `/api/admin/qualifications` | Create qualification |
| PUT | `/api/admin/qualifications/:id` | Update qualification |
| DELETE | `/api/admin/qualifications/:id` | Delete qualification |

### Response Format

**Success (200):**
```json
{
  "message": "Success message",
  "data": { ... }
}
```

**Error (4xx/5xx):**
```json
{
  "error": "Error description"
}
```

---

## Authentication Flow

### Step 1: Login

```
Frontend → POST /api/auth/login
  {
    "username": "admin",
    "password": "password"
  }
  + credentials: 'include'

Backend ← 200 OK
  Set-Cookie: access_token=...
  Set-Cookie: refresh_token=...
  Set-Cookie: csrf_token=...
  Set-Cookie: csrf_token_sig=...
```

### Step 2: Verify Login (Optional)

```
Frontend → GET /api/auth/me
  + credentials: 'include'
  + Cookie: access_token=...; csrf_token=...

Backend ← 200 OK
  {
    "user": {
      "id": 1,
      "username": "admin",
      "role": "admin"
    }
  }
```

### Step 3: CRUD Operations

```
Frontend → POST /api/admin/portfolio
  {
    "title": "My Project",
    "description": "Description"
  }
  + credentials: 'include'
  + X-XSRF-TOKEN: <csrf_token_value>
  + Cookie: access_token=...; csrf_token=...

Backend ← 201 Created
  {
    "message": "Portfolio created successfully",
    "data": { ... }
  }
```

### Step 4: Token Refresh (Auto)

When access token expires (15 min):

```
Frontend → POST /api/auth/refresh
  + credentials: 'include'
  + Cookie: refresh_token=...

Backend ← 200 OK
  Set-Cookie: access_token=<new_token>
  Set-Cookie: csrf_token=<new_token>

Frontend → Retry original request
```

### Step 5: Logout

```
Frontend → POST /api/auth/logout
  + credentials: 'include'

Backend ← 200 OK (cookies cleared)

Frontend → Redirect to /login
```

---

## Troubleshooting

### ❌ Problem: "CSRF token header (X-XSRF-TOKEN) required"

**Cause:** Not sending X-XSRF-TOKEN header on POST/PUT/DELETE

**Solution:**
```javascript
// ❌ Wrong
fetch('https://adirdk.cloud/api/admin/portfolio', {
  method: 'POST',
  credentials: 'include'  // Missing CSRF header!
})

// ✅ Correct
const csrfToken = document.cookie
  .split('; ')
  .find(c => c.startsWith('csrf_token='))
  ?.split('=')[1];

fetch('https://adirdk.cloud/api/admin/portfolio', {
  method: 'POST',
  credentials: 'include',
  headers: {
    'X-XSRF-TOKEN': csrfToken  // ← Add this!
  }
})
```

---

### ❌ Problem: Cookies not appearing in document.cookie

**Cause:** Missing `credentials: 'include'`

**Solution:**
```javascript
// ❌ Wrong
fetch('https://adirdk.cloud/api/auth/login', {
  method: 'POST'
  // Missing credentials!
})

// ✅ Correct
fetch('https://adirdk.cloud/api/auth/login', {
  method: 'POST',
  credentials: 'include'  // ← Add this!
})
```

---

### ❌ Problem: 401 Unauthorized on auth/me

**Cause:** access_token cookie not being sent or expired

**Solution:**

1. Check if login was successful (check Set-Cookie in Network tab)
2. Verify `credentials: 'include'` is set
3. Check if cookies expired (access_token expires in 15 min)

```javascript
// Debug: Check cookies
console.log('Cookies:', document.cookie);

// Debug: Check Network tab
// Login request → Response Headers → Check Set-Cookie
```

---

### ❌ Problem: "Cookie rejected because... SameSite is Lax/Strict"

**Cause:** Wrong SameSite configuration or missing credentials

**Solution:**

This is normal in browser console warnings. To verify it's working:

1. Check Network → Response Headers → Set-Cookie
2. On next request, check Request Headers → Cookie
3. If cookies are being sent, it's working fine

Backend is configured with `SameSite=None + Secure=true`, so frontend must use `credentials: 'include'`.

---

## Security Best Practices

### 1. ✅ Always use HTTPS in Production

```javascript
// Development
const API_URL = 'http://localhost:8080';

// Production
const API_URL = 'https://adirdk.cloud';
```

### 2. ✅ Handle Token Refresh Automatically

Implement interceptors to refresh tokens before they expire:

```javascript
// Axios example
api.interceptors.response.use(
  response => response,
  async error => {
    if (error.response?.status === 401) {
      // Try refresh
      await api.post('/api/auth/refresh');
      // Retry original request
      return api(error.config);
    }
    return Promise.reject(error);
  }
);
```

### 3. ✅ Validate CSRF Token Exists

Before sending CRUD requests, verify CSRF token is available:

```javascript
const csrfToken = document.cookie
  .split('; ')
  .find(c => c.startsWith('csrf_token='))
  ?.split('=')[1];

if (!csrfToken && ['POST', 'PUT', 'DELETE'].includes(method)) {
  // Redirect to login
  window.location.href = '/login';
  return;
}
```

### 4. ✅ Clear Cookies on Logout

Always call logout endpoint to clear cookies server-side:

```javascript
// Good
await fetch('https://adirdk.cloud/api/auth/logout', {
  method: 'POST',
  credentials: 'include'
});

// Then redirect
window.location.href = '/login';
```

### 5. ✅ Use HttpOnly Cookies

Backend already sets HttpOnly on sensitive tokens. Don't try to read them in JS:

```javascript
// ❌ Cannot read (HttpOnly)
const accessToken = document.cookie
  .split('; ')
  .find(c => c.startsWith('access_token='));  // Returns null

// ✅ Can read (not HttpOnly)
const csrfToken = document.cookie
  .split('; ')
  .find(c => c.startsWith('csrf_token='));  // Works
```

### 6. ✅ Implement Proper Error Handling

```javascript
try {
  const data = await apiCall('/api/admin/portfolio', {
    method: 'POST',
    body: payload
  });
  console.log('Success:', data);
} catch (error) {
  if (error.message.includes('401')) {
    // Redirect to login
    window.location.href = '/login';
  } else {
    // Show user-friendly error
    console.error('Operation failed:', error.message);
  }
}
```

---

## Testing Checklist

Before deploying to production:

- [ ] Login works and cookies appear in document.cookie
- [ ] `document.cookie` shows: `csrf_token`, `access_token`, `refresh_token`
- [ ] Creating portfolio works (CSRF validation passes)
- [ ] Updating portfolio works
- [ ] Deleting portfolio works
- [ ] Token refresh works automatically after 15 minutes
- [ ] Logout clears cookies
- [ ] 401 errors redirect to login
- [ ] HTTPS is enforced in production
- [ ] CORS errors are resolved

---

## Quick Reference

### Fetch API Template

```javascript
const response = await fetch('https://adirdk.cloud/endpoint', {
  method: 'POST',
  credentials: 'include',  // ← ALWAYS
  headers: {
    'Content-Type': 'application/json',
    'X-XSRF-TOKEN': csrfToken,  // ← For POST/PUT/DELETE
  },
  body: JSON.stringify(data)
});
```

### Axios Template

```javascript
axios.defaults.baseURL = 'https://adirdk.cloud';
axios.defaults.withCredentials = true;  // ← ALWAYS

const response = await axios.post('/api/admin/portfolio', data);
// CSRF token added automatically by interceptor
```

### React Hook Template

```javascript
const { call } = useApi();
const data = await call('/api/admin/portfolio', {
  method: 'POST',
  body: payload
});
```

---

## Support

For issues or questions:

1. Check the **Troubleshooting** section
2. Verify **credentials: 'include'** is present
3. Check **X-XSRF-TOKEN** header for CRUD operations
4. Review **Network tab** in DevTools
5. Check **Backend logs** for detailed errors

---

**Last Updated:** 2026-04-18
**Backend Version:** 1.0
**Compatibility:** All modern browsers
