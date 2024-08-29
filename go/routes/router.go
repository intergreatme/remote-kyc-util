/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package routes

import (
	"net/http"

	"github.com/intergreatme/remote-kyc-util/handlers"
)

func Router(handler *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.RootHandler)
	mux.HandleFunc("POST /api/allowlist", handler.AllowlistHandler)
	mux.HandleFunc("POST /api/cancel", handler.CancelHandler)
	mux.HandleFunc("POST /api/update", handler.UpdateHandler)
	mux.HandleFunc("POST /api/completion", handler.CompletionHandler)
	mux.HandleFunc("POST /api/status", handler.StatusHandler)
	mux.HandleFunc("POST /api/feedback", handler.FeedbackHandler)
	mux.HandleFunc("POST /api/getfile", handler.GetFile)
	mux.HandleFunc("POST /api/getlivelinessfile", handler.GetLivelinessFile)

	mux.HandleFunc("POST /api/verify", handler.VerifyHandler)

	return mux
}
