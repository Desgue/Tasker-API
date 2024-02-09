package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

func getPublicKey(url string) jwk.Set {
	set, err := jwk.Fetch(context.Background(), url)
	if err != nil {
		log.Printf("failed to parse JWK: %s", err)

	}
	return set
}
func verifyJwtMiddleware(next http.Handler) http.Handler {
	// get the public key from the cognito endpoint
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	cognito_jwk_url, ok := os.LookupEnv("COGNITO_JWK_URL")
	if !ok {
		log.Fatalln("Cognito Url not found in .env file")
	}
	// get the cognito endpoint
	cognito_issuer, ok := os.LookupEnv("COGNITO_ISSUER")
	if !ok {
		log.Fatalln("Cognito Issuer not found in .env file")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Authenticating the user")
		tokenString := r.Header.Get("Authorization")
		tokenByte := []byte(tokenString)
		verifyKeySet := getPublicKey(cognito_jwk_url)

		token, err := jwt.Parse(tokenByte, jwt.WithKeySet(verifyKeySet), jwt.WithValidate(true))
		if err != nil {
			log.Println("Error parsing the token: ", err)
			WriteJson(w, http.StatusUnauthorized, ApiLog{Err: "Error parsing the token", StatusCode: http.StatusUnauthorized})
			return
		}
		// Compare token claims for: Expire time, issuer, token_use
		if token.Expiration().Unix() < time.Now().Unix() {
			log.Println("Auth failed due to Token expired")
			WriteJson(w, http.StatusUnauthorized, ApiLog{Err: "Token expired", StatusCode: http.StatusUnauthorized})
			return
		}
		if token.Issuer() != cognito_issuer {
			log.Println("Auth failied due to Invalid issuer")
			WriteJson(w, http.StatusUnauthorized, ApiLog{Err: "Invalid issuer", StatusCode: http.StatusUnauthorized})
			return
		}
		if token.PrivateClaims()["token_use"] != "access" {
			log.Println("Auth failied due to Invalid token use")
			WriteJson(w, http.StatusUnauthorized, ApiLog{Err: "Invalid token use", StatusCode: http.StatusUnauthorized})
			return
		}
		log.Println("User authenticated successfully")

		// Check if user is present on the database, if not create a new user
		cognitoId := token.Subject()
		if err := validadeUserDb(cognitoId); err != nil {
			log.Println("Error validating user in the database: ", err)
			WriteJson(w, http.StatusInternalServerError, ApiLog{Err: "Error validating user in the database", StatusCode: http.StatusInternalServerError})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validadeUserDb(cognitoId string) error {
	connStr, ok := os.LookupEnv("DB_CONNSTR")
	if !ok {
		log.Fatalln("DB_CONNSTR not found in .env file")
	}

	db, err := NewPostgresStore(connStr)
	if err != nil {
		log.Println("Error initializing database", err)
		return err
	}
	defer db.db.Close()

	userStore := NewPostgresUserStore(db.db)

	ok, _ = userStore.CheckUser(cognitoId)
	if !ok {
		log.Println("User not found in the database, creating a new user")
		userStore.CreateUser(cognitoId)
	}
	return nil
}
