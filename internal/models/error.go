package models

import "errors"

// common errors
var (
	ErrDBClosed = errors.New("db is already closed")
	ErrNotFound = errors.New("not found")
)
