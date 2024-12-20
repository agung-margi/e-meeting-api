package middleware

import (
	"e-meeting-api/pkg/util"
	"e-meeting-api/presenter/response"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var secret = os.Getenv("JWT_SECRET")

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("Unauthorized"))
		}

		token, err := util.VerifyToken(secret, tokenString)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("Unauthorized"))
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			c.Set("is_admin", claims["is_admin"])
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("Unauthorized"))

	}
}
