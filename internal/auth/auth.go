package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/MrDarling002/PMAP/internal/models"
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

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        bearerToken := c.Request().Header.Get("Authorization")
        if bearerToken == "" {
            return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
        }
        if len(bearerToken) < 8 || bearerToken[:7] != "Bearer " {
            return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
        }
        tokenString := bearerToken[7:]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            secret := os.Getenv("JWT_SECRET")
            if secret == "" {
                return nil, fmt.Errorf("JWT_SECRET environment variable not set")
            }
            return []byte(secret), nil
        })
        if err != nil {
            return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
        }
        if !token.Valid {
            return echo.NewHTTPError(http.StatusUnauthorized, "Token is invalid")
        }
        return next(c)
    }
}
