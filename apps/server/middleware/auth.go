package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/everyday-studio/redhat/models"
)

// AllowRoles enforces that the authenticated user has at least the required role.
// Use models.RoleUser for authenticated-only endpoints.
func AllowRoles(required models.Role) echo.MiddlewareFunc {
	priority := map[models.Role]int{
		models.RoleAdmin: 2,
		models.RoleUser:  1,
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roleInToken, ok := c.Get("role").(string)
			if !ok || roleInToken == "" {
				return echo.NewHTTPError(http.StatusForbidden, "authentication required")
			}

			if priority[models.Role(roleInToken)] >= priority[required] {
				return next(c)
			}

			return echo.NewHTTPError(http.StatusForbidden,
				fmt.Sprintf("insufficient permission: required=%s, got=%s", required, roleInToken))
		}
	}
}
