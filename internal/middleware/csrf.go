package middleware

import (
	"backend-portfolio/internal/auth"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CSRFCookieName = "csrf_token"
	CSRFHeaderName = "X-XSRF-TOKEN"
	csrfTokenLen   = 32
)

// CSRFService manages CSRF token generation, signing, and validation.
// It uses HMAC-SHA256 with the JWT access secret as the signing key.
type CSRFService struct {
	hmacKey []byte
	cookie  *auth.CookieManager
}

// NewCSRFService creates a CSRFService.
func NewCSRFService(hmacKey []byte, cookie *auth.CookieManager) *CSRFService {
	return &CSRFService{hmacKey: hmacKey, cookie: cookie}
}

// SetCSRFCookie writes a readable CSRF token cookie (for JS) and a signed
// HttpOnly cookie (for server-side validation).
func (s *CSRFService) SetCSRFCookie(c *gin.Context) {
	token, err := generateCSRFToken()
	if err != nil {
		return
	}
	signature := s.sign(token)

	// Readable cookie — JS reads this and puts it in the X-XSRF-TOKEN header
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    token,
		Path:     "/",
		Domain:   s.cookie.Domain(),
		MaxAge:   s.cookie.AccessMaxAge(),
		Secure:   s.cookie.Secure(),
		HttpOnly: false, // JS must read this
		SameSite: s.cookie.SameSite(),
	})

	// Signed cookie — HttpOnly, used for server-side comparison
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     CSRFCookieName + "_sig",
		Value:    signature,
		Path:     "/",
		Domain:   s.cookie.Domain(),
		MaxAge:   s.cookie.AccessMaxAge(),
		Secure:   s.cookie.Secure(),
		HttpOnly: true,
		SameSite: s.cookie.SameSite(),
	})
}

// ClearCSRFCookie removes CSRF cookies.
func (s *CSRFService) ClearCSRFCookie(c *gin.Context) {
	for _, name := range []string{CSRFCookieName, CSRFCookieName + "_sig"} {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Domain:   s.cookie.Domain(),
			MaxAge:   -1,
			Secure:   s.cookie.Secure(),
			HttpOnly: name != CSRFCookieName,
			SameSite: s.cookie.SameSite(),
		})
	}
}

// Middleware returns a Gin middleware that validates the Signed Double Submit Cookie.
func (s *CSRFService) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			c.Next()
			return
		}

		// Skip CSRF for Bearer-token clients (not vulnerable to CSRF)
		if c.GetHeader("Authorization") != "" {
			c.Next()
			return
		}

		cookieToken, err := c.Cookie(CSRFCookieName)
		if err != nil || cookieToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token cookie missing"})
			c.Abort()
			return
		}

		cookieSig, err := c.Cookie(CSRFCookieName + "_sig")
		if err != nil || cookieSig == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF signature missing"})
			c.Abort()
			return
		}

		if !s.verify(cookieToken, cookieSig) {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token signature invalid"})
			c.Abort()
			return
		}

		headerToken := c.GetHeader(CSRFHeaderName)
		if headerToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token header (X-XSRF-TOKEN) required"})
			c.Abort()
			return
		}

		if !hmac.Equal([]byte(headerToken), []byte(cookieToken)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *CSRFService) sign(token string) string {
	mac := hmac.New(sha256.New, s.hmacKey)
	mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *CSRFService) verify(token, signature string) bool {
	expected := s.sign(token)
	return hmac.Equal([]byte(expected), []byte(signature))
}

func generateCSRFToken() (string, error) {
	b := make([]byte, csrfTokenLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
