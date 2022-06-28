package internal

import (
	"context"
)

// Store represents a common interface to database [in-memory or physical]
type Store interface {
	// Migrate is intended to run migrations scripts in any
	Migrate(ctx context.Context) error
	// GetTopModeByArea retrieves top modes of a game being played in the given area
	GetTopModeByArea(ctx context.Context, gameId int, area string) ([]string, error)
	// AddPlayer adds new plays to game, returns the player id
	AddPlayer(ctx context.Context, gameId, modeId int, area string) (int, error)
	// RemovePlayer removes the player from the game
	RemovePlayer(ctx context.Context, playerId int) error
	// Close closes database connection
	Close() error
}
