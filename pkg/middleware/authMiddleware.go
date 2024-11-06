package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the token from the user_token cookie
		token, err := c.Cookie("token")
		if err != nil {
			// If the cookie is not present, reject the request
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No token provided"})
			return
		}

		fmt.Println("Token: ", token)

		// Optionally, set user information in context for later handlers
		c.Set("user_token", token)

		// Token is valid, proceed to the next handler
		c.Next()
	}
}
