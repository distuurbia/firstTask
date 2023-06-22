// Package handler contains handler methods and handler tests
package handler

import (
	"context"
	"net/http"

	"github.com/distuurbia/firstTask/internal/model"
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

type UserService interface {
	SignIn(ctx context.Context, user *model.User) (error)
	Login(ctx context.Context, user *model.User)(string, string, error)
}

// EntityHandler contains Service interface
type EntityHandler struct {
	persSrvc PersonService
	userSrvc UserService
}

// NewHandler accepts Service interface and returns an object of *PersonHandler
func NewHandler(persSrvc PersonService, userSrvc UserService) *EntityHandler {
	return &EntityHandler{persSrvc: persSrvc, userSrvc: userSrvc}
}

// Create calls Create method of Service by handler
func (handl *EntityHandler) Create(c echo.Context) error {
	var createdPerson model.Person
	createdPerson.ID = uuid.New()
	err := c.Bind(&createdPerson)
	if err != nil {
		return err
	}
	err = handl.persSrvc.Create(c.Request().Context(), &createdPerson)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create check if id is UUID format")
	}
	return c.JSON(http.StatusCreated, createdPerson)
}

// ReadRow calls ReadRow method of Service by handler
func (handl *EntityHandler) ReadRow(c echo.Context) error {
	id := uuid.MustParse(c.Param("id"))
	readPerson, err := handl.persSrvc.ReadRow(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read check if id is UUID format or that such person exist")
	}
	return c.JSON(http.StatusOK, readPerson)
}

// GetAll calls GetAll method of Service by handler
func (handl *EntityHandler) GetAll(c echo.Context) error {
	persAll, err := handl.persSrvc.GetAll(c.Request().Context())
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
	err = handl.persSrvc.Update(c.Request().Context(), &updatedPerson)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to update check if id is UUID format or that such person exist")
	}
	return c.JSON(http.StatusOK, updatedPerson)
}

// Delete calls Delete method of Service by handler
func (handl *EntityHandler) Delete(c echo.Context) error {
	id := uuid.MustParse(c.Param("id"))
	err := handl.persSrvc.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "falled to delete check if id is UUID format or that such person exist")
	}
	return c.NoContent(http.StatusNoContent)
}

// SignIn calls SignIn method of Service by handler
func (handl *EntityHandler) SignIn(c echo.Context) error {
	var createdUser model.User
	createdUser.ID = uuid.New()
	err := c.Bind(&createdUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind info")
	}
	err = handl.userSrvc.SignIn(c.Request().Context(), &createdUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to signIn")
	}
	return c.JSON(http.StatusCreated, "ID: " + createdUser.ID.String() + "Username: " + createdUser.Username)
}

// Login calls Login method of Service by handler
func (handl *EntityHandler) Login(c echo.Context) error {
	var loginedUser model.User
	err := c.Bind(&loginedUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind info")
	}
	accessToken, refreshToken, err := handl.userSrvc.Login(c.Request().Context(), &loginedUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to login")
	}

	return c.JSON(http.StatusOK, "access token: " + accessToken + "refresh token: " + refreshToken)
}
