package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jarellano2x5/finance-go/config"
)

func GenerateJWT(id uuid.UUID, nm string) (string, error) {
	cfg := config.LoadConfig()
	var sec = []byte(cfg.JwtScrt)

	claims := jwt.MapClaims{
		"user_id": id.String(),
		"name":    nm,
		"exp":     time.Now().Add(time.Hour * 10).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sec)
}

func ParseJWT(tks string) (jwt.MapClaims, error) {
	cfg := config.LoadConfig()
	var sec = []byte(cfg.JwtScrt)

	token, err := jwt.Parse(tks, func(t *jwt.Token) (interface{}, error) {
		return sec, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	return claims, nil
}
