package controllers

import (
    "net/http"
    "github.com/MrDarling002/PMAP/internal/auth"
    "github.com/MrDarling002/PMAP/internal/models"
    "github.com/MrDarling002/PMAP/db"

    "golang.org/x/crypto/bcrypt"
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt/v4"
)

func CreateUser(c echo.Context) error {
    user := new(models.User)
    if err := c.Bind(user); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }
    if user.Email == "" || user.Username == "" || user.Password == "" {
        return echo.NewHTTPError(http.StatusBadRequest, "Email, username, and password are required")
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
    }
    user.Password = string(hashedPassword)
    result, err := db.DB.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
    }
    userID, err := result.LastInsertId()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user ID")
    }
    user.ID = int(userID)

    return c.JSON(http.StatusCreated, user)
}

func LoginUser(c echo.Context) error {
    user := new(models.User)
    if err := c.Bind(user); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }
   if user.Email == "" || user.Password == "" {
        return echo.NewHTTPError(http.StatusBadRequest, "Email and password are required")
    }
    var storedUser models.User
    err := db.DB.QueryRow("SELECT id, username, email, password FROM users WHERE email = $1", user.Email).Scan(&storedUser.ID, &storedUser.Username, &storedUser.Email, &storedUser.Password)
    if err != nil {
        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
    }
    if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
    }
    token, err := auth.GenerateToken(&storedUser)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
    }

    return c.JSON(http.StatusOK, map[string]string{
        "token": token,
    })
}

func GetUser(c echo.Context) error {
    userID := c.Get("user").(jwt.MapClaims)["sub"].(float64)
    var user models.User
    err := db.DB.QueryRow("SELECT id, username, email FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Username, &user.Email)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "User not found")
    }
    return c.JSON(http.StatusOK, user)
}

func UpdateUser(c echo.Context) error {
    userID := c.Get("user").(jwt.MapClaims)["sub"].(float64)
    user := new(models.User)
    if err := c.Bind(user); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }
    if user.Email == "" || user.Username == "" {
        return echo.NewHTTPError(http.StatusBadRequest, "Email and username are required")
    }
    _, err := db.DB.Exec("UPDATE users SET username = $1, email = $2 WHERE id = $3", user.Username, user.Email, userID)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
    }

    return c.NoContent(http.StatusOK)
}

func DeleteUser(c echo.Context) error {
    userID := c.Get("user").(jwt.MapClaims)["sub"].(float64)
    _, err := db.DB.Exec("DELETE FROM users WHERE id = $1", userID)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
    }
    return c.NoContent(http.StatusOK)
}