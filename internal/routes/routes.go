package routes

import (
    "github.com/MrDarling002/PMAP/internal/auth"
    "github.com/MrDarling002/PMAP/internal/controllers"

    "github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
    e.POST("/users", controllers.CreateUser)
    e.POST("/login", controllers.LoginUser)

    e.GET("/user", controllers.GetUser, auth.Middleware)
}