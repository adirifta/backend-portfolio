// Package auth provides JWT token generation, validation, and cookie management.
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents the JWT payload for both access and refresh tokens.
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTService handles token generation and validation using separate secrets
// for access and refresh tokens (defense in depth).
type JWTService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTService creates a JWTService with the given secrets and expiry durations.
func NewJWTService(accessSecret, refreshSecret string, accessExpiryMin, refreshExpiryMin int) *JWTService {
	return &JWTService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessExpiry:  time.Duration(accessExpiryMin) * time.Minute,
		refreshExpiry: time.Duration(refreshExpiryMin) * time.Minute,
	}
}

// GenerateAccessToken creates a short-lived access token.
func (s *JWTService) GenerateAccessToken(userID uint, role string) (string, error) {
	return s.generateToken(userID, role, s.accessSecret, s.accessExpiry)
}

// GenerateRefreshToken creates a long-lived refresh token.
func (s *JWTService) GenerateRefreshToken(userID uint, role string) (string, error) {
	return s.generateToken(userID, role, s.refreshSecret, s.refreshExpiry)
}

// ValidateAccessToken validates and parses an access token string.
func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString, s.accessSecret)
}

// ValidateRefreshToken validates and parses a refresh token string.
func (s *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString, s.refreshSecret)
}

// AccessSecret returns the access token HMAC key (used by CSRF for signing).
func (s *JWTService) AccessSecret() []byte {
	return s.accessSecret
}

func (s *JWTService) generateToken(userID uint, role string, secret []byte, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "backend-portfolio",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (s *JWTService) validateToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
