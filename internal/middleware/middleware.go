// Package middleware need fop an authorization in our requests
package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/distuurbia/firstTask/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware makes an authorization through access token
func JWTMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
			}

			tokenString := extractTokenFromHeader(authHeader)
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
			}
			token, err := ValidateToken(tokenString, cfg.SecretKey)
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				exp := claims["exp"].(float64)
				if exp < float64(time.Now().Unix()) {
					return echo.NewHTTPError(http.StatusUnauthorized, "Token is expired")
				}
			}
			return next(c)
		}
	}
}

// extractTokenFromHeader extractes token from authHeader
func extractTokenFromHeader(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(strings.ToLower(parts[0]), "bearer") {
		return ""
	}
	return parts[1]
}

// ValidateToken parses tokenString and checks if signing method is ok and return jwt token with filled Valid field
func ValidateToken(tokenString, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
