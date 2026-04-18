# Frontend Implementation Examples

Contoh kode siap pakai untuk berbagai use case.

---

## 📦 Vanilla JavaScript (Fetch API)

### Complete API Client

```javascript
// src/api/client.js

const API_BASE_URL = process.env.REACT_APP_API_URL || 'https://adirdk.cloud';
const CSRF_COOKIE_NAME = 'csrf_token';
const CSRF_HEADER_NAME = 'X-XSRF-TOKEN';

/**
 * Extract cookie value by name
 */
function getCookie(name) {
  const cookies = document.cookie.split('; ');
  const cookie = cookies.find(c => c.startsWith(name + '='));
  return cookie ? cookie.split('=')[1] : null;
}

/**
 * Generic API call wrapper with CSRF and auto-refresh
 */
export async function apiCall(endpoint, options = {}) {
  const {
    method = 'GET',
    body = null,
    headers = {},
    skipRefresh = false,
  } = options;

  // Add CSRF token for state-changing requests
  const csrfToken = ['POST', 'PUT', 'DELETE'].includes(method)
    ? getCookie(CSRF_COOKIE_NAME)
    : null;

  if (csrfToken && !headers[CSRF_HEADER_NAME]) {
    headers[CSRF_HEADER_NAME] = csrfToken;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method,
    credentials: 'include',  // CRITICAL: Include cookies
    headers: {
      'Content-Type': 'application/json',
      ...headers,
    },
    ...(body && { body: JSON.stringify(body) }),
  });

  // Handle 401: Try token refresh
  if (response.status === 401 && !skipRefresh) {
    const refreshed = await refreshAccessToken();
    if (refreshed) {
      // Retry original request
      return apiCall(endpoint, { ...options, skipRefresh: true });
    }
    // Redirect to login
    handleUnauthorized();
    return null;
  }

  // Handle other errors
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    const error = new Error(errorData.error || `HTTP ${response.status}`);
    error.status = response.status;
    throw error;
  }

  return response.json();
}

/**
 * Auth: Login
 */
export async function login(username, password) {
  return apiCall('/api/auth/login', {
    method: 'POST',
    body: { username, password },
  });
}

/**
 * Auth: Logout
 */
export async function logout() {
  await apiCall('/api/auth/logout', { method: 'POST' });
  handleUnauthorized();
}

/**
 * Auth: Get current user
 */
export async function getMe() {
  return apiCall('/api/auth/me');
}

/**
 * Auth: Refresh access token
 */
export async function refreshAccessToken() {
  try {
    await apiCall('/api/auth/refresh', {
      method: 'POST',
      skipRefresh: true,  // Prevent infinite loop
    });
    return true;
  } catch (error) {
    console.error('Token refresh failed:', error);
    return false;
  }
}

/**
 * CRUD: Portfolio
 */
export const Portfolio = {
  create: (data) =>
    apiCall('/api/admin/portfolio', { method: 'POST', body: data }),

  update: (id, data) =>
    apiCall(`/api/admin/portfolio/${id}`, { method: 'PUT', body: data }),

  delete: (id) =>
    apiCall(`/api/admin/portfolio/${id}`, { method: 'DELETE' }),

  deleteMedia: (portfolioId, mediaId) =>
    apiCall(`/api/admin/portfolio-media/${portfolioId}/${mediaId}`, {
      method: 'DELETE',
    }),
};

/**
 * CRUD: Skills
 */
export const Skills = {
  create: (data) =>
    apiCall('/api/admin/skills', { method: 'POST', body: data }),

  update: (id, data) =>
    apiCall(`/api/admin/skills/${id}`, { method: 'PUT', body: data }),

  delete: (id) =>
    apiCall(`/api/admin/skills/${id}`, { method: 'DELETE' }),
};

/**
 * CRUD: Qualifications
 */
export const Qualifications = {
  create: (data) =>
    apiCall('/api/admin/qualifications', { method: 'POST', body: data }),

  update: (id, data) =>
    apiCall(`/api/admin/qualifications/${id}`, { method: 'PUT', body: data }),

  delete: (id) =>
    apiCall(`/api/admin/qualifications/${id}`, { method: 'DELETE' }),
};

/**
 * Handle unauthorized (401) response
 */
function handleUnauthorized() {
  // Clear local storage
  localStorage.removeItem('user');

  // Dispatch custom event (if using events)
  window.dispatchEvent(new CustomEvent('logout'));

  // Redirect to login
  window.location.href = '/login';
}
```

### Usage Example

```javascript
// src/main.js

import * as api from './api/client.js';

// Login
async function handleLogin() {
  try {
    const response = await api.login('admin', 'password');
    console.log('Logged in:', response.user);
    localStorage.setItem('user', JSON.stringify(response.user));
  } catch (error) {
    console.error('Login failed:', error.message);
  }
}

// Create portfolio
async function handleCreatePortfolio() {
  try {
    const result = await api.Portfolio.create({
      title: 'My Awesome Project',
      description: 'This is an awesome project',
    });
    console.log('Portfolio created:', result);
  } catch (error) {
    console.error('Failed to create portfolio:', error.message);
  }
}

// Update portfolio
async function handleUpdatePortfolio(id) {
  try {
    const result = await api.Portfolio.update(id, {
      title: 'Updated Title',
      description: 'Updated Description',
    });
    console.log('Portfolio updated:', result);
  } catch (error) {
    console.error('Failed to update portfolio:', error.message);
  }
}

// Delete portfolio
async function handleDeletePortfolio(id) {
  try {
    await api.Portfolio.delete(id);
    console.log('Portfolio deleted');
  } catch (error) {
    console.error('Failed to delete portfolio:', error.message);
  }
}

// Logout
async function handleLogout() {
  try {
    await api.logout();
  } catch (error) {
    console.error('Logout failed:', error.message);
  }
}
```

---

## ⚛️ React (TypeScript)

### API Hook

```typescript
// src/hooks/useApi.ts

import { useState, useCallback } from 'react';

interface ApiOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  body?: Record<string, any>;
  headers?: Record<string, string>;
  skipRefresh?: boolean;
}

interface ApiError extends Error {
  status?: number;
}

export function useApi() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const apiBaseUrl = process.env.REACT_APP_API_URL || 'https://adirdk.cloud';

  const getCookie = useCallback((name: string): string | null => {
    const cookies = document.cookie.split('; ');
    const cookie = cookies.find(c => c.startsWith(name + '='));
    return cookie ? cookie.split('=')[1] : null;
  }, []);

  const handleUnauthorized = useCallback(() => {
    localStorage.removeItem('user');
    window.location.href = '/login';
  }, []);

  const call = useCallback(
    async (endpoint: string, options: ApiOptions = {}) => {
      const {
        method = 'GET',
        body = null,
        headers = {},
        skipRefresh = false,
      } = options;

      setLoading(true);
      setError(null);

      try {
        const csrfToken = ['POST', 'PUT', 'DELETE'].includes(method)
          ? getCookie('csrf_token')
          : null;

        const response = await fetch(`${apiBaseUrl}${endpoint}`, {
          method,
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
            ...(csrfToken && { 'X-XSRF-TOKEN': csrfToken }),
            ...headers,
          },
          ...(body && { body: JSON.stringify(body) }),
        });

        if (response.status === 401 && !skipRefresh) {
          const refreshResp = await fetch(`${apiBaseUrl}/api/auth/refresh`, {
            method: 'POST',
            credentials: 'include',
          });

          if (refreshResp.ok) {
            return call(endpoint, { ...options, skipRefresh: true });
          } else {
            handleUnauthorized();
            return null;
          }
        }

        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}));
          throw new Error(errorData.error || `HTTP ${response.status}`);
        }

        const data = await response.json();
        return data;
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [getCookie, handleUnauthorized, apiBaseUrl]
  );

  return { call, loading, error };
}
```

### Auth Context

```typescript
// src/context/AuthContext.tsx

import React, { createContext, useContext, useState, useCallback } from 'react';
import { useApi } from '../hooks/useApi';

interface User {
  id: number;
  username: string;
  role: string;
}

interface AuthContextType {
  user: User | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  getMe: () => Promise<User | null>;
  loading: boolean;
  error: string | null;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { call, loading, error } = useApi();
  const [user, setUser] = useState<User | null>(null);

  const login = useCallback(
    async (username: string, password: string) => {
      const data = await call('/api/auth/login', {
        method: 'POST',
        body: { username, password },
      });
      setUser(data.user);
      localStorage.setItem('user', JSON.stringify(data.user));
    },
    [call]
  );

  const logout = useCallback(async () => {
    await call('/api/auth/logout', { method: 'POST' });
    setUser(null);
    localStorage.removeItem('user');
    window.location.href = '/login';
  }, [call]);

  const getMe = useCallback(async () => {
    try {
      const data = await call('/api/auth/me');
      setUser(data.user);
      localStorage.setItem('user', JSON.stringify(data.user));
      return data.user;
    } catch {
      setUser(null);
      localStorage.removeItem('user');
      return null;
    }
  }, [call]);

  return (
    <AuthContext.Provider value={{ user, login, logout, getMe, loading, error }}>
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

### Login Component

```typescript
// src/pages/LoginPage.tsx

import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [localError, setLocalError] = useState('');

  const { login, loading, error } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLocalError('');

    try {
      await login(username, password);
      navigate('/dashboard');
    } catch (err) {
      setLocalError(
        err instanceof Error ? err.message : 'Login failed'
      );
    }
  };

  return (
    <div className="login-container">
      <form onSubmit={handleSubmit}>
        <h1>Login</h1>
        
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
          disabled={loading}
        />

        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          disabled={loading}
        />

        <button type="submit" disabled={loading}>
          {loading ? 'Logging in...' : 'Login'}
        </button>

        {(error || localError) && (
          <p className="error">{error || localError}</p>
        )}
      </form>
    </div>
  );
}
```

### Portfolio CRUD Component

```typescript
// src/components/PortfolioForm.tsx

import React, { useState } from 'react';
import { useApi } from '../hooks/useApi';

export function PortfolioForm() {
  const { call, loading, error } = useApi();
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSuccessMessage('');

    try {
      const result = await call('/api/admin/portfolio', {
        method: 'POST',
        body: { title, description },
      });

      setSuccessMessage('Portfolio created successfully!');
      setTitle('');
      setDescription('');

      // Refresh portfolio list
      window.dispatchEvent(new CustomEvent('portfolioCreated'));
    } catch (err) {
      // Error handled by useApi
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>Create Portfolio</h2>

      <input
        type="text"
        placeholder="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        required
        disabled={loading}
      />

      <textarea
        placeholder="Description"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        required
        disabled={loading}
      />

      <button type="submit" disabled={loading}>
        {loading ? 'Creating...' : 'Create'}
      </button>

      {error && <p className="error">{error}</p>}
      {successMessage && <p className="success">{successMessage}</p>}
    </form>
  );
}
```

---

## 🟢 Vue 3 (Composition API)

### API Plugin

```typescript
// src/api.ts

import { ref, Ref } from 'vue';

export interface UseApiReturn {
  call: (endpoint: string, options?: any) => Promise<any>;
  loading: Ref<boolean>;
  error: Ref<string | null>;
}

export function useApi(): UseApiReturn {
  const loading = ref(false);
  const error = ref<string | null>(null);

  const apiBaseUrl = import.meta.env.VITE_API_URL || 'https://adirdk.cloud';

  const getCookie = (name: string): string | null => {
    const cookies = document.cookie.split('; ');
    const cookie = cookies.find(c => c.startsWith(name + '='));
    return cookie ? cookie.split('=')[1] : null;
  };

  const handleUnauthorized = () => {
    localStorage.removeItem('user');
    window.location.href = '/login';
  };

  const call = async (
    endpoint: string,
    options: any = {}
  ): Promise<any> => {
    const { method = 'GET', body = null, headers = {}, skipRefresh = false } = options;

    loading.value = true;
    error.value = null;

    try {
      const csrfToken =
        ['POST', 'PUT', 'DELETE'].includes(method) ? getCookie('csrf_token') : null;

      const response = await fetch(`${apiBaseUrl}${endpoint}`, {
        method,
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          ...(csrfToken && { 'X-XSRF-TOKEN': csrfToken }),
          ...headers,
        },
        ...(body && { body: JSON.stringify(body) }),
      });

      if (response.status === 401 && !skipRefresh) {
        const refreshResp = await fetch(`${apiBaseUrl}/api/auth/refresh`, {
          method: 'POST',
          credentials: 'include',
        });

        if (refreshResp.ok) {
          return call(endpoint, { ...options, skipRefresh: true });
        } else {
          handleUnauthorized();
          return null;
        }
      }

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP ${response.status}`);
      }

      return await response.json();
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return { call, loading, error };
}
```

### Auth Store (Pinia)

```typescript
// src/stores/auth.ts

import { defineStore } from 'pinia';
import { useApi } from '../api';

interface User {
  id: number;
  username: string;
  role: string;
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
  }),

  actions: {
    async login(username: string, password: string) {
      const { call } = useApi();
      const data = await call('/api/auth/login', {
        method: 'POST',
        body: { username, password },
      });
      this.user = data.user;
      localStorage.setItem('user', JSON.stringify(data.user));
      return data;
    },

    async logout() {
      const { call } = useApi();
      await call('/api/auth/logout', { method: 'POST' });
      this.user = null;
      localStorage.removeItem('user');
      window.location.href = '/login';
    },

    async getMe() {
      const { call } = useApi();
      try {
        const data = await call('/api/auth/me');
        this.user = data.user;
        localStorage.setItem('user', JSON.stringify(data.user));
        return data.user;
      } catch {
        this.user = null;
        localStorage.removeItem('user');
        return null;
      }
    },
  },
});
```

### Login Component

```vue
// src/views/LoginView.vue

<template>
  <div class="login-container">
    <form @submit.prevent="handleLogin">
      <h1>Login</h1>
      
      <input
        v-model="username"
        type="text"
        placeholder="Username"
        required
        :disabled="loading"
      />

      <input
        v-model="password"
        type="password"
        placeholder="Password"
        required
        :disabled="loading"
      />

      <button type="submit" :disabled="loading">
        {{ loading ? 'Logging in...' : 'Login' }}
      </button>

      <p v-if="error" class="error">{{ error }}</p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useApi } from '@/api';

const username = ref('');
const password = ref('');
const authStore = useAuthStore();
const router = useRouter();
const { loading, error } = useApi();

async function handleLogin() {
  try {
    await authStore.login(username.value, password.value);
    router.push('/dashboard');
  } catch (err) {
    console.error('Login failed:', err);
  }
}
</script>
```

---

## 🔴 Angular (HttpClientModule)

### Interceptor untuk CSRF

```typescript
// src/app/interceptors/csrf.interceptor.ts

import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
  HttpErrorResponse,
} from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, switchMap } from 'rxjs/operators';
import { AuthService } from '../services/auth.service';

@Injectable()
export class CsrfInterceptor implements HttpInterceptor {
  constructor(private authService: AuthService) {}

  intercept(
    request: HttpRequest<unknown>,
    next: HttpHandler
  ): Observable<HttpEvent<unknown>> {
    // Add CSRF token for state-changing requests
    if (['POST', 'PUT', 'DELETE'].includes(request.method)) {
      const csrfToken = this.getCookie('csrf_token');
      if (csrfToken) {
        request = request.clone({
          setHeaders: {
            'X-XSRF-TOKEN': csrfToken,
          },
        });
      }
    }

    return next.handle(request).pipe(
      catchError((error: HttpErrorResponse) => {
        if (error.status === 401) {
          return this.authService.refreshToken().pipe(
            switchMap(() => next.handle(request)),
            catchError(() => {
              this.authService.logout();
              return throwError(() => error);
            })
          );
        }
        return throwError(() => error);
      })
    );
  }

  private getCookie(name: string): string | null {
    const cookies = document.cookie.split('; ');
    const cookie = cookies.find(c => c.startsWith(name + '='));
    return cookie ? cookie.split('=')[1] : null;
  }
}
```

### App Configuration

```typescript
// src/app/app.config.ts

import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { CsrfInterceptor } from './interceptors/csrf.interceptor';
import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    importProvidersFrom(HttpClientModule),
    {
      provide: HTTP_INTERCEPTORS,
      useClass: CsrfInterceptor,
      multi: true,
    },
  ],
};
```

---

## 📚 Summary

**Vanilla JS:** Gunakan `apiCall()` wrapper function
**React:** Gunakan `useApi()` hook + `useAuth()` context
**Vue:** Gunakan `useApi()` composable + Pinia store
**Angular:** Gunakan HttpClient dengan interceptor

Semua contoh sudah include:
- ✅ `credentials: 'include'`
- ✅ X-XSRF-TOKEN header untuk CRUD
- ✅ Auto token refresh pada 401
- ✅ Error handling
- ✅ Loading states
