package controllers

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {
    return c.String(http.StatusOK, "Create User Endpoint")
}

func LoginUser(c echo.Context) error {
    return c.String(http.StatusOK, "Login User Endpoint")
}

func GetUser(c echo.Context) error {
    return c.String(http.StatusOK, "Get User Endpoint")
}