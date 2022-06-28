package models

// databse migration/bootstrapping scripts
const (
	createTableScript_Games = `CREATE TABLE IF NOT EXISTS games (
		id 			SERIAL PRIMARY KEY,
		name 		VARCHAR UNIQUE NOT NULL)`
	createTableScript_GameModes = `CREATE TABLE IF NOT EXISTS game_modes (
		id 			SERIAL PRIMARY KEY,
		game_id 	INTEGER NOT NULL REFERENCES games(id) ON DELETE CASCADE,
		mode 		VARCHAR NOT NULL,
		CONSTRAINT constraint_uinque_game_mode UNIQUE (game_id, mode))`
	createTableScript_GamesBeingPlayed = `CREATE TABLE IF NOT EXISTS games_being_played (
		id 			SERIAL PRIMARY KEY,
		game_id 	INTEGER NOT NULL REFERENCES games(id) ON DELETE CASCADE,
		mode_id 	INTEGER NOT NULL REFERENCES game_modes(id) ON DELETE CASCADE,
		area 		CHAR(3) NOT NULL)`
	createIndexScript_GamesBeingPlayed     = `CREATE INDEX IF NOT EXISTS games_being_played_idx ON games_being_played (game_id, area)`
	truncateTableScript_Games              = `DELETE FROM games`
	resetSerialScript_Games                = `ALTER SEQUENCE games_id_seq RESTART WITH 1`
	resetSerialScript_GamesMode            = `ALTER SEQUENCE game_modes_id_seq RESTART WITH 1`
	resetSerialScript_GamesBeingPlayedMode = `ALTER SEQUENCE game_modes_id_seq RESTART WITH 1`
	insertTableScript_Games                = `INSERT INTO games (name) VALUES ('FAMOUS SHOOTER')`
	insertTableScript_GamesMode            = `INSERT INTO game_modes (game_id, mode) VALUES (1, 'Battle Royal'), (1, 'Team Deathmatch'), (1, 'Capture the flag')`
)

// migration scripts
var Scripts []string = []string{
	createTableScript_Games,
	createTableScript_GameModes,
	createTableScript_GamesBeingPlayed,
	createIndexScript_GamesBeingPlayed,
	truncateTableScript_Games,
	resetSerialScript_Games,
	resetSerialScript_GamesMode,
	resetSerialScript_GamesBeingPlayedMode,
	insertTableScript_Games,
	insertTableScript_GamesMode,
}
