// Package handler contains handler methods and handler tests
package handler

import (
	"context"
	"net/http"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/distuurbia/firstTask/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// PersonService is an interface that contains CRUD methods and GetAll of the service
type PersonService interface {
	Create(ctx context.Context, pers *model.Person) error
	ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error)
	GetAll(ctx context.Context) ([]model.Person, error)
	Update(ctx context.Context, pers *model.Person) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// UserService is an interface that contains methods of service for user
type UserService interface {
	SignUp(ctx context.Context, user *model.User) error
	Login(ctx context.Context, user *model.User) (service.TokenPair, error)
	Refresh(ctx context.Context, tokenPair service.TokenPair) (service.TokenPair, error)
}

// EntityHandler contains Service interface
type EntityHandler struct {
	srvcPers PersonService
	srvcUser UserService
}

// NewHandler accepts Service interface and returns an object of *HandlerEntity
func NewHandler(srvcPers PersonService, srvcUser UserService) *EntityHandler {
	return &EntityHandler{srvcPers: srvcPers, srvcUser: srvcUser}
}

// Create calls Create method of Service by handler
func (handl *EntityHandler) Create(c echo.Context) error {
	var createdPerson model.Person
	createdPerson.ID = uuid.New()
	err := c.Bind(&createdPerson)
	if err != nil {
		return err
	}
	err = handl.srvcPers.Create(c.Request().Context(), &createdPerson)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create check if id is UUID format")
	}
	return c.JSON(http.StatusCreated, createdPerson)
}

// ReadRow calls ReadRow method of Service by handler
func (handl *EntityHandler) ReadRow(c echo.Context) error {
	id := uuid.MustParse(c.Param("id"))
	readPerson, err := handl.srvcPers.ReadRow(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read check if id is UUID format or that such person exist")
	}
	return c.JSON(http.StatusOK, readPerson)
}

// GetAll calls GetAll method of Service by handler
func (handl *EntityHandler) GetAll(c echo.Context) error {
	persAll, err := handl.srvcPers.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to get all persons")
	}
	return c.JSON(http.StatusOK, persAll)
}

// Update calls Update method of Service by handler
func (handl *EntityHandler) Update(c echo.Context) error {
	var updatedPerson model.Person
	updatedPerson.ID = uuid.MustParse(c.Param("id"))
	err := c.Bind(&updatedPerson)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind info in person")
	}
	err = handl.srvcPers.Update(c.Request().Context(), &updatedPerson)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to update check if id is UUID format or that such person exist")
	}
	return c.JSON(http.StatusOK, updatedPerson)
}

// Delete calls Delete method of Service by handler
func (handl *EntityHandler) Delete(c echo.Context) error {
	id := uuid.MustParse(c.Param("id"))
	err := handl.srvcPers.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "falled to delete check if id is UUID format or that such person exist")
	}
	return c.NoContent(http.StatusNoContent)
}

// SignIn calls SignIn method of Service by handler
func (handl *EntityHandler) SignIn(c echo.Context) error {
	bindInfo := struct {
		Name string `json:"username"`
		Pass string `json:"password"`
	}{}
	err := c.Bind(&bindInfo)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind info")
	}
	var createdUser model.User
	createdUser.ID = uuid.New()
	createdUser.Username = bindInfo.Name
	createdUser.Password = []byte(bindInfo.Pass)
	err = handl.srvcUser.SignUp(c.Request().Context(), &createdUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to signIn")
	}
	return c.JSON(http.StatusCreated, "ID: "+createdUser.ID.String()+"Username: "+createdUser.Username)
}

// Login calls Login method of Service by handler
func (handl *EntityHandler) Login(c echo.Context) error {
	bindInfo := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := c.Bind(&bindInfo)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind info")
	}
	var loginedUser model.User
	loginedUser.Username = bindInfo.Username
	loginedUser.Password = []byte(bindInfo.Password)
	tokenPair, err := handl.srvcUser.Login(c.Request().Context(), &loginedUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to login")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access token":  tokenPair.AccessToken,
		"refresh token": tokenPair.RefreshToken,
	})
}

// Refresh refreshes pair of access and refresh tokens
func (handl *EntityHandler) Refresh(c echo.Context) error {
	bindInfo := struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{}
	err := c.Bind(&bindInfo)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind info")
	}
	var tokenPair service.TokenPair
	tokenPair.AccessToken = bindInfo.AccessToken
	tokenPair.RefreshToken = bindInfo.RefreshToken
	tokenPair, err = handl.srvcUser.Refresh(c.Request().Context(), tokenPair)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to refresh")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"access token":  tokenPair.AccessToken,
		"refresh token": tokenPair.RefreshToken,
	})
}
