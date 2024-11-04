package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/user"
	"github.com/valentinesamuel/go_task-mgt-api/pkg"
	"net/http"
	"sync"
	"time"
)

type authHandlerImpl struct {
	repo user.UserRepository
	mu   sync.RWMutex
}

func NewAuthHandler(repo user.UserRepository) AuthHandler {
	if repo == nil {
		panic("repository cannot be nil")
	}
	return &authHandlerImpl{
		repo: repo,
	}
}

var validate = validator.New()

func (h *authHandlerImpl) RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var reqBody models.User
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest,
			pkg.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid request payload",
				err.Error()))
		return
	}

	if err := validate.Struct(&reqBody); err != nil {
		validationErrors := pkg.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest,
			pkg.NewErrorResponse(
				http.StatusBadRequest,
				"Validation errors",
				validationErrors))
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	exisingUser, _ := h.repo.GetByEmail(ctx, reqBody.Email)
	if exisingUser != nil {
		c.JSON(http.StatusConflict,
			pkg.NewErrorResponse(
				http.StatusConflict,
				"User already exists",
				"User with the provided email already exists"))
		return
	}

	// Hash the password before saving
	hashedPassword, err := pkg.Encrypt([]byte(reqBody.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			pkg.NewErrorResponse(
				http.StatusInternalServerError,
				"Internal server error",
				err.Error()))
		return
	}

	token, err := pkg.GenerateToken(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			pkg.NewErrorResponse(
				http.StatusInternalServerError,
				"Internal server error",
				err.Error()))
		return
	}

	reqBody.Password = hashedPassword
	reqBody.Token = token
	newUser, err := h.repo.Create(ctx, &reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			pkg.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed to create user",
				err.Error()))
		return
	}

	c.SetCookie("token", token, 36000, "/", "localhost", false, true)

	newUser.Password = ""
	c.JSON(http.StatusCreated, pkg.NewSuccessResponse(
		http.StatusCreated,
		"User created successfully",
		newUser))
}

func (h *authHandlerImpl) LoginUser(c *gin.Context) {
	panic("implement me")
}

func (h *authHandlerImpl) LogoutUser(c *gin.Context) {
	panic("implement me")
}
