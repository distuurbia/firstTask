package handler

import (
	"net/http"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/distuurbia/firstTask/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PersonHandler struct {
	srvc *service.PersonService
}

func NewHandler(srvc *service.PersonService) *PersonHandler {
	return &PersonHandler{srvc: srvc}
}
func (handl *PersonHandler) Create(c echo.Context) error {
	var createdPerson model.Person
	createdPerson.Id = uuid.New()
	err := c.Bind(&createdPerson)
	if err != nil {
		return err
	}
	err = handl.srvc.Create(c.Request().Context(), &createdPerson)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create check if id is UUID format")
	}
	return c.JSON(http.StatusCreated, createdPerson)

}

func (handl *PersonHandler) ReadRow(c echo.Context) error {
	id := uuid.MustParse(c.Param("id"))
	readPerson, err := handl.srvc.ReadRow(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read check if id is UUID format or that such person exist")
	}
	return c.JSON(http.StatusOK, readPerson)
}

func (handl *PersonHandler) Update(c echo.Context) error {
	var updatedPerson model.Person
	updatedPerson.Id = uuid.MustParse(c.Param("id"))
	err := c.Bind(&updatedPerson)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind info in person")
	}
	err = handl.srvc.Update(c.Request().Context(), &updatedPerson)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to update check if id is UUID format or that such person exist")
	}
	return c.JSON(http.StatusOK, updatedPerson)
}

func (handl *PersonHandler) Delete(c echo.Context) error {
	id := uuid.MustParse(c.Param("id"))
	err := handl.srvc.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "falled to delete check if id is UUID format or that such person exist")
	}
	return c.NoContent(http.StatusNoContent)
}
