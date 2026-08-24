package utils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TokenTypeBearer = "Bearer"

type JWTClaims struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	SecretKey string
	ExpiredIn time.Duration
}

func NewJWTManager(secret string, expiredIn time.Duration) *JWTManager {
	return &JWTManager{
		SecretKey: secret,
		ExpiredIn: expiredIn,
	}
}

func (j *JWTManager) GenerateToken(userID uint64, email, role string) (string, error) {
	now := time.Now()

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ExpiredIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTManager) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
