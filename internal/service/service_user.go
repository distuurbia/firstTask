// Package service realize bisnes-logic of the microservice
package service

import (
	"context"
	"crypto/sha256"
	"os"
	"time"

	// "crypto/sha256"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RepositoryUser is an interface that contains CRUD methods and GetAll
type RepositoryUser interface {
	SignUp(ctx context.Context, user *model.User) (error)
	GetPasswordByUsername(ctx context.Context, user *model.User) ([]byte, error)
	AddRefreshToken(ctx context.Context, user *model.User) error
}

// ServiceUser contains UserRepository interface
type ServiceUser struct {
	rpsUser RepositoryUser
}


// NewUserService accepts UserRepository object and returnes an object of type *UserService
func NewServiceUser(rpsUser RepositoryUser) *ServiceUser {
	return &ServiceUser{rpsUser: rpsUser}
}

// Expiration time for an access and a refresh tokens
const (
	accessTokenExpiration  = 15 * time.Minute
	refreshTokenExpiration = 72 * time.Hour
)

// TokenPair contains an Access and a refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type JWTCustomClaims struct {
	ID  uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

// SignIn is a method of UserService that calls  method of Repository
func (srvUser *ServiceUser) SignUp(ctx context.Context, user *model.User) (error){
	var err error
	user.Password, err = srvUser.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("ServiceUser -> HashPassword -> error: %w", err)
	}
	err = srvUser.rpsUser.SignUp(ctx, user)
	if err != nil {
		return fmt.Errorf("ServiceUser -> UserRepository -> SignIn -> error: %w", err)
	}
	return nil
}

// Login is a method of UserService that calls method of Repository
func (srvUser *ServiceUser) Login(ctx context.Context, user *model.User)(string, string, error){
	hash, err := srvUser.rpsUser.GetPasswordByUsername(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("ServiceUser ->  Login -> RepositoryUser -> GetPasswordByUsernsame -> error: %w", err)
	}
	verified, err := srvUser.CheckPasswordHash(hash, user.Password)
	if err != nil || !verified {
		return "", "", fmt.Errorf("ServiceUser ->  Login -> CheckPasswordHash -> error: %w", err)
	}
	tokenPair, err := srvUser.GenerateTokenPair()
	if err != nil {
		return "", "", fmt.Errorf("ServiceUser ->  Login -> GenerateTokenPair -> error: %w", err)
	}
	sum := sha256.Sum256([]byte(tokenPair.RefreshToken))
	hashedRefreshToken, err := srvUser.HashPassword(sum[:])
	if err != nil {
		return "", "", fmt.Errorf("ServiceUser -> Login -> HashPassword -> error: %w", err)
	}
	user.RefreshToken = string(hashedRefreshToken)

	err = srvUser.rpsUser.AddRefreshToken(context.Background(), user)
	if err != nil {
		return "", "", fmt.Errorf("ServiceUsere ->  Login -> RepositoryUser -> AddRefreshToken -> error: %w", err)
	}
	return tokenPair.AccessToken, user.RefreshToken, nil
}

func (srvUser *ServiceUser) HashPassword(password []byte) ([]byte, error) {
    bytes, err := bcrypt.GenerateFromPassword(password, 14)
	if err != nil {
		return bytes, fmt.Errorf("%w", err)
	}
    return bytes, nil
}

func (srvUser *ServiceUser) CheckPasswordHash(hash, password []byte) (bool, error) {
    err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil{
		return false, fmt.Errorf("%w", err)
	}
    return true, nil
} 

func (srvUser *ServiceUser) GenerateTokenPair() (TokenPair, error) {
	accessToken, err := srvUser.GenerateJWTToken(accessTokenExpiration)
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser -> GenerateTokenPair -> accessToken -> GenerateJWTToken -> error: %w", err)
	}

	refreshToken, err := srvUser.GenerateJWTToken(refreshTokenExpiration)
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser -> GenerateTokenPair -> refreshToken -> GenerateJWTToken -> error: %w", err)
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (srvUser *ServiceUser) GenerateJWTToken(expiration time.Duration) (string, error) {
	claims := &jwt.MapClaims{
		"exp": time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("ServiceUser -> GenerateJWTToken -> token.SignedString -> error: %w", err)
	}

	return tokenString, nil
}


