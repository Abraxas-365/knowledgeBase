package conf

import (
	"os"
	"strings"
)

type Conf struct {
	GoogleConf
	CorsConf
	RedirectAfterLogin string
	DatabaseURL        string
	Port               string
}

type GoogleConf struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
}

type CorsConf struct {
	AllowOrigins string
}

func Load() Conf {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	port = ":" + port

	uri := os.Getenv("DATABASE_URL")
	if uri == "" {
		panic("DATABASE_URL is not set")
	}
	if !strings.Contains(uri, "sslmode") {
		uri += "?sslmode=disable"
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		panic("GOOGLE_CLIENT_ID is not set")
	}
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientSecret == "" {
		panic("GOOGLE_CLIENT_SECRET is not set")
	}
	googleRedirectURI := os.Getenv("GOOGLE_REDIRECT_URI")
	if googleRedirectURI == "" {
		panic("GOOGLE_REDIRECT_URI is not set")
	}

	redirectAfterLogin := os.Getenv("REDIRECT_AFTER_LOGIN")
	if redirectAfterLogin == "" {
		panic("REDIRECT_AFTER_LOGIN is not set")
	}
	allowOrigins := os.Getenv("ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:3001, http://localhost:3000"
	}

	return Conf{
		GoogleConf: GoogleConf{
			GoogleClientID:     googleClientID,
			GoogleClientSecret: googleClientSecret,
			GoogleRedirectURI:  googleRedirectURI,
		},
		CorsConf: CorsConf{
			AllowOrigins: allowOrigins,
		},
		RedirectAfterLogin: redirectAfterLogin,
		DatabaseURL:        uri,
		Port:               port,
	}

}
