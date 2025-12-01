package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims custom JWT claims
type JWTClaims struct {
	UserId    int    `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserLevel string `json:"user_level"`
	jwt.RegisteredClaims
}

// JWTService handles JWT operations
type JWTService struct {
	secretKey    string
	expireHours  int
}

// NewJWTService creates new JWT service
func NewJWTService(secretKey string, expireHours int) *JWTService {
	return &JWTService{
		secretKey:   secretKey,
		expireHours: expireHours,
	}
}

// GenerateToken generates JWT token
func (s *JWTService) GenerateToken(userId int, email, level string) (string, error) {
	claims := JWTClaims{
		UserId:    userId,
		UserEmail: email,
		UserLevel: level,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.expireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken validates JWT token
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
