package middleware

import (
	"e-meeting-api/pkg/util"
	"e-meeting-api/presenter/response"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("Unauthorized"))
		}

		// Memeriksa apakah header Authorization dimulai dengan "Bearer "
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("Invalid token format"))
		}

		// Mengambil token yang ada setelah "Bearer " (menghapus kata "Bearer " dari string)
		tokenString = tokenString[7:]

		// Verifikasi token
		token, err := util.VerifyToken(tokenString)
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

func IsAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isAdmin := c.Get("is_admin").(bool)
		if !isAdmin {
			return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("Unauthorized"))
		}
		return next(c)
	}
}
