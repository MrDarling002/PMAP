package models

import (
    "fmt"
    "net/http"
    "time"
	"os"
    "github.com/golang-jwt/jwt/v4"
    "github.com/labstack/echo/v4"
)

func GenerateToken(user *models.User) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    })
    secret := os.Getenv("JWT_SECRET")
    return token.SignedString([]byte(secret))
}

func Middleware(c echo.Context) error { bearerToken =
    bearerToken := c.Request().Header.Get("Authorization")
    if bearerToken == "" {
        return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
    }
    tokenString := bearerToken[7:]
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil
			fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    if err != nil {
        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
    }
    if !token.Valid {
        return echo.NewHTTPError(http.StatusUnauthorized, "Token is invalid")
    }
    return c.Next()
}