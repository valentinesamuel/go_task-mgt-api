package auth

import (
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
	LogoutUser(c *gin.Context)
}
