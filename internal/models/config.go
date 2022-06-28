package models

import "time"

// Config wraps app level config
type Config struct {
	DB                                      DB
	ServerPort                              string        `default:"1234" split_words:"true"`
	TimeToWaitForSupportingServicesToComeUp time.Duration `default:"0s" split_words:"true"`
}

// DB represents the database config
type DB struct {
	Host            string        `required:"true" split_words:"true"`
	Port            int           `required:"true" split_words:"true"`
	User            string        `required:"true" split_words:"true"`
	Password        string        `required:"true" split_words:"true"`
	Database        string        `required:"true" split_words:"true"`
	Driver          string        `required:"true" split_words:"true"`
	MaxIdleConns    int           `default:"5" split_words:"true"`
	MaxOpenConns    int           `default:"10" split_words:"true"`
	MaxLifeTimeConn time.Duration `default:"1m" split_words:"true"`
	RunMigration    bool          `default:"true" split_words:"true"`
}
