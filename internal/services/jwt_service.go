package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/config"
)

type JWTService struct {
	cfg config.Config
}

func NewJWTService(cfg config.Config) *JWTService {
	return &JWTService{cfg: cfg}
}

func (s *JWTService) NewAccessToken(userID uint) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": s.cfg.JWTIssuer,
		"sub": userID,
		"typ": "access",
		"iat": now.Unix(),
		"exp": now.Add(time.Minute * time.Duration(s.cfg.AccessMinutes)).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.cfg.AccessSecret))
}

func (s *JWTService) ParseAccess(token string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.AccessSecret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	if claims["typ"] != "access" {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}
