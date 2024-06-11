package routes

import (
	"net/http"

	"github.com/intergreatme/remote-kyc-util/handlers"
)

func Router() {
	http.HandleFunc("/api/allowlist", handlers.AllowlistHandler)

	// Add Status API, perhaps feedback?
}
