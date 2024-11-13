package auth

import (
	"github.com/adeyemialameen04/unwind-be/internal/config"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func NewGoogleAuth(cfg *config.Config) {
	store := sessions.NewCookieStore([]byte(cfg.GoogleSigningKey))
	store.MaxAge(cfg.GoogleMaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	if cfg.Environment == "development" {
		store.Options.Secure = false
	} else {
		store.Options.Secure = true
	}

	gothic.Store = store
	goth.UseProviders(
		google.New(cfg.GoogleClientID, cfg.GoogleClientSecret, "http://localhost:8080/api/v1/auth/google/callback", "email", "profile"),
	)
}
