package utils

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("expired token")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrInvalidIssuer    = errors.New("invalid issuer")
	ErrInvalidAudience  = errors.New("invalid audience")
	ErrInvalidSubject   = errors.New("invalid subject")
	ErrInvalidIssuedAt  = errors.New("invalid issued at")
	ErrInvalidExpiresAt = errors.New("invalid expires at")
)

// TokenClaims defines the JWT claims used across the application.
// It embeds jwt.RegisteredClaims and adds a domain-specific UserID.
type TokenClaims struct {
	UserID   string `json:"uid"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed JWT token for the provided userID.
// It uses HMAC-SHA256 and reads configuration from environment variables:
// - JWT_SECRET: signing key (default: "dev_secret")
// - JWT_EXPIRES_IN_MIN: expiration in minutes (default: 60)
func GenerateJWT(userID string, email string, username string) (string, error) {
	secret := GetEnv("JWT_SECRET", "dev_secret")
	issuer := GetEnv("JWT_ISSUER", "knowstack")
	audience := GetEnv("JWT_AUDIENCE", "knowstack")

	now := time.Now()
	expiresInMinutes := GetEnvAsInt("JWT_EXPIRES_IN_MIN", 60)
	expiresAt := now.Add(time.Duration(expiresInMinutes) * time.Minute)

	claims := TokenClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{audience},
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// VerifyJWT validates the token signature and expiration and returns parsed claims.
// It reads JWT_SECRET from the environment (default: "dev_secret").
func VerifyJWT(tokenString string) (*TokenClaims, error) {
	secret := GetEnv("JWT_SECRET", "dev_secret")
	issuer := GetEnv("JWT_ISSUER", "knowstack")
	audience := GetEnv("JWT_AUDIENCE", "knowstack")
	parsedToken, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok || !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Issuer != issuer {
		return nil, ErrInvalidIssuer
	}

	if claims.Subject != claims.UserID {
		return nil, ErrInvalidSubject
	}

	aud := claims.Audience
	if len(aud) == 0 {
		return nil, ErrInvalidAudience
	}
	matched := false
	for _, a := range aud {
		if a == audience {
			matched = true
			break
		}
	}
	if !matched {
		return nil, ErrInvalidAudience
	}

	return claims, nil
}

// ExtractBearerToken returns the token part from a typical Authorization header value.
// If the header is not in the expected format, it returns an empty string.
func ExtractBearerToken(authorizationHeader string) string {
	const prefix = "Bearer "
	if len(authorizationHeader) > len(prefix) && authorizationHeader[:len(prefix)] == prefix {
		return authorizationHeader[len(prefix):]
	}
	return ""
}
