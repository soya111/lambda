package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"zephyr/pkg/middleware"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/joho/godotenv"
)

var (
	Auth0Domain   string
	Auth0Audience string
	Auth0Secret   string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Auth0Domain = os.Getenv("AUTH0_DOMAIN")
	Auth0Audience = os.Getenv("AUTH0_AUDIENCE")
	Auth0Secret = os.Getenv("AUTH0_SECRET")
}

func main() {
	http.Handle("/api/private", middleware.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS Headers.
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Hello from a private endpoint! You need to be authenticated to see this."}`))

			// アクセスログを出力する
			log.Println("Access to private endpoint")
			v := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			fmt.Printf("%#v\n", v)
		}),
	))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
