package database

import (
	"context"
	"database/sql"
	"lila_games/internal/models"
	"strconv"
)

// Connection wraps a physical database connection
type Connection struct {
	db     *sql.DB
	closed bool
}

// New creates new databse connection
// implements internal.Store
// returns an error detailing the cause of failure in case the operation fails
func NewConnection(cfg models.DB) (*Connection, error) {
	connStr := "host=" + cfg.Host + " port=" + strconv.Itoa(cfg.Port) + " user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.Database + " sslmode=disable binary_parameters=yes"
	db, err := sql.Open(cfg.Driver, connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.MaxLifeTimeConn)
	return &Connection{db, false}, nil
}

// Migrate runs the migrations needs to bootstrap the game database
// returns an error detailing the cause of failure in case the operation fails
func (c *Connection) Migrate(ctx context.Context) error {
	if c.closed {
		return models.ErrDBClosed
	}
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	success := true
	defer func() {
		if success {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	for _, script := range models.Scripts {
		_, err = tx.ExecContext(ctx, script)
		if err != nil {
			success = false
			return err
		}
	}
	return nil
}

// GetTopModeByArea retrieves top modes of the game being played in a given area
// returns an error detailing the cause of failure in case the operation fails
func (c *Connection) GetTopModeByArea(ctx context.Context, gameId int, area string) ([]string, error) {
	if c.closed {
		return nil, models.ErrDBClosed
	}
	rows, err := c.db.QueryContext(ctx, `SELECT mode FROM (SELECT gm.mode as mode, count(*) AS count FROM games_being_played AS gbp 
	JOIN game_modes AS gm ON gbp.mode_id = gm.id 
	WHERE gbp.game_id = $1 AND gbp.area = $2 group by mode order by count desc) data`, gameId, area)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		mode  string
		modes []string = []string{}
	)
	for rows.Next() {
		err = rows.Scan(&mode)
		if err != nil {
			return nil, err
		}
		modes = append(modes, mode)
	}
	if len(modes) == 0 {
		return nil, models.ErrNotFound
	}
	return modes, rows.Err()
}

// AddPlayer adds a player to the game's provided mode in a given area
// returns an error detailing the cause of failure in case the operation fails
func (c *Connection) AddPlayer(ctx context.Context, gameId, modeId int, area string) (int, error) {
	if c.closed {
		return 0, models.ErrDBClosed
	}
	rows, err := c.db.QueryContext(ctx, "INSERT INTO games_being_played (game_id, mode_id, area) VALUES ($1,$2,$3) RETURNING id", gameId, modeId, area)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var id int
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	}
	return int(id), rows.Err()
}

// RemovePlayer removes the play from the game he is currently playing,
// returns an error detailing the cause of failure in case the operation fails
func (c *Connection) RemovePlayer(ctx context.Context, playerId int) error {
	if c.closed {
		return models.ErrDBClosed
	}
	rows, err := c.db.QueryContext(ctx, "DELETE FROM games_being_played WHERE id = $1 RETURNING id", playerId)
	if err != nil {
		return err
	}
	defer rows.Close()
	var id int
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
	}
	if int(id) != playerId {
		return models.ErrNotFound
	}
	return nil
}

func (c *Connection) Close() error {
	if c.closed {
		return models.ErrDBClosed
	}
	return c.db.Close()
}
