package routes

import (
	categoryHandler "boilerplate/internal/category"
	userHandler "boilerplate/internal/user"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	userHandler *userHandler.UserHandler,
	categoryHandler *categoryHandler.CategoryHandler,
	authMiddleware echo.MiddlewareFunc,
	adminMiddleware echo.MiddlewareFunc,
) {
	// Public routes
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	// Protected routes
	protected := e.Group("")
	protected.Use(authMiddleware)
	{
		// User routes
		protected.POST("/logout", userHandler.Logout)
		// users routes
		users := protected.Group("/admin/v1/user")
		{
			users.GET("/me", userHandler.GetMe)
			users.PUT("/update", userHandler.UpdateProfile)
			users.DELETE("/delete", userHandler.DeleteAccount)
			users.GET("", userHandler.GetAllUsers)
		}
		// Category routes
		categories := protected.Group("/admin/v1/categories")
		categories.Use(adminMiddleware)
		{
			categories.POST("", categoryHandler.Create)
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}
	}
}
