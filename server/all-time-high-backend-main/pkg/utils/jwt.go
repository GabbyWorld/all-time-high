// internal/utils/jwt.go
package utils

import (
	"time"

	"github.com/GabbyWorld/all-time-high-backend/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

type JWTManager struct {
	secret      string
	tokenExpiry time.Duration
}

func NewJWTManager(secret string, expiry string) (*JWTManager, error) {
	duration, err := time.ParseDuration(expiry)
	if err != nil {
		return nil, err
	}
	return &JWTManager{
		secret:      secret,
		tokenExpiry: duration,
	}, nil
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken生成JWT
func (j *JWTManager) GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-web-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// VerifyToken验证并解析JWT
func (j *JWTManager) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}
