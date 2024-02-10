package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadENV() {
	// Load variables from .env file
	if os.Getenv("APP_ENV") != "production" {

		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
		connStr, ok = os.LookupEnv("LOCAL_DB")
		if !ok {
			log.Fatal("LOCAL_DB not found in .env file")
		}
		hostPort, ok = os.LookupEnv("PORT")
		if !ok {
			hostPort = "8000"
		}

		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
		cognito_jwk_url, ok = os.LookupEnv("COGNITO_JWK_URL")
		if !ok {
			log.Fatalln("Cognito Url not found in .env file")
		}
		// get the cognito endpoint
		cognito_issuer, ok = os.LookupEnv("COGNITO_ISSUER")
		if !ok {
			log.Fatalln("Cognito Issuer not found in .env file")
		}

		listenAddr = "localhost:" + hostPort

	} else {
		connStr = os.Getenv("DATABASE_URL")
		hostPort = os.Getenv("PORT")
		listenAddr = `0.0.0.0:` + hostPort
		cognito_jwk_url = os.Getenv("COGNITO_JWK_URL")
		cognito_issuer = os.Getenv("COGNITO_ISSUER")
	}
}
