package services

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/everyday-studio/redhat/config"
	"github.com/everyday-studio/redhat/kit/security"
	"github.com/everyday-studio/redhat/models"
	"github.com/everyday-studio/redhat/repository/postgres"
)

const steamAuthURL = "https://api.steampowered.com/ISteamUserAuth/AuthenticateUserTicket/v1/"

// steamResponse mirrors the ISteamUserAuth/AuthenticateUserTicket response envelope.
type steamResponse struct {
	Response struct {
		Params *struct {
			Result          string `json:"result"`
			SteamID         string `json:"steamid"`
			OwnerSteamID    string `json:"ownersteamid"`
			VACBanned       bool   `json:"vacbanned"`
			PublisherBanned bool   `json:"publisherbanned"`
		} `json:"params"`
		Error *struct {
			ErrorCode int    `json:"errorcode"`
			ErrorDesc string `json:"errordesc"`
		} `json:"error"`
	} `json:"response"`
}

type AuthService struct {
	userRepo       *postgres.UserRepository
	privateKey     *rsa.PrivateKey
	accessTokenTTL time.Duration
	steamAPIKey    string
	steamAppID     int
	httpClient     *http.Client
}

func NewAuthService(
	userRepo *postgres.UserRepository,
	privateKey *rsa.PrivateKey,
	accessTokenTTL time.Duration,
	steamCfg config.SteamConfig,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		privateKey:     privateKey,
		accessTokenTTL: accessTokenTTL,
		steamAPIKey:    steamCfg.APIKey,
		steamAppID:     steamCfg.AppID,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
	}
}

// SteamLogin validates the Steam Session Ticket against the Steam Web API,
// upserts the user on success, and returns a signed access token.
func (s *AuthService) SteamLogin(ticket string, steamID string) (accessToken string, user *models.User, err error) {
	if ticket == "" || steamID == "" {
		return "", nil, fmt.Errorf("%w: ticket and steam_id are required", models.ErrInvalidInput)
	}

	if err := s.validateSteamTicket(ticket, steamID); err != nil {
		return "", nil, err
	}

	newID := uuid.New().String()
	user, err = s.userRepo.Upsert(newID, steamID, steamID)
	if err != nil {
		return "", nil, fmt.Errorf("SteamLogin upsert: %w", err)
	}

	accessToken, err = security.GenerateAccessToken(user.ID, string(models.RoleUser), s.privateKey, s.accessTokenTTL)
	if err != nil {
		return "", nil, fmt.Errorf("SteamLogin token generation: %w", err)
	}

	return accessToken, user, nil
}

// validateSteamTicket calls ISteamUserAuth/AuthenticateUserTicket and confirms
// that the authenticated SteamID matches the one the client claimed.
func (s *AuthService) validateSteamTicket(ticket string, claimedSteamID string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, steamAuthURL, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to build steam request: %v", models.ErrInternal, err)
	}

	q := req.URL.Query()
	q.Set("key", s.steamAPIKey)
	q.Set("appid", strconv.Itoa(s.steamAppID))
	q.Set("ticket", ticket)
	req.URL.RawQuery = q.Encode()

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: steam api unreachable: %v", models.ErrInternal, err)
	}
	defer resp.Body.Close()

	// Steam returns 403 HTML when the API key lacks permission for the app.
	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("%w: steam api key does not have permission for appid %d (need Publisher Web API Key)", models.ErrUnauthorized, s.steamAppID)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: steam api returned unexpected status %d", models.ErrInternal, resp.StatusCode)
	}

	var result steamResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("%w: failed to decode steam response: %v", models.ErrInternal, err)
	}

	if result.Response.Error != nil {
		return fmt.Errorf("%w: steam error %d: %s",
			models.ErrUnauthorized,
			result.Response.Error.ErrorCode,
			result.Response.Error.ErrorDesc,
		)
	}

	if result.Response.Params == nil || result.Response.Params.Result != "OK" {
		return fmt.Errorf("%w: steam ticket validation failed", models.ErrUnauthorized)
	}

	if result.Response.Params.SteamID != claimedSteamID {
		return fmt.Errorf("%w: steam_id mismatch (got %s, claimed %s)",
			models.ErrUnauthorized,
			result.Response.Params.SteamID,
			claimedSteamID,
		)
	}

	return nil
}
