# Frontend Integration Documentation

**Dokumentasi lengkap untuk mengintegrasikan frontend dengan backend authentication.**

Direktori ini berisi panduan komprehensif tentang cara menggunakan sistem CSRF + cookie authentication backend dengan frontend apapun.

---

## 📚 Dokumentasi yang Tersedia

### 1. **[FRONTEND-SETUP.md](./FRONTEND-SETUP.md)** - Panduan Setup Lengkap ⭐

**Baca ini dulu!** Panduan komprehensif yang mencakup:
- Overview dan architecture
- Setup untuk Vanilla JS, React, Vue, Angular, Axios
- API integration guide
- Complete authentication flow
- Security best practices
- Testing checklist

**Waktu baca:** 30-45 menit  
**Untuk:** Semua developer frontend

---

### 2. **[FRONTEND-EXAMPLES.md](./FRONTEND-EXAMPLES.md)** - Contoh Kode Siap Pakai

Contoh implementasi lengkap untuk setiap framework:

- **Vanilla JavaScript** - API client dengan auto-refresh
- **React** - Hooks, Context, Components
- **Vue 3** - Composition API, Pinia store
- **Angular** - HttpClient, Interceptors
- **Axios** - Global setup & interceptors

**Waktu baca:** 20 menit  
**Untuk:** Developer yang ingin copy-paste kode

---

### 3. **[CSRF-TROUBLESHOOTING.md](./CSRF-TROUBLESHOOTING.md)** - Debug & Troubleshooting

Panduan debug untuk masalah yang paling umum:

- Diagnosis checklist
- Testing dengan Postman
- Debug di browser DevTools
- Common issues & solutions
- Network tab verification
- Manual testing steps

**Waktu baca:** 15-20 menit  
**Untuk:** Ketika ada error atau masalah

---

### 4. **[API.md](./API.md)** - API Endpoints Documentation

Dokumentasi detail tentang semua endpoints yang tersedia (backend).

---

### 5. **[visitor-tracking.md](./visitor-tracking.md)** - Visitor Tracking API

Dokumentasi tentang visitor tracking API (public endpoint).

---

## 🚀 Quick Start

### Backend Setup

```bash
# 1. Verify environment variables
docker exec portfolio-backend-app-1 cat .env

# Harus ada:
COOKIE_DOMAIN=.adirdk.cloud
COOKIE_SECURE=true
COOKIE_SAMESITE=None
ALLOWED_ORIGINS=https://dashboard.adirdk.com,...
```

### Frontend Setup (Choose One)

#### Vanilla JavaScript
1. Copy kode dari [FRONTEND-EXAMPLES.md - Vanilla JS section](./FRONTEND-EXAMPLES.md#-vanilla-javascript-fetch-api)
2. Buat file `src/api/client.js`
3. Gunakan `apiCall()` di aplikasi Anda

#### React
1. Copy hooks dari [FRONTEND-EXAMPLES.md - React section](./FRONTEND-EXAMPLES.md#-react-typescript)
2. Buat Context provider untuk auth
3. Wrap aplikasi dengan `<AuthProvider>`

#### Vue 3
1. Copy composable dari [FRONTEND-EXAMPLES.md - Vue section](./FRONTEND-EXAMPLES.md#-vue-3-composition-api)
2. Setup Pinia store
3. Install provider di `main.ts`

#### Angular
1. Copy interceptor dari [FRONTEND-EXAMPLES.md - Angular section](./FRONTEND-EXAMPLES.md#-angular-httpclientmodule)
2. Configure di `app.config.ts`
3. Gunakan HttpClient normally

---

## 📋 Implementation Checklist

Sebelum go-live, verifikasi:

### Backend Configuration
- [ ] `COOKIE_DOMAIN=.adirdk.cloud` (dengan titik di depan)
- [ ] `COOKIE_SECURE=true`
- [ ] `COOKIE_SAMESITE=None`
- [ ] `ALLOWED_ORIGINS` include frontend domain
- [ ] HTTPS enabled di production

### Frontend Code
- [ ] ✅ `credentials: 'include'` di semua fetch calls
- [ ] ✅ `withCredentials: true` untuk axios
- [ ] ✅ `X-XSRF-TOKEN` header di POST/PUT/DELETE
- [ ] ✅ Auto-refresh token pada 401 error
- [ ] ✅ Redirect to login pada 401
- [ ] ✅ Error handling di semua API calls

### Testing
- [ ] Login works, cookies appear in `document.cookie`
- [ ] Can create/update/delete data
- [ ] Token refresh works after 15 minutes
- [ ] Logout clears cookies
- [ ] Postman works (to verify backend is OK)
- [ ] No CORS errors in console
- [ ] No auth errors in Network tab

---

## 🎯 Most Common Issues & Fixes

### ❌ "CSRF token header required"

**Cause:** Missing `X-XSRF-TOKEN` header on POST/PUT/DELETE

**Fix:**
```javascript
const csrfToken = document.cookie
  .split('; ')
  .find(c => c.startsWith('csrf_token='))
  ?.split('=')[1];

fetch(url, {
  headers: {
    'X-XSRF-TOKEN': csrfToken  // ← Add this
  }
})
```

→ See [CSRF-TROUBLESHOOTING.md#issue-1](./CSRF-TROUBLESHOOTING.md#issue-1-csrf-token-header-x-xsrf-token-required)

---

### ❌ Cookies not in document.cookie

**Cause:** Missing `credentials: 'include'`

**Fix:**
```javascript
fetch(url, {
  credentials: 'include'  // ← Add this!
})
```

Or for axios:
```javascript
axios.defaults.withCredentials = true;
```

→ See [CSRF-TROUBLESHOOTING.md#issue-2](./CSRF-TROUBLESHOOTING.md#issue-2-cookies-not-in-documentcookie)

---

### ❌ 401 Unauthorized

**Cause:** Cookies not being sent in requests

**Fix:** Add `credentials: 'include'` to ALL requests

→ See [CSRF-TROUBLESHOOTING.md#issue-3](./CSRF-TROUBLESHOOTING.md#issue-3-401-unauthorized-on-get-apiauth-me)

---

## 🧪 Testing dengan Postman

Fastest way to verify backend works:

1. **Login:**
   ```
   POST https://adirdk.cloud/api/auth/login
   Body: {"username":"admin","password":"xxx"}
   ```

2. **Check cookies in Postman → Cookies tab**

3. **Create Portfolio:**
   ```
   POST https://adirdk.cloud/api/admin/portfolio
   Header: X-XSRF-TOKEN: <csrf_token from cookies>
   Body: {"title":"Test","description":"Desc"}
   ```

**Jika Postman works tapi frontend tidak:**
→ Issue pasti di frontend (biasanya missing `credentials: 'include'`)

---

## 🔍 Quick Debug in Browser

```javascript
// 1. Check cookies
console.log('Cookies:', document.cookie);

// 2. Extract CSRF token
const csrf = document.cookie
  .split('; ')
  .find(c => c.startsWith('csrf_token='))
  ?.split('=')[1];
console.log('CSRF Token:', csrf);

// 3. Test API call
fetch('https://adirdk.cloud/api/auth/me', {
  credentials: 'include'
})
.then(r => r.json())
.then(d => console.log('User:', d));
```

**Lalu check Network tab untuk:**
- Request headers: `Cookie: ...` present?
- Response headers: `Set-Cookie: ...` pada login?

→ See [CSRF-TROUBLESHOOTING.md#-debug-in-browser](./CSRF-TROUBLESHOOTING.md#-debug-in-browser)

---

## 📚 Reading Order (Recommended)

Untuk developer baru:

1. **Start here:** [FRONTEND-SETUP.md - Overview](./FRONTEND-SETUP.md#overview) (5 min)
2. **Understand flow:** [FRONTEND-SETUP.md - Architecture](./FRONTEND-SETUP.md#architecture) (10 min)
3. **Choose framework:** [FRONTEND-SETUP.md - Setup by Framework](./FRONTEND-SETUP.md#setup-by-framework) (15 min)
4. **Copy code:** [FRONTEND-EXAMPLES.md](./FRONTEND-EXAMPLES.md) (grab your framework section)
5. **Test it:** [CSRF-TROUBLESHOOTING.md - Manual Testing](./CSRF-TROUBLESHOOTING.md#-manual-testing-steps) (5 min)
6. **Debug if needed:** [CSRF-TROUBLESHOOTING.md - Common Issues](./CSRF-TROUBLESHOOTING.md#-common-issues--fixes)

---

## 🔐 Security Notes

✅ **Backend handles security:**
- HttpOnly cookies untuk sensitive tokens
- HMAC-SHA256 untuk CSRF token signing
- SameSite=None + Secure flag
- Double-submit cookie pattern
- Automatic token rotation

✅ **Frontend must do:**
- Use HTTPS in production
- Never store access token di localStorage (stay in cookie)
- Always validate CSRF token exists
- Handle token refresh automatically
- Redirect to login on 401

---

## 🆘 Still Having Issues?

### Step-by-step Debug Process

1. **Test backend with Postman** → Verify backend works
2. **Check browser console** → Look for CORS/auth errors
3. **Check Network tab** → Verify cookies & headers
4. **Read troubleshooting** → Find your error pattern
5. **Check docker logs** → `docker logs portfolio-backend-app-1`

```bash
# If all else fails, restart backend
docker-compose -f docker-compose.prod.yml restart app
```

---

## 📞 Support Resources

**Inside this documentation:**
- Overview → [FRONTEND-SETUP.md](./FRONTEND-SETUP.md)
- Code examples → [FRONTEND-EXAMPLES.md](./FRONTEND-EXAMPLES.md)
- Troubleshooting → [CSRF-TROUBLESHOOTING.md](./CSRF-TROUBLESHOOTING.md)
- Backend API → [API.md](./API.md)

**External resources:**
- [MDN - Cookies](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies)
- [MDN - SameSite attribute](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie/SameSite)
- [OWASP - CSRF Prevention](https://owasp.org/www-community/attacks/csrf)

---

## 📝 Version Info

- **Backend Version:** 1.0
- **Documentation Updated:** 2026-04-18
- **Supported Frameworks:** Vanilla JS, React, Vue 3, Angular, Axios
- **Node Compatibility:** 14+
- **Browser Compatibility:** All modern browsers (Chrome, Firefox, Safari, Edge)

---

## 🎓 Key Concepts

### Cookies vs Headers

| Aspect | Cookie | Header |
|--------|--------|--------|
| **Sent automatically** | ✅ Yes (if `credentials: 'include'`) | ❌ No |
| **HttpOnly** | ✅ Can be | ❌ No |
| **JS readable** | ✅ (unless HttpOnly) | N/A |
| **For auth** | ✅ access_token, refresh_token | ✅ X-XSRF-TOKEN |

### CSRF Protection Flow

```
1. Browser stores: csrf_token (readable), csrf_token_sig (HttpOnly)
2. JS reads csrf_token from document.cookie
3. JS sends csrf_token in X-XSRF-TOKEN header
4. Backend validates:
   - Token matches header & cookie
   - Signature is valid
5. Backend allows request
```

### Cross-Site vs Same-Site

```
Same-site:    dashboard.adirdk.cloud → api.adirdk.cloud ✅ Can use Lax
Cross-site:   dashboard.adirdk.com → adirdk.cloud ❌ Need None + Secure
```

---

## ✅ Pre-Deployment Checklist

```bash
# Run this before deploying to production

# 1. Verify backend config
docker exec portfolio-backend-app-1 cat .env | grep COOKIE

# 2. Test with curl
curl -X POST https://adirdk.cloud/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"xxx"}' \
  -v | grep -i "set-cookie"

# 3. Test CORS
curl -X OPTIONS https://adirdk.cloud/api/admin/portfolio \
  -H "Origin: https://dashboard.adirdk.com" \
  -H "Access-Control-Request-Method: POST" \
  -v | grep -i "access-control"

# Expected: Access-Control-Allow-Credentials: true
# Expected: Access-Control-Allow-Origin: https://dashboard.adirdk.com
```

---

## 🙋 FAQ

**Q: Harus pakai cookies apa bisa pakai Bearer token?**  
A: Bisa, tapi harus ubah backend ke Bearer token validation. Cookies lebih aman dari CSRF attacks.

**Q: Why SameSite=None?**  
A: Frontend dan backend cross-site (`dashboard.adirdk.com` vs `adirdk.cloud`). Lax hanya work same-site.

**Q: Berapa lama token expired?**  
A: Access token 15 menit, refresh token 7 hari. Auto-refresh implemented di contoh.

**Q: Bisa hide CSRF token?**  
A: Tidak, frontend perlu baca untuk kirim di header. CSRF token purposely non-HttpOnly.

**Q: Apa itu credentials: 'include'?**  
A: Flag untuk fetch/axios agar kirim & terima cookies dengan cross-site requests.

---

**Happy coding! 🚀**

For questions or improvements, check the troubleshooting guide or backend logs.
