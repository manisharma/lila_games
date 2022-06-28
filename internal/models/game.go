package models

// Game represents a Game entity
type Game struct {
	Id             int      `json:"id"`
	Name           string   `json:"name"`
	SupportedModes []string `json:"modes"`
}

// GameBeingPlayed represents a player playing the game
type GameBeingPlayed struct {
	Id       int    `json:"id"`
	GameId   int    `json:"gameId"`
	PlayerId int    `json:"playerId"`
	Mode     string `json:"mode"`
	Area     string `json:"area"`
}
