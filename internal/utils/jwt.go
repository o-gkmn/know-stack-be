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
	ErrTokenExpired     = errors.New("token expired")
)

// TokenClaims defines the JWT claims used across the application.
// It embeds jwt.RegisteredClaims and adds a domain-specific UserID.
type TokenClaims struct {
	UserID   string   `json:"uid"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	RoleID   uint     `json:"role_id"`
	Claims   []string `json:"claim_ids"`
	jwt.RegisteredClaims
}

type RefreshTokenClaim struct {
	UserID  string `json:"uid"`
	TokenID string `json:"tokenID"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a signed JWT token for the provided userID.
// It uses HMAC-SHA256 and reads configuration from environment variables:
// - JWT_SECRET: signing key (default: "dev_secret")
// - JWT_EXPIRES_IN_MIN: expiration in minutes (default: 60)
func GenerateAccessToken(userID string, email string, username string, roleID uint, claimNames []string) (string, error) {
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
		RoleID:   roleID,
		Claims:   claimNames,
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

// VerifyAccessToken validates the token signature and expiration and returns parsed claims.
// It reads JWT_SECRET from the environment (default: "dev_secret").
func VerifyAccessToken(tokenString string) (*TokenClaims, error) {
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

// Generates a refresh token with longer expire times
func GenerateRefreshToken(userID, tokenID string, remember bool) (string, error) {
	secret := GetEnv("JWT_REFRESH_SECRET", "dev_refresh_secret")
	audiance := GetEnv("JWT_AUDIENCE", "knowstack")
	issuer := GetEnv("JWT_ISSUER", "knowstack")

	now := time.Now()
	refreshExpiresInDays := GetEnvAsInt("JWT_REFRESH_EXPIRES_IN_DAYS", 7)
	if remember {
		refreshExpiresInDays = GetEnvAsInt("JWT_REFRESH_EXPIRES_IN_DAYS_REMEMBER", 30)
	}

	expireAt := now.Add(time.Duration(refreshExpiresInDays) * time.Hour * 24)

	claims := RefreshTokenClaim{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{audiance},
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateRefreshToken validates the refresh token signature and expirations and returns parsed claims
func ValidateRefreshToken(token string) (*RefreshTokenClaim, error) {
	secret := GetEnv("JWT_SECRET", "dev_seecret")
	issuer := GetEnv("JWT_ISSUER", "knowstack")
	audience := GetEnv("JWT_AUDIENCE", "knowstack")

	parsedToken, err := jwt.ParseWithClaims(token, &RefreshTokenClaim{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignature
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*RefreshTokenClaim)
	if !ok || !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Issuer != issuer {
		return nil, ErrInvalidIssuer
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

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrTokenExpired
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
