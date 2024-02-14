package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

func getPublicKey(url string) (jwk.Set, error) {
	set, err := jwk.Fetch(context.Background(), url)
	if err != nil {
		log.Printf("failed to parse JWK: %s", err)
		return nil, err
	}
	return set, nil
}
func verifyJwtMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Authenticating the user")
		log.Println("Parsing Authorization Header next")

		header := r.Header.Get("Authorization")
		if header == "" {
			log.Println("Auth failed due to missing Authorization Header")
			WriteJson(w, http.StatusUnauthorized, ApiLog{Err: "Missing Authorization Header", StatusCode: http.StatusUnauthorized})
			return
		}
		if !strings.HasPrefix(header, "Bearer ") {
			log.Println("Auth failed due to Invalid Authorization Header")
			WriteJson(w, http.StatusUnauthorized, ApiLog{Err: "Invalid Authorization Header", StatusCode: http.StatusUnauthorized})
			return
		}

		tokenString := strings.Split(header, "Bearer ")[1]
		log.Println("Authorization Header parsed")

		tokenByte := []byte(tokenString)

		log.Println("Fetching the public key")
		verifyKeySet, err := getPublicKey(cognito_jwk_url)
		if err != nil {
			log.Println("Error fetching the public key: ", err)
			WriteJson(w, http.StatusInternalServerError, ApiLog{Err: "Error fetching the public key", StatusCode: http.StatusInternalServerError})
			return
		}

		log.Println("Parsing the token with the public key to validate")
		token, err := jwt.Parse(tokenByte, jwt.WithKeySet(verifyKeySet), jwt.WithValidate(true))
		if err != nil {
			log.Println("Error parsing the token with err message: ", err)
			WriteJson(w, http.StatusUnauthorized, ApiLog{Err: fmt.Sprintf("Error parsing the token with err message: %s", err.Error()), StatusCode: http.StatusUnauthorized})
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
			WriteJson(
				w,
				http.StatusInternalServerError,
				ApiLog{Err: fmt.Sprintf("Error validating user in the database with err msg: %s",
					err.Error()),
					StatusCode: http.StatusInternalServerError,
				})
			return
		}
		log.Println("Setting header with cognito Id")
		r.Header.Set("CognitoId", cognitoId)
		log.Println("Serving next handler")
		next.ServeHTTP(w, r)
	})
}

func validadeUserDb(cognitoId string) error {
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
