package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type JWTConfig struct {
	SigningKey string
}
type JWTCustomClaims struct {
	ID  uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}


func JWTAuth(signingKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			token := extractTokenFromHeader(authHeader)
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
			}
			claims, err := validateToken(token, []byte(signingKey))
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			c.Set("ID", claims.ID)
			return next(c)
		}
		
	}
}

func extractTokenFromHeader(authHeader string) string{
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func validateToken(tokenString string, SigningKey []byte) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SigningKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(JWTCustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid token claims")
	}
	return &claims, nil
}