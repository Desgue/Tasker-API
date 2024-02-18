package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	repo "github.com/Desgue/ttracker-api/internal/repository"
	"github.com/Desgue/ttracker-api/internal/util"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// JWT MIDDLEWARE

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
		verifyKeySet, err := getPublicKey(util.Cognito_jwk_url)
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
		if token.Issuer() != util.Cognito_issuer {
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

// HELPER FUNCTIONS FOR JWT MIDDLEWARE

func validadeUserDb(cognitoId string) error {
	db, err := repo.NewPostgresStore(util.ConnStr)
	if err != nil {
		log.Println("Error initializing database", err)
		return err
	}
	defer db.DB.Close()

	userStore := repo.NewPostgresUserStore(db.DB)

	ok, _ := userStore.CheckUser(cognitoId)
	if !ok {
		log.Println("User not found in the database, creating a new user")
		userStore.CreateUser(cognitoId)
	}
	return nil
}

func getPublicKey(url string) (jwk.Set, error) {
	set, err := jwk.Fetch(context.Background(), url)
	if err != nil {
		log.Printf("failed to parse JWK: %s", err)
		return nil, err
	}
	return set, nil
}

// LOGGING MIDDLEWARE

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("Logging the request")
		uri := r.RequestURI
		method := r.Method
		referrer := r.Referer()
		userAgent := r.UserAgent()

		log.Printf(` 
		%s -> %s 
		Referrer: %s 
		User-Agent: %s`,
			method, uri, referrer, userAgent)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		log.Println("Serving next handler")
		next.ServeHTTP(w, r)
	})
}
