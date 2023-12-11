package middlewares

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/o1egl/paseto"
)

type CustomClaims struct {
	UserID int    `json:"uid"`
	Role   string `json:"role"`
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token is missing")
		}

		// Split the token string by space to get the actual token value
		tokenParts := strings.Fields(token)
		if len(tokenParts) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
		}

		// Extract the actual token value
		actualToken := tokenParts[1]

		var claims CustomClaims

		if err := paseto.NewV2().Decrypt(actualToken, []byte("YELLOW SUBMARINE, BLACK WIZARDRY"), &claims, nil); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		return next(c)
	}
}
