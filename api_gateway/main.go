package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type UserRequestHandler func(http.ResponseWriter, *http.Request)

func main() {
	err := godotenv.Load(filepath.Join(".env"))
	if err != nil {
		log.Printf("No env file detected. Please follow steps in readme.")
		return
	}

	tokenEndpoint := os.Getenv("TOKEN_ENDPOINT")
	if tokenEndpoint == "" {
		log.Println("TOKEN_ENDPOINT is not found. Exiting.")
		return
	}

	formattedEndpoint := fmt.Sprintf(":%s", tokenEndpoint)

	router := mux.NewRouter()
	router.Handle("/secure/storage", UserRequestHandler(SecureLandingPage)).Methods(http.MethodPost)
	router.Handle("/public/storage", UserRequestHandler(LandingPage)).Methods(http.MethodPost)

	log.Printf("api gateway is warming up ... 🌞🌞🌞")
	http.ListenAndServe(formattedEndpoint, router)
}

func SecureLandingPage(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode("Hello World !! Welcome to Secure")
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
}

func LandingPage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello world !! Welcome to Public ")
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
}

func (u UserRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")
	if err != nil {
		log.Printf(" missing license key")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = godotenv.Load(filepath.Join("..", "LICENSE"))
	if err != nil {
		log.Printf("No license key file found. Are you authorized?")
		return
	}

	validKey := os.Getenv("LICENSE")
	if validKey == "" {
		log.Println("No license key file found. Exiting.")
		return
	}

	var LICENSE = []byte(validKey)

	tokenStr := cookie.Value
	fmt.Printf("token is - %+v", cookie.Value)
	token, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid license detected ")
		}
		// v2 upgrade to include auditor && issuer check
		return LICENSE, nil
	})

	if token.Valid {
		u(w, r)
	} else {
		// no token || invalid token
		http.Error(w, "invalid license ", http.StatusUnauthorized)
		return
	}
}
