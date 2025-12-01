package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pmii/pmii-backend/shared/jwt"
	sharedResponse "github.com/pmii/pmii-backend/shared/utils"
)

// AuthMiddleware validates JWT token
func AuthMiddleware(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			sharedResponse.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required", nil)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sharedResponse.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format", nil)
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			sharedResponse.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", err)
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserId)
		c.Set("user_email", claims.UserEmail)
		c.Set("user_level", claims.UserLevel)

		c.Next()
	}
}

// AdminOnlyMiddleware ensures only admin can access
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userLevel, exists := c.Get("user_level")
		if !exists || userLevel != "1" { // 1 = Admin
			sharedResponse.ErrorResponse(c, http.StatusForbidden, "Admin access only", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
