package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mohit2530/jwt-tokens/api"
)

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
	router.HandleFunc("/validate", api.Login).Methods(http.MethodPost)
	router.HandleFunc("/refreshToken", api.RefreshToken).Methods(http.MethodPost)

	log.Printf("token-generator is warming up.. 🌞🌞🌞")
	http.ListenAndServe(formattedEndpoint, router)
}
