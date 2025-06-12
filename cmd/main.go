package main

import (
	"fmt"
	"log"

	"boilerplate/config"
	"boilerplate/container"
	"boilerplate/internal/category"
	"boilerplate/internal/user"
	"boilerplate/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize container
	ctn, err := container.NewContainer()
	if err != nil {
		log.Fatal("Cannot initialize container:", err)
	}

	// Get echo instance
	e := ctn.Get(container.EchoDefName).(*echo.Echo)

	// Get handlers
	userHandler := ctn.Get(container.UserHandlerDefName).(*user.UserHandler)
	categoryHandler := ctn.Get(container.CategoryHandlerDefName).(*category.CategoryHandler)

	// Get middleware
	authMiddleware := ctn.Get(container.AuthMiddlewareDefName).(echo.MiddlewareFunc)
	adminMiddleware := ctn.Get(container.AdminAuthMiddlewareDefName).(echo.MiddlewareFunc)

	// Setup routes
	routes.SetupRoutes(e, userHandler, categoryHandler, authMiddleware, adminMiddleware)

	// Get config and start server
	cfg := ctn.Get(container.ConfigDefName).(config.Config)
	port := fmt.Sprintf(":%s", cfg.ServerPort)
	if err := e.Start(port); err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
