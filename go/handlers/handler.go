/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package handlers

import (
	"database/sql"

	"github.com/intergreatme/remote-kyc-util/config"
)

// Handler struct to hold the dependencies
type Handler struct {
	DB     *sql.DB
	Config config.Configuration
}

// NewHandler initializes and returns a Handler with the given dependencies
func NewHandler(db *sql.DB, config config.Configuration) *Handler {
	return &Handler{
		DB:     db,
		Config: config,
	}
}
