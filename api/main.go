package main

import (
    "log"
    "github.com/MrDarling002/PMAP/db"
    "github.com/MrDarling002/PMAP/internal/routes"
    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    db.InitDB()
    routes.SetupRoutes(e)
    log.Fatal(e.Start(":8080"))
}