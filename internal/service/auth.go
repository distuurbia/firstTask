// Package handler contains handler methods and handler tests
package service

import (
	"fmt"
	"os"
	"time"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


var mySigningKey = []byte(os.Getenv("SECRET_KEY"))

type jwtCustomClaims struct {
	ID  uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return string(bytes), fmt.Errorf("%w", err)
	}
    return string(bytes), nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil{
		return false, fmt.Errorf("%w", err)
	}
    return true, err
} 

func GetJWT(user *model.User, expTime float64) (string, error) {
	claims := &jwtCustomClaims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expTime))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(mySigningKey)
	if err != nil {
		return signedToken, fmt.Errorf("%w", err)
	}
	return signedToken, nil
}

// func SignIn() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		username := c.FormValue("username")
// 		password := c.FormValue("password")

// 		if username == "you" && password == "passed" {
// 			token, err := GenerateToken()

// 			if err != nil {
// 				return err
// 			}

// 			return c.JSON(http.StatusOK, Token{
// 				Token: token,
// 			})
// 		}

// 		return echo.ErrUnauthorized
// 	}
// }

// func GenerateToken() (string, error) {
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	claims := token.Claims.(jwt.MapClaims)
  
// 	claims["authorized"] = true
// 	claims["client"] = "Krissanawat"
// 	claims["aud"] = "billing.jwtgo.io"
// 	claims["iss"] = "jwtgo.io"
// 	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
  
// 	tokenString, err := token.SignedString(mySigningKey)
  
// 	if err != nil {
// 		return "", fmt.Errorf("Something Went Wrong: %w", err)
// 	}
  
// 	return tokenString, nil
// }