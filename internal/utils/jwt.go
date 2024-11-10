package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/adeyemialameen04/unwind-be/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type EmailID struct {
	Email     string `json:"email"`
	ID        string `json:"id"`
	ProfileId string `json:"profileId"`
}

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token has expired")
	ErrInvalidTokenType = errors.New("invalid token type")
)

func CreateJWT(data EmailID, refresh bool, cfg *config.Config) (string, error) {
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
		"data": EmailID{
			ID:        data.ID,
			Email:     data.Email,
			ProfileId: data.ProfileId,
		},
		"jti":     uuid.New(),
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
			fmt.Println("Debug 1: Invalid signing method")
			return nil, ErrInvalidToken
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			isRefresh, ok := claims["refresh"].(bool)
			if !ok {
				fmt.Println("Debug 2: refresh claim not found or not boolean")
				return nil, ErrInvalidToken
			}
			if isRefresh {
				return []byte(cfg.RefreshJwtKey), nil
			}
			return []byte(cfg.AccessJwtKey), nil
		}
		fmt.Println("Debug 3: Claims casting failed")
		return nil, ErrInvalidToken
	})
	if err != nil {
		fmt.Println("Debug 4: JWT parsing failed with error:", err)
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		fmt.Println("Debug 5: Token claims invalid or token not valid")
		return nil, ErrInvalidToken
	}

	if exp, ok := claims["expires"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			fmt.Println("Debug 6: Token expired")
			return nil, ErrTokenExpired
		}
	} else {
		fmt.Println("Debug 7: expires claim not found or not float64")
		return nil, ErrInvalidToken
	}

	if _, ok := claims["id"].(string); !ok {
		fmt.Println("Debug 8: id claim not found or not string")
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func ValidateAccessToken(tokenStr string, cfg *config.Config) (jwt.MapClaims, error) {
	claims, err := ValidateTokens(tokenStr, cfg)
	if err != nil {
		fmt.Println("Debug 9: ValidateTokens failed:", err)
		return nil, err
	}

	if refresh, ok := claims["refresh"].(bool); !ok || refresh {
		fmt.Println("Debug 10: Invalid token type (refresh token used instead of access token)")
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

func ValidateRefreshToken(tokenStr string, cfg *config.Config) (jwt.MapClaims, error) {
	claims, err := ValidateTokens(tokenStr, cfg)
	if err != nil {
		return nil, err
	}

	if refresh, ok := claims["refresh"].(bool); !ok || !refresh {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

func ExtractID(claims jwt.MapClaims) (string, string, error) {
	data, ok := claims["data"].(map[string]interface{})
	if !ok {
		return "", "", ErrInvalidToken
	}

	userId, ok := data["id"].(string)
	if !ok {
		return "", "", ErrInvalidToken
	}
	profileId, ok := data["profileId"].(string)
	if !ok {
		return "", "", ErrInvalidToken
	}

	return userId, profileId, nil
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func GenerateTokenPair(data EmailID, cfg *config.Config) (TokenPair, error) {
	accessToken, err := CreateJWT(data, false, cfg)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := CreateJWT(data, true, cfg)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
