# CSRF & Cookie Troubleshooting Guide

Panduan lengkap untuk debug dan fix CSRF token + cookie issues.

---

## 🔍 Diagnosis Checklist

### Step 1: Verify Backend Configuration

Akses container dan cek env:

```bash
docker exec -it portfolio-backend-app-1 sh
cat .env
```

**Expected:**
```bash
COOKIE_DOMAIN=.adirdk.cloud
COOKIE_SECURE=true
COOKIE_SAMESITE=None
ALLOWED_ORIGINS=https://dashboard.adirdk.com,https://adirdk.cloud
```

If not set correctly:
```bash
# Edit in docker
vi .env

# Add/fix:
COOKIE_DOMAIN=.adirdk.cloud
COOKIE_SECURE=true
COOKIE_SAMESITE=None

# Restart
docker restart portfolio-backend-app-1
```

---

### Step 2: Test with Postman

#### Login Request

```
Method: POST
URL: https://adirdk.cloud/api/auth/login

Headers:
  Content-Type: application/json

Body (JSON):
{
  "username": "admin",
  "password": "your_password"
}

Settings:
  ✅ Cookies: Enabled (default in Postman)
```

**Expected Response:**
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

**Check Response Headers → Cookies:**
```
Set-Cookie: access_token=...; Domain=.adirdk.cloud; Path=/; SameSite=None; Secure; HttpOnly
Set-Cookie: csrf_token=...; Domain=.adirdk.cloud; Path=/; SameSite=None; Secure
Set-Cookie: refresh_token=...; Domain=.adirdk.cloud; Path=/; SameSite=None; Secure; HttpOnly
Set-Cookie: csrf_token_sig=...; Domain=.adirdk.cloud; Path=/; SameSite=None; Secure; HttpOnly
```

---

#### Get Current User

```
Method: GET
URL: https://adirdk.cloud/api/auth/me

Headers: (none needed)

Cookies: Should be automatically included from login
```

**Check Cookies tab in Postman** - should see all 4 cookies from login.

---

#### Create Portfolio (CRUD)

```
Method: POST
URL: https://adirdk.cloud/api/admin/portfolio

Headers:
  Content-Type: application/json
  X-XSRF-TOKEN: <paste csrf_token from Cookies tab>

Body (JSON):
{
  "title": "Test Project",
  "description": "Test Description"
}

Cookies: Should be automatically included
```

**Get csrf_token from Cookies tab in Postman:**
1. After login, go to **Cookies** tab
2. Find `csrf_token` (not HttpOnly, readable)
3. Copy its value
4. Paste in `X-XSRF-TOKEN` header

---

### Step 3: Debug in Browser

Open **DevTools (F12)** → **Network Tab** → Perform login

#### Login Request

Click login request in Network tab:

**Request Headers:**
```
POST /api/auth/login
Host: adirdk.cloud
Content-Type: application/json
Origin: https://dashboard.adirdk.com
```

**Response Headers:**
```
Set-Cookie: csrf_token=abc123...; Domain=.adirdk.cloud; SameSite=None; Secure
Set-Cookie: access_token=xyz789...; Domain=.adirdk.cloud; SameSite=None; Secure; HttpOnly
Set-Cookie: refresh_token=...; ...
Set-Cookie: csrf_token_sig=...; ...
Access-Control-Allow-Credentials: true
Access-Control-Allow-Origin: https://dashboard.adirdk.com
```

**If Response Headers don't have Set-Cookie:**
```
❌ Problem: Backend not setting cookies
→ Check docker-compose.prod.yml COOKIE_* settings
→ Verify ALLOWED_ORIGINS includes dashboard.adirdk.com
```

---

#### After Login - Check Cookies Appear

```javascript
// In browser console
document.cookie
// Expected output: "csrf_token=abc...; access_token=xyz...; ..."
```

**If empty:**
```
❌ Problem: Cookies not sent to browser
→ Check if login request headers had: credentials: 'include'
→ Check CORS headers had: Access-Control-Allow-Credentials: true
```

---

#### CRUD Request

Click portfolio creation request in Network tab:

**Request Headers:**
```
POST /api/admin/portfolio
Host: adirdk.cloud
Content-Type: application/json
X-XSRF-TOKEN: abc123...  ← Must be present!
Cookie: csrf_token=...; access_token=...; ...  ← Sent by browser
Origin: https://dashboard.adirdk.com
```

**Response:**
```
200 OK
{
  "message": "Portfolio created successfully",
  ...
}
```

**If getting 403 Forbidden:**
```
❌ CSRF token header (X-XSRF-TOKEN) required
→ Add X-XSRF-TOKEN header to request
→ Value from: document.cookie where name starts with 'csrf_token='
```

**If getting 401 Unauthorized:**
```
❌ Invalid/expired access_token
→ Try GET /api/auth/me first to verify token is valid
→ If fails, login again
```

---

## 🐛 Common Issues & Fixes

### Issue 1: "CSRF token header (X-XSRF-TOKEN) required"

**Symptom:** 403 error on POST/PUT/DELETE requests

**Root Cause:** Missing `X-XSRF-TOKEN` header

**Debug:**
```bash
# Check if frontend is sending header
# DevTools → Network → CRUD request → Request Headers → Look for X-XSRF-TOKEN
```

**Fix:**

**Vanilla JS:**
```javascript
const csrfToken = document.cookie
  .split('; ')
  .find(c => c.startsWith('csrf_token='))
  ?.split('=')[1];

fetch('https://adirdk.cloud/api/admin/portfolio', {
  method: 'POST',
  credentials: 'include',
  headers: {
    'X-XSRF-TOKEN': csrfToken  // ← Add this
  }
})
```

**Axios:**
```javascript
// Add to request interceptor
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
```

---

### Issue 2: Cookies not in document.cookie

**Symptom:** `document.cookie` returns empty or missing access_token/csrf_token

**Root Cause:** Missing `credentials: 'include'`

**Debug:**
```javascript
// Check cookies were set in response
// DevTools → F12 → Network → login request → Response Headers
// Should see Set-Cookie headers
```

**Fix:**

**Vanilla JS:**
```javascript
// ❌ Wrong
fetch('https://adirdk.cloud/api/auth/login', {
  method: 'POST'
})

// ✅ Correct
fetch('https://adirdk.cloud/api/auth/login', {
  method: 'POST',
  credentials: 'include'  // ← Add this!
})
```

**Axios:**
```javascript
// Global setting
axios.defaults.withCredentials = true;

// Or per request
axios.post('/api/auth/login', data, {
  withCredentials: true  // ← Add this!
})
```

**React:**
```javascript
// In fetch hook
const response = await fetch(url, {
  credentials: 'include'  // ← Add this!
})
```

---

### Issue 3: 401 Unauthorized on GET /api/auth/me

**Symptom:** Login succeeds, but immediately get 401 on auth/me

**Root Cause:** Cookies not being sent in subsequent requests

**Debug:**
```
DevTools → Network → auth/me request
Check Request Headers:
  ✅ Should have: Cookie: access_token=...; csrf_token=...
  ❌ If missing: credentials: 'include' not set
```

**Fix:** Add `credentials: 'include'` to ALL requests

```javascript
// Template
fetch('any-url', {
  credentials: 'include'  // ← ALWAYS add this!
})
```

---

### Issue 4: "Cookie rejected because it is in a cross-site context"

**Symptom:** Browser warning in console, cookies not working

**Root Cause:** Cross-site request without proper SameSite configuration

**Debug:**
```
DevTools → Console
Look for: "Cookie rejected because it is in a cross-site context"
```

**Verify Setup:**

```bash
# Check backend env
docker exec portfolio-backend-app-1 cat .env | grep COOKIE

# Must be:
# COOKIE_SECURE=true
# COOKIE_SAMESITE=None
# COOKIE_DOMAIN=.adirdk.cloud
```

**Fix:** Update docker-compose.prod.yml

```yaml
environment:
  - COOKIE_DOMAIN=.adirdk.cloud
  - COOKIE_SECURE=true
  - COOKIE_SAMESITE=None
```

Then restart:
```bash
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d
```

---

### Issue 5: CORS Error "Access-Control-Allow-Credentials"

**Symptom:** CORS error in browser console about credentials

**Root Cause:** Backend not allowing credentials for CORS

**Debug:**
```
DevTools → Console → Check for CORS error message
DevTools → Network → Response Headers → Look for:
  Access-Control-Allow-Credentials: true
```

**Fix:** Verify backend router.go has:

```go
r.Use(cors.New(cors.Config{
  AllowCredentials: true,  // ← Must be true
  AllowOrigins: allowedOrigins,
}))
```

If not present, update and rebuild backend.

---

### Issue 6: Token Expires and Can't Refresh

**Symptom:** Works for 15 minutes, then 401 Unauthorized

**Root Cause:** Token expired, refresh endpoint not called or failed

**Debug:**
```
1. Set timer for 15 minutes
2. Try API call after 15 min
3. Should auto-refresh and retry
4. Check Network → POST /api/auth/refresh
```

**Fix:** Implement auto-refresh interceptor

**Vanilla JS:**
```javascript
export async function apiCall(endpoint, options = {}) {
  let response = await fetch(url, options);
  
  if (response.status === 401) {
    // Try refresh
    const refreshResp = await fetch('https://adirdk.cloud/api/auth/refresh', {
      method: 'POST',
      credentials: 'include'
    });
    
    if (refreshResp.ok) {
      // Retry original request
      response = await fetch(url, options);
    } else {
      // Go to login
      window.location.href = '/login';
    }
  }
  
  return response.json();
}
```

**Axios:**
```javascript
api.interceptors.response.use(
  response => response,
  async error => {
    if (error.response?.status === 401 && !error.config._retry) {
      error.config._retry = true;
      await api.post('/api/auth/refresh');
      return api(error.config);
    }
    window.location.href = '/login';
    return Promise.reject(error);
  }
);
```

---

## 🧪 Manual Testing Steps

### Full Flow Test

```javascript
// 1. Login
fetch('https://adirdk.cloud/api/auth/login', {
  method: 'POST',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'admin', password: 'xxx' })
})
.then(r => r.json())
.then(d => {
  console.log('✅ Login response:', d);
  console.log('✅ Cookies set:', document.cookie);
  return d;
});

// 2. Get current user
// Wait ~1 second then:
fetch('https://adirdk.cloud/api/auth/me', {
  method: 'GET',
  credentials: 'include'
})
.then(r => r.json())
.then(d => console.log('✅ Current user:', d));

// 3. Create portfolio
// Get CSRF token first:
const csrfToken = document.cookie.split('; ').find(c => c.startsWith('csrf_token='))?.split('=')[1];

fetch('https://adirdk.cloud/api/admin/portfolio', {
  method: 'POST',
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json',
    'X-XSRF-TOKEN': csrfToken
  },
  body: JSON.stringify({
    title: 'Test Project',
    description: 'Test Description'
  })
})
.then(r => r.json())
.then(d => console.log('✅ Created portfolio:', d));

// 4. Logout
fetch('https://adirdk.cloud/api/auth/logout', {
  method: 'POST',
  credentials: 'include'
})
.then(r => r.json())
.then(d => {
  console.log('✅ Logged out:', d);
  console.log('✅ Cookies after logout:', document.cookie);
});
```

---

## 📊 Network Tab Checklist

For each request type, verify in Network tab:

### Login Request ✅
- [ ] Method: POST
- [ ] URL: /api/auth/login
- [ ] Request Headers: `Content-Type: application/json`
- [ ] Response Headers: `Set-Cookie` (4 cookies)
- [ ] Response Status: 200 OK

### Auth/Me Request ✅
- [ ] Method: GET
- [ ] URL: /api/auth/me
- [ ] Request Headers: `Cookie: ...` (from login)
- [ ] Response Status: 200 OK

### CRUD Create ✅
- [ ] Method: POST
- [ ] URL: /api/admin/portfolio
- [ ] Request Headers: `X-XSRF-TOKEN: ...`
- [ ] Request Headers: `Cookie: ...`
- [ ] Response Status: 201 Created

### CRUD Update ✅
- [ ] Method: PUT
- [ ] URL: /api/admin/portfolio/:id
- [ ] Request Headers: `X-XSRF-TOKEN: ...`
- [ ] Request Headers: `Cookie: ...`
- [ ] Response Status: 200 OK

### CRUD Delete ✅
- [ ] Method: DELETE
- [ ] URL: /api/admin/portfolio/:id
- [ ] Request Headers: `X-XSRF-TOKEN: ...`
- [ ] Request Headers: `Cookie: ...`
- [ ] Response Status: 200 OK

---

## 🔐 Quick Security Verification

```bash
# 1. Verify backend uses HTTPS
curl -v https://adirdk.cloud/health
# Should NOT have any warning about certificate

# 2. Verify CORS is correct
curl -v -X OPTIONS https://adirdk.cloud/api/admin/portfolio \
  -H "Origin: https://dashboard.adirdk.com" \
  -H "Access-Control-Request-Method: POST"
# Should have: Access-Control-Allow-Origin: https://dashboard.adirdk.com
# Should have: Access-Control-Allow-Credentials: true

# 3. Verify cookie attributes
curl -X POST https://adirdk.cloud/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"xxx"}' \
  -v | grep -i "set-cookie"
# Should show: SameSite=None; Secure
```

---

## 📝 Log Analysis

### Backend Logs

```bash
# View logs
docker logs -f portfolio-backend-app-1

# Look for errors:
# CSRF token signature invalid
# → X-XSRF-TOKEN header doesn't match cookie
# → Fix: Send header with value from csrf_token cookie

# CSRF token cookie missing
# → credentials: 'include' not set on request
# → Fix: Add credentials: 'include' to fetch options

# Authorization header not found
# → access_token cookie not being sent
# → Fix: Add credentials: 'include' to fetch options
```

---

## ✅ Pre-Deployment Checklist

- [ ] Login works and sets all 4 cookies
- [ ] `document.cookie` shows csrf_token, access_token, refresh_token
- [ ] Can create/update/delete portfolio
- [ ] Token auto-refreshes after 15 minutes
- [ ] 401 errors redirect to login
- [ ] HTTPS is enforced (no mixed content warnings)
- [ ] CORS headers are correct (checked with curl)
- [ ] No console errors in DevTools
- [ ] All Network requests have correct headers
- [ ] Logout clears cookies and redirects

---

## 📞 Still Having Issues?

1. **Check Postman works** → If Postman works, issue is in frontend
2. **Check backend logs** → `docker logs portfolio-backend-app-1`
3. **Check DevTools** → Network tab + Console
4. **Verify .env** → `docker exec portfolio-backend-app-1 cat .env`
5. **Check CORS** → Use curl to test

---

**Last Updated:** 2026-04-18
