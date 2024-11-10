package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/config"
	"github.com/adeyemialameen04/unwind-be/internal/utils"
	"github.com/gin-gonic/gin"
)

func AccessTokenMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			server.SendUnauthorized(c, nil, server.WithMessage("Missing authentication token"))
			c.Abort()
			return
		}

		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			server.SendUnauthorized(c, errors.New("Invalid authentication token"), server.WithMessage("Invalid authentication token"))
			c.Abort()
			return
		}

		tokenString = tokenParts[1]
		claims, err := utils.ValidateAccessToken(tokenString, cfg)
		fmt.Println(claims)
		if err != nil {
			switch err {
			case utils.ErrTokenExpired:
				server.SendUnauthorized(c, err, server.WithMessage("Token has expired"))
			case utils.ErrInvalidTokenType:
				server.SendUnauthorized(c, err, server.WithMessage("Invalid token type"))
			default:
				server.SendUnauthorized(c, err, server.WithMessage("Invalid"))
			}
			c.Abort()
			return
		}

		data, err := utils.ExtractDataFromToken(claims)
		if err != nil {
			server.SendUnauthorized(c, err)
			c.Abort()
			return
		}
		c.Set("userId", data.ID)
		c.Set("profileId", data.ProfileId)
		c.Set("email", data.Email)
	}
}

func RefreshToenMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			server.SendUnauthorized(c, nil, server.WithMessage("Missing authentication token"))
			c.Abort()
			return
		}

		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			server.SendUnauthorized(c, errors.New("Invalid authentication token"), server.WithMessage("Invalid authentication token"))
			c.Abort()
			return
		}

		tokenString = tokenParts[1]
		claims, err := utils.ValidateRefreshToken(tokenString, cfg)
		if err != nil {
			switch err {
			case utils.ErrTokenExpired:
				server.SendUnauthorized(c, err, server.WithMessage("Refresh token has expired"))
			case utils.ErrInvalidTokenType:
				server.SendUnauthorized(c, err, server.WithMessage("Invalid refresh token"))
			default:
				server.SendUnauthorized(c, err, server.WithMessage("Invalid token"))
			}
			c.Abort()
			return
		}
		data, err := utils.ExtractDataFromToken(claims)
		if err != nil {
			server.SendUnauthorized(c, err)
			c.Abort()
			return
		}
		c.Set("userId", data.ID)
		c.Set("profileId", data.ProfileId)
		c.Set("email", data.Email)
	}
}
