package handler

import (
	"net/http"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/distuurbia/firstTask/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Person struct {
	s *service.Person
}

func NewHandler(s *service.Person) *Person {
	return &Person{s: s}
}
func (h *Person) Create(c echo.Context) error {
	var createdPerson model.Person
	createdPerson.Id = uuid.New()
	err := c.Bind(&createdPerson)
	if err != nil {
		return err
	}
	err = h.s.Create(c.Request().Context(), &createdPerson)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "check if id is UUID format")
	}
	return c.JSON(http.StatusCreated, createdPerson)

}

func (h *Person) ReadRow(c echo.Context) error {
	id := uuid.MustParse(c.Param("id"))
	readPerson, err := h.s.ReadRow(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "check if id is UUID format or that such person exist")
	}
	return c.JSON(http.StatusOK, readPerson)
}

func (h *Person) Update(c echo.Context) error {
	var updatedPerson model.Person
	updatedPerson.Id = uuid.MustParse(c.Param("id"))
	err := c.Bind(&updatedPerson)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "check if id is UUID format or that such person exist")
	}
	err = h.s.Update(c.Request().Context(), &updatedPerson)
	return c.JSON(http.StatusOK, updatedPerson)
}

func (h *Person) Delete(c echo.Context) error {
	id := uuid.MustParse(c.Param("id"))
	err := h.s.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "check if id is UUID format or that such person exist")
	}
	return c.NoContent(http.StatusNoContent)
}
