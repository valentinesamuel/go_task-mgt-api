package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/valentinesamuel/go_task-mgt-api/internal/user"
	"sync"
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

func (h *authHandlerImpl) Register(c *gin.Context) {
	panic("implement me")
}

func (h *authHandlerImpl) Login(c *gin.Context) {
	panic("implement me")
}

func (h *authHandlerImpl) Logout(c *gin.Context) {
	panic("implement me")
}
