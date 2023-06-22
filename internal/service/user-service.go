// Package service realize bisnes-logic of the microservice
package service

import (
	"context"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
)

// UserRepository is an interface that contains CRUD methods and GetAll
type UserRepository interface {
	SignIn(ctx context.Context, user *model.User) (error)
	GetPassword(ctx context.Context, user *model.User) (string, error)
	AddRefreshToken(ctx context.Context, user *model.User) error
}

// UserService contains UserRepository interface
type UserService struct {
	userRps UserRepository
}

// NewUserService accepts UserRepository object and returnes an object of type *UserService
func NewUserService(userRps UserRepository) *UserService {
	return &UserService{userRps: userRps}
}

// SignIn is a method of UserService that calls  method of Repository
func (userSrv *UserService) SignIn(ctx context.Context, user *model.User) (error){
	var err error
	user.Password, err = HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = userSrv.userRps.SignIn(ctx, user)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// Login is a method of UserService that calls method of Repository
func (userSrv *UserService) Login(ctx context.Context, user *model.User)(string, string, error){
	hash, err := userSrv.userRps.GetPassword(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}
	verified, err := CheckPasswordHash(user.Password, hash)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}
	if verified {
		refreshToken, err := GetJWT(user, 72)
		if err != nil {
			return "", "", fmt.Errorf("%w", err)
		}
		user.RefreshToken = refreshToken
		// user.RefreshToken, err = HashPassword(refreshToken)
		// if err != nil {
		// 	return "", "", fmt.Errorf("%w", err)
		// }
		err = userSrv.userRps.AddRefreshToken(context.Background(), user)
		if err != nil {
			return "", "", fmt.Errorf("%w", err)
		}
		accessToken, err := GetJWT(user, 0.5)
		if err != nil {
			return "", "", fmt.Errorf("%w", err)
		}
		return accessToken, refreshToken, nil
	}
	return "", "", nil
}



