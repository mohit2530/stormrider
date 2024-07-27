package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/mohit2530/jwt-tokens/types"
)

func Login(w http.ResponseWriter, r *http.Request) {

	creds, err := verifyUser(r)
	if err != nil {
		log.Printf("unable to decode the request body. err - %+v", err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	types.CreateJwt(w, r, creds)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {

	creds, err := verifyUser(r)
	if err != nil {
		log.Printf("unable to decode the request body. err - %+v", err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = refreshToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	types.CreateJwt(w, r, creds)
}

func verifyUser(r *http.Request) (types.Credentials, error) {

	var credentials types.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Printf("unable to decode the request body. err - %+v", err)
		return types.Credentials{}, err
	}

	var fakeUsers = map[string]string{
		"admin": "password",
	}

	// simulates a test env. NOT OK TO USE IN PROD ENV
	expectedPwd, ok := fakeUsers[credentials.Username]
	if !ok || expectedPwd != credentials.Password {
		err := errors.New("unable to verify password ")
		return types.Credentials{}, err
	}

	return credentials, nil
}

func refreshToken(r *http.Request) error {

	claims := &types.Claims{}
	// control jwt token; prevent refresh if > 30 seconds of expiry
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return errors.New(" subject within jwt time limit")
	}
	expTime := time.Now().Add(7 * time.Minute)
	claims.ExpiresAt = expTime.Unix()
	return nil
}
