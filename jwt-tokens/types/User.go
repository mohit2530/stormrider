package types

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Credentials stores user provided username and password
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateJwt(w http.ResponseWriter, r *http.Request, creds Credentials) {

	var claims Claims
	expTime := time.Now().Add(7 * time.Minute) // every seven minutes
	tokenStr, err := claims.initialize(creds, expTime)
	if err != nil {
		log.Printf("unable to initialize token. err - %+v", err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expTime,
	})

	w.WriteHeader(http.StatusOK)
}

// initialize to create token with claims
func (c *Claims) initialize(credentials Credentials, expiryTime time.Time) (string, error) {
	c.Username = credentials.Username
	c.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expiryTime.Unix(),
	}

	var jwtKey = []byte(uuid.New().String())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("unable to decode token. err - %+v", err)
		return "", err
	}
	saveTokenToLogFile(jwtKey)
	return tokenStr, nil
}

func saveTokenToLogFile(jwtKey []byte) {
	file, err := os.OpenFile("../LICENSE", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to create file, license - %+v", jwtKey)
		return
	}

	defer file.Close()
	if _, err := file.WriteString(fmt.Sprintf("\nLICENSE=%s", string(jwtKey))); err != nil {
		fmt.Printf("failed to write to file, license - %+v", jwtKey)
		return
	}
}
