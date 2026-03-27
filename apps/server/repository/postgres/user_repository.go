package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/everyday-studio/redhat/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	const query = `
		SELECT id, steam_id, nickname, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.SteamID, &user.Nickname, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("GetByID: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetBySteamID(steamID string) (*models.User, error) {
	const query = `
		SELECT id, steam_id, nickname, created_at, updated_at
		FROM users
		WHERE steam_id = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, steamID).Scan(
		&user.ID, &user.SteamID, &user.Nickname, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("GetBySteamID: %w", err)
	}
	return user, nil
}

// Upsert inserts a new user using the caller-provided id.
// On conflict (existing steam_id) it only bumps updated_at and returns the existing row.
// UUID generation is the responsibility of the service layer.
func (r *UserRepository) Upsert(id string, steamID string, nickname string) (*models.User, error) {
	const query = `
		INSERT INTO users (id, steam_id, nickname)
		VALUES ($1, $2, $3)
		ON CONFLICT (steam_id) DO UPDATE
			SET updated_at = NOW()
		RETURNING id, steam_id, nickname, created_at, updated_at`

	user := &models.User{}
	err := r.db.QueryRow(query, id, steamID, nickname).Scan(
		&user.ID, &user.SteamID, &user.Nickname, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			return nil, fmt.Errorf("Upsert pq error (%s): %w", pqErr.Code, err)
		}
		return nil, fmt.Errorf("Upsert: %w", err)
	}
	return user, nil
}
