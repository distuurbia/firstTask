// Package handler contains handler methods and handler tests
package handler

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/distuurbia/firstTask/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
	validate *validator.Validate
}

// NewHandler accepts Service interface and returns an object of *HandlerEntity
func NewHandler(srvcPers PersonService, srvcUser UserService, validate *validator.Validate) *EntityHandler {
	return &EntityHandler{srvcPers: srvcPers, srvcUser: srvcUser, validate: validate}
}

// Create calls Create method of Service by handler
// @Summary Create a new person
// @Description Creates a new person
// @Tags Person
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for authentication"
// @Param personRequest body model.PersonRequest true "personRequest value (model.PersonRequest)"
// @Success 201 {object} model.Person
// @Failure 400 {object} error
// @Router /persons [post]
func (handl *EntityHandler) Create(c echo.Context) error {
	var createdPerson model.Person
	createdPerson.ID = uuid.New()
	err := c.Bind(&createdPerson)
	if err != nil {
		logrus.Errorf("EntityHandler -> Create -> c.Bind -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind")
	}
	err = handl.validate.StructCtx(c.Request().Context(), createdPerson)
	if err != nil {
		logrus.Errorf("EntityHandler -> Create -> validate -> StructCtx -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate")
	}
	err = handl.srvcPers.Create(c.Request().Context(), &createdPerson)
	if err != nil {
		logrus.Errorf("EntityHandler -> Create -> srvcPers.Create -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create check if id is UUID format")
	}
	return c.JSON(http.StatusCreated, createdPerson)
}

// ReadRow calls ReadRow method of Service by handler
// @Summary Get a person by ID
// @Description Get a person by ID
// @Tags Person
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for authentication"
// @Param id path string true "Person ID"
// @Success 200 {object} model.Person
// @Failure 400 {object} error
// @Router /persons/{id} [get]
func (handl *EntityHandler) ReadRow(c echo.Context) error {
	id := c.Param("id")
	err := handl.validate.VarCtx(c.Request().Context(), id, "required,uuid")
	if err != nil {
		logrus.Errorf("EntityHandler -> ReadRow -> validate -> VarCtx -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to valid id")
	}
	readPerson, err := handl.srvcPers.ReadRow(c.Request().Context(), uuid.MustParse(id))
	if err != nil {
		logrus.Errorf("EntityHandler -> ReadRow -> srvcPers.ReadRow -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read check if id is UUID format or that such person exist")
	}
	return c.JSON(http.StatusOK, readPerson)
}

// GetAll calls GetAll method of Service by handler
// @Summary Get all persons
// @Description Get all persons
// @Tags Person
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {array} model.Person
// @Failure 400 {object} error
// @Router /persons [get]
func (handl *EntityHandler) GetAll(c echo.Context) error {
	persAll, err := handl.srvcPers.GetAll(c.Request().Context())
	if err != nil {
		logrus.Errorf("EntityHandler -> GetAll -> srvcPers.GetAll -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to get all persons")
	}
	return c.JSON(http.StatusOK, persAll)
}

// Update calls Update method of Service by handler
// @Summary Update a person by ID
// @Description Update a person by ID
// @Tags Person
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for authentication"
// @Param id path string true "Person ID"
// @Param personRequest body model.PersonRequest true "personRequest value (model.PersonRequest)"
// @Success 200 {object} model.Person
// @Failure 400 {object} error
// @Router /persons/{id} [put]
func (handl *EntityHandler) Update(c echo.Context) error {
	var updatedPerson model.Person
	id := c.Param("id")
	err := handl.validate.VarCtx(c.Request().Context(), id, "required,uuid")
	if err != nil {
		logrus.Errorf("EntityHandler -> Update -> validate -> VarCtx -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to valid id")
	}
	updatedPerson.ID = uuid.MustParse(id)
	err = c.Bind(&updatedPerson)
	if err != nil {
		logrus.Errorf("EntityHandler -> Update -> c.Bind -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind")
	}
	err = handl.validate.StructCtx(c.Request().Context(), updatedPerson)
	if err != nil {
		logrus.Errorf("EntityHandler -> Update -> validate -> StructCtx -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate")
	}
	err = handl.srvcPers.Update(c.Request().Context(), &updatedPerson)
	if err != nil {
		logrus.Errorf("EntityHandler -> Update -> srvcPers.Update -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to update check if id is UUID format or that such person exist")
	}
	return c.JSON(http.StatusOK, updatedPerson)
}

// Delete calls Delete method of Service by handler
// @Summary Delete a person by ID
// @Description Delete a person by ID
// @Tags Person
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token for authentication"
// @Param id path string true "Person ID"
// @Success 204
// @Failure 400 {object} error
// @Router /persons/{id} [delete]
func (handl *EntityHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	err := handl.validate.VarCtx(c.Request().Context(), id, "required,uuid")
	if err != nil {
		logrus.Errorf("EntityHandler -> Update -> validate -> VarCtx -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to valid id")
	}
	err = handl.srvcPers.Delete(c.Request().Context(), uuid.MustParse(id))
	if err != nil {
		logrus.Errorf("EntityHandler -> Delete -> srvcPers.Delete -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "falled to delete check if id is UUID format or that such person exist")
	}
	return c.NoContent(http.StatusNoContent)
}

// SignUp calls SignUp method of Service by handler
// @Summary Sign up a new user
// @Description Sign up a new user
// @Tags User
// @Accept json
// @Produce json
// @Param userRequest body model.UserRequest true "userRequest value (model.PersonRequest)"
// @Success 201 {string} string
// @Failure 400 {object} error
// @Router /signUp [post]
func (handl *EntityHandler) SignUp(c echo.Context) error {
	bindInfo := struct {
		Username string `json:"username" validate:"required,min=4,max=15"`
		Password string `json:"password"`
	}{}
	err := c.Bind(&bindInfo)
	if err != nil {
		logrus.Errorf("EntityHandler -> SignUp -> c.Bind -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind")
	}
	var createdUser model.User
	createdUser.ID = uuid.New()
	createdUser.Username = bindInfo.Username
	createdUser.Password = []byte(bindInfo.Password)
	err = handl.validate.StructCtx(c.Request().Context(), createdUser)
	if err != nil {
		logrus.Errorf("EntityHandler -> SignUp -> validate.Struct -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate")
	}
	err = handl.srvcUser.SignUp(c.Request().Context(), &createdUser)
	if err != nil {
		logrus.Errorf("EntityHandler -> SignUp -> srvcUser.SignUp -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to signUp")
	}
	return c.JSON(http.StatusCreated, "ID: "+createdUser.ID.String())
}

// Login authenticates user and returns access and refresh tokens
// @Summary User login
// @Description Authenticates a user with the provided username and password
// @Tags User
// @Accept json
// @Produce json
// @Param userRequest body model.UserRequest true "userRequest value (model.PersonRequest)"
// @Success 200 {object} service.TokenPair
// @Failure 400 {object} error
// @Router /login [post]
func (handl *EntityHandler) Login(c echo.Context) error {
	bindInfo := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := c.Bind(&bindInfo)
	if err != nil {
		logrus.Errorf("EntityHandler -> Login -> c.Bind -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind", err)
	}
	var loginedUser model.User
	loginedUser.Username = bindInfo.Username
	loginedUser.Password = []byte(bindInfo.Password)
	err = handl.validate.StructCtx(c.Request().Context(), loginedUser)
	if err != nil {
		logrus.Errorf("EntityHandler -> SignUp -> validate.Struct -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate")
	}
	tokenPair, err := handl.srvcUser.Login(c.Request().Context(), &loginedUser)
	if err != nil {
		logrus.Errorf("EntityHandler -> Login -> srvcUser.Login -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to login")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"access token":  tokenPair.AccessToken,
		"refresh token": tokenPair.RefreshToken,
	})
}

// Refresh refreshes pair of access and refresh tokens
// @Summary Refreshes access and refresh tokens
// @Description Refreshes the access and refresh tokens using the provided refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Refresh token"
// @Success 200 {object} service.TokenPair
// @Failure 400 {object} error
// @Router /refresh [post]
func (handl *EntityHandler) Refresh(c echo.Context) error {
	bindInfo := struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{}
	err := c.Bind(&bindInfo)
	if err != nil {
		logrus.Errorf("EntityHandler -> Refresh -> c.Bind -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind")
	}
	var tokenPair service.TokenPair
	tokenPair.AccessToken = bindInfo.AccessToken
	tokenPair.RefreshToken = bindInfo.RefreshToken
	tokenPair, err = handl.srvcUser.Refresh(c.Request().Context(), tokenPair)
	if err != nil {
		logrus.Errorf("EntityHandler -> Refresh -> srvcUser.Refresh -> error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to refresh")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"access token":  tokenPair.AccessToken,
		"refresh token": tokenPair.RefreshToken,
	})
}

// DownloadImage downloads image from server
// @Summary Download an image
// @Description Downloads the specified image from the server
// @Tags Images
// @Accept json
// @Produce octet-stream
// @Param imageName path string true "Image filename"
// @Success 200 {file} octet-stream "Image file"
// @Failure 404 {object} error
// @Router /images/download/{imageName} [get]
func (handl *EntityHandler) DownloadImage(c echo.Context) error {
	imgName := c.Param("imageName")
	imgPath := "images/download/" + imgName
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		logrus.Errorf("EntityHandler -> DownloadImage -> os.Stat -> error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, "image not found")
	}
	imgPath = filepath.Clean(imgPath)
	img, err := os.Open(imgPath)
	if err != nil {
		logrus.Errorf("EntityHandler -> DownloadImage -> os.Open -> error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, "image not found")
	}
	defer func() {
		if err = img.Close(); err != nil {
			logrus.Errorf("EntityHandler -> DownloadImage -> img.Close -> error: %v", err)
		}
	}()
	c.Response().Header().Set("Content-Type", "image/png")
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+imgName)
	_, err = io.Copy(c.Response(), img)
	if err != nil {
		logrus.Errorf("EntityHandler -> DownloadImage -> io.copy -> error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, "image not found")
	}
	return nil
}

// UploadImage uploads image to server
// @Summary Upload an image
// @Description Uploads an image to the server
// @Tags Images
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file"
// @Success 200 {string} string "OK"
// @Failure 404 {object} error
// @Router /images/upload [post]
func (handl *EntityHandler) UploadImage(c echo.Context) error {
	image, err := c.FormFile("image")
	if err != nil {
		logrus.Errorf("EntityHandler -> UploadImage -> c.FormFile -> error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, "image not found")
	}
	src, err := image.Open()
	if err != nil {
		logrus.Errorf("EntityHandler -> UploadImage -> image.Open -> error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, "image not found")
	}
	defer func() {
		if err = src.Close(); err != nil {
			logrus.Errorf("EntityHandler -> uploadImage -> src.Close -> error: %v", err)
		}
	}()
	dstPath := "images/upload/" + image.Filename
	dstPath = filepath.Clean(dstPath)
	dst, err := os.Create(dstPath)
	if err != nil {
		logrus.Errorf("EntityHandler -> UploadImage -> os.Create -> error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, "failed to create")
	}
	defer func() {
		if err = dst.Close(); err != nil {
			logrus.Errorf("EntityHandler -> uploadImage -> dst.Close -> error: %v", err)
		}
	}()
	if _, err = io.Copy(dst, src); err != nil {
		logrus.Errorf("EntityHandler -> UploadImage -> io.Copy -> error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, "failed to copy")
	}
	c.Response().Header().Set("Content-Type", "image/png")
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+dst.Name())
	_, err = io.Copy(c.Response(), dst)
	if err != nil {
		logrus.Errorf("EntityHandler -> UploadImage -> io.Copy -> error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, "failed to copy")
	}
	return nil
}
