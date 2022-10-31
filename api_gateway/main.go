package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var LICENSE = []byte("2530d6a4-5d42-4758-b331-2fbbfed27bf9")

type UserRequestHandler func(http.ResponseWriter, *http.Request)

func main() {

	router := mux.NewRouter()

	router.Handle("/secure/storage", UserRequestHandler(SecureLandingPage)).Methods(http.MethodPost)
	router.Handle("/public/storage", UserRequestHandler(LandingPage)).Methods(http.MethodPost)

	log.Printf(" gateway is warming up ")
	http.ListenAndServe(":3009", router)
}

func SecureLandingPage(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode("Hello World !! Welcome to Secure")
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	return
}

func LandingPage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello world !! Welcome to Public ")
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	return
}

func (u UserRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// v2 ugrade for user validation

	cookie, err := r.Cookie("token")

	if err != nil {
		log.Printf(" missing license key ")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenStr := cookie.Value
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
