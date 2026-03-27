package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/everyday-studio/redhat/models"
	"github.com/everyday-studio/redhat/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(e *echo.Echo, authService *services.AuthService) *AuthHandler {
	h := &AuthHandler{authService: authService}

	api := e.Group("/api")
	api.POST("/auth/steam", h.SteamLogin)

	return h
}

// SteamLoginRequest is the request body for POST /api/auth/steam.
type SteamLoginRequest struct {
	Ticket  string `json:"ticket"   validate:"required"`
	SteamID string `json:"steam_id" validate:"required"`
}

// SteamLoginResponse is returned on a successful steam authentication.
type SteamLoginResponse struct {
	AccessToken string      `json:"access_token"`
	User        *models.User `json:"user"`
}

// SteamLogin authenticates a user via a Steam Session Ticket.
//
// POST /api/auth/steam
func (h *AuthHandler) SteamLogin(c echo.Context) error {
	var req SteamLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	if req.Ticket == "" || req.SteamID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "ticket and steam_id are required"})
	}

	accessToken, user, err := h.authService.SteamLogin(req.Ticket, req.SteamID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUnauthorized), errors.Is(err, models.ErrInvalidInput):
			return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, SteamLoginResponse{
		AccessToken: accessToken,
		User:        user,
	})
}
