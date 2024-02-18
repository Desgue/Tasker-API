package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	ConnStr         string
	HostPort        string
	ListenAddr      string
	Ok              bool
	Cognito_jwk_url string
	Cognito_issuer  string
	IsProd          bool
)

func LoadENV() {
	// Load variables from .env file
	if os.Getenv("APP_ENV") != "production" {
		log.Println("Loading .env file")
		IsProd = false

		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
		ConnStr, Ok = os.LookupEnv("LOCAL_DB")
		if !Ok {
			log.Fatal("LOCAL_DB not found in .env file")
		}
		HostPort, Ok = os.LookupEnv("PORT")
		if !Ok {
			HostPort = "8000"
		}

		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
		Cognito_jwk_url, Ok = os.LookupEnv("COGNITO_JWK_URL")
		if !Ok {
			log.Fatalln("Cognito Url not found in .env file")
		}
		// get the cognito endpoint
		Cognito_issuer, Ok = os.LookupEnv("COGNITO_ISSUER")
		if !Ok {
			log.Fatalln("Cognito Issuer not found in .env file")
		}

		ListenAddr = "localhost:" + HostPort

	} else {
		log.Println("Loading environment variables")
		IsProd = true
		ConnStr = os.Getenv("DATABASE_URL")
		HostPort = os.Getenv("PORT")
		ListenAddr = `0.0.0.0:` + HostPort
		Cognito_jwk_url = os.Getenv("COGNITO_JWK_URL")
		Cognito_issuer = os.Getenv("COGNITO_ISSUER")
	}
}
