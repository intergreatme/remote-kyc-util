package handlers

import (
	"net/http"
)

func (h *Handler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the status logic here
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("StatusHandler response"))
}
