package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const AuthCookieName = "auth_token"

type Provider struct {
	salt                string
	jwtSecret           []byte
	JwtTTL              time.Duration
	JwtRenewalThreshold time.Duration
}

func NewProvider(salt string,
	jwtSecret []byte,
	jwtTTL time.Duration,
	jwtRenewalThreshold time.Duration,
) *Provider {
	return &Provider{
		salt:                salt,
		jwtSecret:           jwtSecret,
		JwtTTL:              jwtTTL,
		JwtRenewalThreshold: jwtRenewalThreshold,
	}
}

// HashPassword hashes a password using the salt and returns the hash.
func (s *Provider) HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password + s.salt))
	return hex.EncodeToString(hash[:])
}

// VerifyPassword verifies a plaintext password against a salted password hash.
func (s *Provider) VerifyPassword(password string, hash string) bool {
	computedHash := s.HashPassword(password)
	return computedHash == hash
}

// GenerateToken generates a JWT token for a user.
func (s *Provider) GenerateToken(userUUID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userUUID.String(),
		"exp": time.Now().Add(s.JwtTTL).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString(s.jwtSecret)
}

// VerifyToken verifies a JWT token and returns the user UUID and whether the token needs to be renewed.
func (s *Provider) VerifyToken(token string) (uuid.UUID, bool, error) {
	claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return uuid.UUID{}, false, err
	}

	renewalRequired, err := s.renewalRequired(claims.Claims)
	if err != nil {
		return uuid.UUID{}, false, err
	}

	sub, err := claims.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, false, err
	}

	uUUID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.UUID{}, false, err
	}

	return uUUID, renewalRequired, nil
}

// SetTokenCookie sets/unsets the JWT token in a cookie Gin context.
// If the token is empty, the cookie is unset.
func (s *Provider) SetTokenCookie(c *gin.Context, token string) {
	if token == "" {
		c.SetCookie(AuthCookieName, "", -1, "/", "", true, true)
		return
	}

	c.SetCookie(AuthCookieName, token, int(s.JwtTTL.Seconds()), "/", "", true, true)
}

// renewalRequired checks if the JWT token needs to be renewed.
func (s *Provider) renewalRequired(claims jwt.Claims) (bool, error) {
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return false, err
	}

	if exp == nil {
		return false, jwt.ErrInvalidKey
	}

	return time.Until(exp.Time) > s.JwtRenewalThreshold, nil
}
