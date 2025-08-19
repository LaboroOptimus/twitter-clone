package authjwt

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Sub   string `json:"sub"`
	Jti   string `json:"jti"`
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

type TokenService interface {
	GenerateAccessToken(ctx context.Context, userID uint, nickname string) (token string, exp time.Time, err error)
	GenerateRefreshToken(ctx context.Context, userID uint, nickname string) (token string, exp time.Time, jti string, err error)

	ParseAndValidateAccess(tokenStr string) (*Claims, error)
	ParseAndValidateRefresh(tokenStr string) (*Claims, error)
}

type jwtService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) TokenService {
	return &jwtService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}
func (s *jwtService) GenerateAccessToken(ctx context.Context, userID uint, nickname string) (string, time.Time, error) {
	exp := time.Now().Add(s.accessTTL)

	claims := Claims{
		Sub:   strconv.FormatUint(uint64(userID), 10),
		Jti:   uuid.NewString(),
		Scope: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   nickname,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.accessSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, exp, nil
}

func (s *jwtService) GenerateRefreshToken(ctx context.Context, userID uint, nickname string) (string, time.Time, string, error) {
	exp := time.Now().Add(s.refreshTTL)
	jti := uuid.NewString()

	claims := Claims{
		Sub:   strconv.FormatUint(uint64(userID), 10),
		Jti:   jti,
		Scope: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   nickname,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.refreshSecret)
	if err != nil {
		return "", time.Time{}, "", err
	}

	return signed, exp, jti, nil
}

func (s *jwtService) ParseAndValidateAccess(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.accessSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.Scope != "access" {
		return nil, fmt.Errorf("invalid token scope")
	}

	return claims, nil

}
func (s *jwtService) ParseAndValidateRefresh(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		// Проверка алгоритма
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.refreshSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.Scope != "refresh" {
		return nil, fmt.Errorf("invalid token scope")
	}

	return claims, nil
}
