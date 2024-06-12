package handlers

import "database/sql"

// Handler struct to hold the dependencies
type Handler struct {
	DB     *sql.DB
	Config Config
}

type Config struct {
	ID             string
	PvtKeyPassword string
}

// NewHandler initializes and returns a Handler with the given dependencies
func NewHandler(db *sql.DB, config Config) *Handler {
	return &Handler{
		DB:     db,
		Config: config,
	}
}
