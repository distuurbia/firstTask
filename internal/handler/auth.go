// Package handler contains handler methods and handler tests
package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	//"github.com/golang-jwt/jwt"
	jwt "github.com/dgrijalva/jwt-go"
)

//
var mySigningKey = []byte(os.Getenv("SECURITY_KEY"))

// GetJWT generates JWT token
func GetJWT() (string, error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = "Alek"
	claims["aud"] = "billing.jwtgo.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", fmt.Errorf("something went wrong: %w", err)
	}
	return tokenString, nil
}

// Index do smth there
func Index(w http.ResponseWriter, r *http.Request) {
	validToken, err := GetJWT()
	fmt.Println(validToken)
	if err != nil {
		fmt.Println("failed to generate token: %w", err)
	}
	fmt.Println(w, string(validToken))
}

// 
func HandleRequests() {
	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

