package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CookieNames used for authentication.
const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

// CookieManager handles setting and clearing HTTP cookies for auth tokens.
type CookieManager struct {
	domain        string
	secure        bool
	sameSite      http.SameSite
	accessMaxAge  int // seconds
	refreshMaxAge int // seconds
}

// NewCookieManager creates a CookieManager from application config values.
func NewCookieManager(domain string, secure bool, sameSiteStr string, accessExpiryMin, refreshExpiryMin int) *CookieManager {
	return &CookieManager{
		domain:        domain,
		secure:        secure,
		sameSite:      parseSameSite(sameSiteStr),
		accessMaxAge:  accessExpiryMin * 60,
		refreshMaxAge: refreshExpiryMin * 60,
	}
}

// SetTokenCookies writes HttpOnly access and refresh token cookies.
func (cm *CookieManager) SetTokenCookies(c *gin.Context, accessToken, refreshToken string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    accessToken,
		Path:     "/",
		Domain:   cm.domain,
		MaxAge:   cm.accessMaxAge,
		Secure:   cm.secure,
		HttpOnly: true,
		SameSite: cm.sameSite,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/api/auth/refresh",
		Domain:   cm.domain,
		MaxAge:   cm.refreshMaxAge,
		Secure:   cm.secure,
		HttpOnly: true,
		SameSite: cm.sameSite,
	})
}

// ClearTokenCookies removes both authentication cookies.
func (cm *CookieManager) ClearTokenCookies(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    "",
		Path:     "/",
		Domain:   cm.domain,
		MaxAge:   -1,
		Secure:   cm.secure,
		HttpOnly: true,
		SameSite: cm.sameSite,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Path:     "/api/auth/refresh",
		Domain:   cm.domain,
		MaxAge:   -1,
		Secure:   cm.secure,
		HttpOnly: true,
		SameSite: cm.sameSite,
	})
}

// Domain returns the configured cookie domain.
func (cm *CookieManager) Domain() string { return cm.domain }

// Secure returns whether the Secure flag is set.
func (cm *CookieManager) Secure() bool { return cm.secure }

// SameSite returns the configured SameSite mode.
func (cm *CookieManager) SameSite() http.SameSite { return cm.sameSite }

// AccessMaxAge returns the access token cookie max age in seconds.
func (cm *CookieManager) AccessMaxAge() int { return cm.accessMaxAge }

func parseSameSite(s string) http.SameSite {
	switch strings.ToLower(s) {
	case "lax":
		return http.SameSiteLaxMode
	case "none":
		return http.SameSiteNoneMode
	case "strict":
		return http.SameSiteStrictMode
	default:
		return http.SameSiteStrictMode
	}
}
