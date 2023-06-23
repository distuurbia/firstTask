// Package service realize bisnes-logic of the microservice
package service

import (
	"context"
	"os"
	"time"

	// "crypto/sha256"
	"fmt"

	"github.com/distuurbia/firstTask/internal/middleware"
	"github.com/distuurbia/firstTask/internal/model"
	"github.com/golang-jwt/jwt/v5"
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

// SignIn is a method of UserService that calls  method of Repository
func (userSrv *ServiceUser) SignUp(ctx context.Context, user *model.User) (error){
	var err error
	user.Password, err = userSrv.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("UserService -> HashPassword -> error: %w", err)
	}
	err = userSrv.rpsUser.SignUp(ctx, user)
	if err != nil {
		return fmt.Errorf("UserService -> UserRepository -> SignIn -> error: %w", err)
	}
	return nil
}

// Login is a method of UserService that calls method of Repository
func (userSrv *ServiceUser) Login(ctx context.Context, user *model.User)([]byte, []byte, error){
	hash, err := userSrv.rpsUser.GetPasswordByUsername(ctx, user)
	if err != nil {
		return []byte{}, []byte{}, fmt.Errorf("UserService -> UserRepository -> GetPasswordByUsernsame -> error: %w", err)
	}
	verified, err := userSrv.CheckPasswordHash(user.Password, hash)
	if err != nil || !verified {
		return []byte{}, []byte{}, fmt.Errorf("UserService -> CheckPasswordHash -> error: %w", err)
	}
	refreshToken, err := userSrv.GetJWT(user, 72)
	if err != nil {
		return []byte{}, []byte{}, fmt.Errorf("UserService -> GetJWT -> error: %w", err)
	}
	user.RefreshToken = refreshToken
	// user.RefreshToken = sha256.Sum256(refreshToken)
	err = userSrv.rpsUser.AddRefreshToken(context.Background(), user)
	if err != nil {
		return []byte{}, []byte{}, fmt.Errorf("UserService -> UserRepository -> AddRefreshToken -> error: %w", err)
	}
	accessToken, err := userSrv.GetJWT(user, 0.5)
	if err != nil {
		return []byte{}, []byte{}, fmt.Errorf("UserService -> GetJWT -> error: %w", err)
	}
	return accessToken, refreshToken, nil
}

func (userSrv *ServiceUser) HashPassword(password []byte) ([]byte, error) {
    bytes, err := bcrypt.GenerateFromPassword(password, 14)
	if err != nil {
		return bytes, fmt.Errorf("%w", err)
	}
    return bytes, nil
}

func (userSrv *ServiceUser) CheckPasswordHash(password, hash []byte) (bool, error) {
    err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil{
		return false, fmt.Errorf("%w", err)
	}
    return true, err
} 

//
func (userSrv *ServiceUser) GetJWT(user *model.User, expTime float64) ([]byte, error) {
	claims := &middleware.JWTCustomClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expTime))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return []byte{}, fmt.Errorf("%w", err)
	}
	return []byte(signedToken), nil
}


