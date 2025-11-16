package httpadapter

import (
	"net/http"

	"pr-manager-service/internal/usecase"
)

type HTTPHandler struct {
	svc     *usecase.Service
	appName string
	version string
}

func NewHTTPHandler(svc *usecase.Service, appName, version string) *HTTPHandler {
	return &HTTPHandler{
		svc:     svc,
		appName: appName,
		version: version,
	}
}

func NewRouter(svc *usecase.Service, appName, version string) *http.ServeMux {
	h := NewHTTPHandler(svc, appName, version)

	mux := http.NewServeMux()

	// Teams
	mux.HandleFunc("/team/add", h.handleCreateTeam)
	mux.HandleFunc("/team/get", h.handleGetTeam)

	// Users
	mux.HandleFunc("/users/setIsActive", h.handleSetIsActive)
	mux.HandleFunc("/users/getReview", h.handleGetUserReviews)

	// PullRequests
	mux.HandleFunc("/pullRequest/create", h.handleCreatePullRequest)
	mux.HandleFunc("/pullRequest/merge", h.handleMergePullRequest)
	mux.HandleFunc("/pullRequest/reassign", h.handleReassignReviewer)

	// Stats / Health
	mux.HandleFunc("/stats", h.handleStats)
	mux.HandleFunc("/health", h.handleHealth)

	return mux
}

// Wraps the handler to a new http server
func NewServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
