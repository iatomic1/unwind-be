package middleware

import (
	"errors"
	"strings"

	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/gin-gonic/gin"
)

func AccessTokenMiddleware() gin.HandlerFunc {
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
		// claims, err :=
	}
}
