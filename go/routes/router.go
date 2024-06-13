package routes

import (
	"net/http"

	"github.com/intergreatme/remote-kyc-util/handlers"
)

func Router(handler *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/allowlist", handler.AllowlistHandler)
	mux.HandleFunc("GET /api/completion", handler.CompletionHandler)
	mux.HandleFunc("GET /api/status", handler.StatusHandler)

	return mux
}
