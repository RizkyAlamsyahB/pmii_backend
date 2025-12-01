package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pmii/pmii-backend/shared/jwt"
	"github.com/pmii/user-service/internal/delivery/http/handler"
	"github.com/pmii/user-service/internal/delivery/http/middleware"
)

// SetupRouter configures all routes
func SetupRouter(userHandler *handler.UserHandler, jwtService *jwt.JWTService) *gin.Engine {
	router := gin.Default()

	// Apply global middlewares
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "user-service",
		})
	})

	// API v1 routes (base path /v1)
	v1 := router.Group("/v1")
	{
		// Public routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/logout", userHandler.Logout) // Logout (optional, client-side mainly)
		}

		// Authenticated routes under /v1/auth
		authProtected := v1.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware(jwtService))
		{
			authProtected.GET("/me", userHandler.GetProfile)
		}

		// Protected routes (auth required)
		authenticated := v1.Group("")
		authenticated.Use(middleware.AuthMiddleware(jwtService))
		{
			// User management - Admin only
			users := authenticated.Group("/users")
			users.Use(middleware.AdminOnlyMiddleware())
			{
				users.POST("", userHandler.Create)                            // Create user (Admin only)
				users.GET("", userHandler.GetAll)                              // List users (Admin only)
				users.GET("/:id", userHandler.GetById)                         // Get user by ID (Admin only)
				users.PUT("/:id", userHandler.Update)                          // Update user (Admin only)
				users.POST("/:id/change-password", userHandler.ChangePassword) // Change password (Admin only)
				users.DELETE("/:id", userHandler.Delete)                       // Delete user (Admin only)
			}
		}
	}

	return router
}
