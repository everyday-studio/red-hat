package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/everyday-studio/redhat/middleware"
	"github.com/everyday-studio/redhat/models"
	"github.com/everyday-studio/redhat/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(e *echo.Echo, userService *services.UserService) *UserHandler {
	h := &UserHandler{userService: userService}

	api := e.Group("/api")
	users := api.Group("/users", middleware.AllowRoles(models.RoleUser))
	users.GET("/me", h.GetMe)

	return h
}

// GetMe returns the authenticated user's profile.
//
// GET /api/users/me
func (h *UserHandler) GetMe(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}

	return c.JSON(http.StatusOK, user)
}
