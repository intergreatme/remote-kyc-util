package routes

import (
	"github.com/gorilla/mux"

	"github.com/intergreatme/remote-kyc-util/handlers"
)

func Router(handler *handlers.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/allowlist", handler.AllowlistHandler).Methods("POST")
	router.HandleFunc("/api/completion", handler.CompletionHandler).Methods("POST")
	router.HandleFunc("/api/status", handler.StatusHandler).Methods("GET")

	return router
}
