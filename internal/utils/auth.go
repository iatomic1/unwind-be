package utils

import (
	"time"

	"github.com/adeyemialameen04/unwind-be/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateJWT(userId string, refresh bool, cfg *config.Config) (string, error) {
	var (
		expiration time.Time
		secret     []byte
	)

	if refresh {
		expiration = time.Now().Add(time.Duration(cfg.RefreshExpirationHour) * time.Hour)
		secret = []byte(cfg.RefreshJwtKey)
	} else {
		expiration = time.Now().Add(time.Duration(cfg.AccessExpirationHour) * time.Hour)
		secret = []byte(cfg.AccessJwtKey)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  userId,
		"expires": expiration.Unix(),
		"refresh": refresh,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
