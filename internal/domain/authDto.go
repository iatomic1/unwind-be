package domain

import (
	"github.com/adeyemialameen04/unwind-be/internal/db/repository"
	"github.com/adeyemialameen04/unwind-be/internal/utils"
	"github.com/google/uuid"
)

const (
	UniqueViolation           = "23505"
	UserCreated               = "User Created Successfully"
	ErrEmailAlreadyExist      = "User with this email already exists"
	ErrUsernameAlreadyExist   = "User with this username already exists"
	ErrInvalidEmailOrPassword = "Invalid email or password"
	LoginSuccessful           = "Login Successful"
	ErrGeneratingTokens       = "Error generating token pair"
	TokensRefreshed           = "Tokens Refreshed successfully"
)

type AuthResponse struct {
	utils.TokenPair
	User EmailID `json:"user"`
}

type RegisterRequest struct {
	repository.RegisterUserParams
	Username string `json:"username" binding:"required"`
}

type RefreshTokenResponse struct {
	utils.TokenPair
}

type EmailID struct {
	Email     string `json:"email"`
	ID        string `json:"id"`
	ProfileID string `json:"profileId"`
} // @name EmailID

func ParseIDs(id string, profileId string) (uuid.UUID, uuid.UUID, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	parsedProfileId, err := uuid.Parse(profileId)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return userId, parsedProfileId, nil
}
