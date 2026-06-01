//go:build norag

package handler

import (
	"net/http"
)

func ServeRAGSearch(w http.ResponseWriter, r *http.Request) {
	writeLocalizedErrorf(w, r, http.StatusServiceUnavailable, "RAGNotAvailable")
}

func ServeRAGMessage(w http.ResponseWriter, r *http.Request) {
	writeLocalizedErrorf(w, r, http.StatusServiceUnavailable, "RAGNotAvailable")
}

func ServeRAGSession(w http.ResponseWriter, r *http.Request) {
	writeLocalizedErrorf(w, r, http.StatusServiceUnavailable, "RAGNotAvailable")
}
