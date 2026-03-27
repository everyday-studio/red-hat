package models

import "time"

type Role string

const (
	RoleUser  Role = "User"
	RoleAdmin Role = "Admin"
)

type User struct {
	ID        string    `json:"id"`
	SteamID   string    `json:"steam_id"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
