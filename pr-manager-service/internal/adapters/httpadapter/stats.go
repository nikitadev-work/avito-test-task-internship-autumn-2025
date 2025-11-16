package httpadapter

import (
	"net/http"
	"time"
)

// GET /stats
func (h *HTTPHandler) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resp := map[string]any{
		"service": h.appName,
		"version": h.version,
		"time":    time.Now().UTC().Format(time.RFC3339),
	}

	writeJSON(w, http.StatusOK, resp)
}
