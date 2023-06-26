// Package service realize bisnes-logic of the microservice
package service

import (
	"context"
	"crypto/sha256"
	"time"

	"fmt"

	"github.com/distuurbia/firstTask/internal/config"
	"github.com/distuurbia/firstTask/internal/middleware"
	"github.com/distuurbia/firstTask/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository is an interface that contains CRUD methods and GetAll
type UserRepository interface {
	SignUp(ctx context.Context, user *model.User) error
	GetPasswordAndIDByUsername(ctx context.Context, username string) (uuid.UUID, []byte, error)
	AddRefreshToken(ctx context.Context, user *model.User) error
	GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (string, error)
}

// UserService contains UserRepository interface
type UserService struct {
	rpsUser UserRepository
}

// NewUserService accepts UserRepository object and returnes an object of type *UserService
func NewUserService(rpsUser UserRepository) *UserService {
	return &UserService{rpsUser: rpsUser}
}

// Expiration time for an access and a refresh tokens
const (
	accessTokenExpiration  = 15 * time.Minute
	refreshTokenExpiration = 72 * time.Hour
	bcryptCost             = 14
)

// TokenPair contains an Access and a refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// SignUp is a method of UserService that calls  method of Repository
func (srvUser *UserService) SignUp(ctx context.Context, user *model.User) error {
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
func (srvUser *UserService) Login(ctx context.Context, user *model.User) (TokenPair, error) {
	id, hash, err := srvUser.rpsUser.GetPasswordAndIDByUsername(ctx, user.Username)
	user.ID = id
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser ->  Login -> RepositoryUser -> GetPasswordByUsernsame -> error: %w", err)
	}
	verified, err := srvUser.CheckPasswordHash(hash, user.Password)
	if err != nil || !verified {
		return TokenPair{}, fmt.Errorf("ServiceUser ->  Login -> CheckPasswordHash -> error: %w", err)
	}
	tokenPair, err := srvUser.GenerateTokenPair(user.ID)
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser ->  Login -> GenerateTokenPair -> error: %w", err)
	}
	sum := sha256.Sum256([]byte(tokenPair.RefreshToken))
	hashedRefreshToken, err := srvUser.HashPassword(sum[:])
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser -> Login -> HashPassword -> error: %w", err)
	}
	user.RefreshToken = string(hashedRefreshToken)
	err = srvUser.rpsUser.AddRefreshToken(context.Background(), user)
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUsere ->  Login -> RepositoryUser -> AddRefreshToken -> error: %w", err)
	}
	return tokenPair, nil
}

// Refresh is a method of ServiceUser that refeshes access token and refresh token
func (srvUser *UserService) Refresh(ctx context.Context, tokenPair TokenPair) (TokenPair, error) {
	accessToken, err := middleware.ValidateToken(tokenPair.AccessToken, config.SecretKey)
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser -> Refresh -> accessToken -> middleware -> ValidateToken")
	}
	var accessID uuid.UUID
	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
		accessID = uuid.MustParse(claims["id"].(string))
	}
	refreshToken, err := middleware.ValidateToken(tokenPair.RefreshToken, config.SecretKey)
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser -> Refresh -> refreshToken -> middleware -> ValidateToken")
	}
	var refreshID uuid.UUID
	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		exp := claims["exp"].(float64)
		refreshID = uuid.MustParse(claims["id"].(string))
		if exp < float64(time.Now().Unix()) {
			return TokenPair{}, fmt.Errorf("ServiceUser ->  Refresh -> middleware -> ValidateToken -> error: %w", err)
		}
	}
	if accessID != refreshID {
		return TokenPair{}, fmt.Errorf("user ID in acess token doesn't equal user ID in refresh token")
	}
	hash, err := srvUser.rpsUser.GetRefreshTokenByID(ctx, accessID)
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser ->  Refresh -> RepositoryUser -> GetPasswordByUsernsame -> error: %w", err)
	}
	sum := sha256.Sum256([]byte(tokenPair.RefreshToken))
	verified, err := srvUser.CheckPasswordHash([]byte(hash), sum[:])
	if err != nil || !verified {
		return TokenPair{}, fmt.Errorf("ServiceUser ->  Refresh -> CheckPasswordHash -> error: refreshToken invalid")
	}
	tokenPair, err = srvUser.GenerateTokenPair(accessID)
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser ->  Refresh -> GenerateTokenPair -> error: %w", err)
	}
	sum = sha256.Sum256([]byte(tokenPair.RefreshToken))
	hashedRefreshToken, err := srvUser.HashPassword(sum[:])
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser -> Refresh -> HashPassword -> error: %w", err)
	}
	var user model.User
	user.RefreshToken = string(hashedRefreshToken)
	user.ID = accessID
	err = srvUser.rpsUser.AddRefreshToken(context.Background(), &user)
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUsere ->  Refresh -> RepositoryUser -> AddRefreshToken -> error: %w", err)
	}
	return tokenPair, nil
}

// HashPassword is a method of ServiceUser that makes from bytes hashed value
func (srvUser *UserService) HashPassword(password []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, bcryptCost)
	if err != nil {
		return bytes, fmt.Errorf("ServiceUser -> HashPassword -> GenerateFromPassword -> error: %w", err)
	}
	return bytes, nil
}

// CheckPasswordHash is a method of ServiceUser that checks if hash is equal hash from given password
func (srvUser *UserService) CheckPasswordHash(hash, password []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return false, fmt.Errorf("ServiceUser -> CheckPasswordHash -> CompareHashAndPassword -> error: %w", err)
	}
	return true, nil
}

// GenerateTokenPair generates pair of access and refresh tokens
func (srvUser *UserService) GenerateTokenPair(id uuid.UUID) (TokenPair, error) {
	accessToken, err := srvUser.GenerateJWTToken(accessTokenExpiration, id)
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser -> GenerateTokenPair -> accessToken -> GenerateJWTToken -> error: %w", err)
	}
	refreshToken, err := srvUser.GenerateJWTToken(refreshTokenExpiration, id)
	if err != nil {
		return TokenPair{}, fmt.Errorf("ServiceUser -> GenerateTokenPair -> refreshToken -> GenerateJWTToken -> error: %w", err)
	}
	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateJWTToken is a method of ServiceUser that generate JWT token with given expiration with user id
func (srvUser *UserService) GenerateJWTToken(expiration time.Duration, id uuid.UUID) (string, error) {
	claims := &jwt.MapClaims{
		"exp": time.Now().Add(expiration).Unix(),
		"id":  id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("ServiceUser -> GenerateJWTToken -> token.SignedString -> error: %w", err)
	}

	return tokenString, nil
}
