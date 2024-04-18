package auth

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	MaxAge   = 86400 * 30
	isProd   = false
	Path     = "/"
	HttpOnly = true
)

func New() {
	err := godotenv.Load("D:/code/Go/PetUrKitten/Gateway/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	key := os.Getenv("API_KEY")
	GoogleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	store := sessions.NewCookieStore([]byte(key))
	store.Options = &sessions.Options{
		Path:     Path,
		MaxAge:   MaxAge,
		HttpOnly: HttpOnly,
		Secure:   isProd,
	}

	gothic.Store = store
	goth.UseProviders(
		google.New(GoogleClientID, GoogleClientSecret, "http://localhost:8080/auth/callback/google"),
	)
}
