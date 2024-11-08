package utils

import (
	"errors"
	"time"

	"github.com/adeyemialameen04/unwind-be/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token has expired")
	ErrInvalidTokenType = errors.New("invalid token type")
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
		"id":      userId,
		"expires": expiration.Unix(),
		"refresh": refresh,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateTokens(tokenStr string, cfg *config.Config) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			isRefresh, ok := claims["refresh"].(bool)
			if !ok {
				return nil, ErrInvalidToken
			}

			if isRefresh {
				return []byte(cfg.RefreshJwtKey), nil
			}

			return []byte(cfg.AccessJwtKey), nil
		}

		return nil, ErrInvalidToken
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if exp, ok := claims["expires"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, ErrTokenExpired
		}
	} else {
		return nil, ErrInvalidToken
	}

	if _, ok := claims["id"].(string); !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
