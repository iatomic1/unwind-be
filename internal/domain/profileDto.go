package domain

import (
	"errors"
	"fmt"

	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetProfileIDFromContext(c *gin.Context) (uuid.UUID, error) {
	profileID, exists := c.Get("profileId")
	if !exists {
		server.SendUnauthorized(c, nil, server.WithMessage("Profile ID not found in context"))
		return uuid.Nil, errors.New("profile ID not found")
	}

	profileIDStr, ok := profileID.(string)
	if !ok {
		server.SendInternalServerError(c, nil, server.WithMessage("Invalid profile ID format"))
		return uuid.Nil, errors.New("invalid profile ID type")
	}

	parsedID, err := uuid.Parse(profileIDStr)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to parse profile ID"))
		return uuid.Nil, fmt.Errorf("parse UUID: %w", err)
	}

	return parsedID, nil
}

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userId")
	if !exists {
		server.SendUnauthorized(c, nil, server.WithMessage("User ID not found in context"))
		return uuid.Nil, errors.New("User ID not found")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		server.SendInternalServerError(c, nil, server.WithMessage("Invalid user ID format"))
		return uuid.Nil, errors.New("invalid user ID type")
	}

	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to parse user ID"))
		return uuid.Nil, fmt.Errorf("parse UUID: %w", err)
	}

	return parsedID, nil
}

func ParseIDFromParams(c *gin.Context, id string) (uuid.UUID, error) {
	if id == "" {
		server.SendBadRequest(c, errors.New("ID must be provided"))
		return uuid.Nil, errors.New("User ID not found")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to parse ID"))
		return uuid.Nil, fmt.Errorf("parse UUID: %w", err)
	}

	return parsedID, nil
}
