package httpadapter

import (
	"net/http"

	"pr-manager-service/internal/usecase"
)

type HTTPHandler struct {
	svc *usecase.Service
}

func NewHTTPHandler(svc *usecase.Service) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

func NewRouter(svc *usecase.Service) *http.ServeMux {
	h := NewHTTPHandler(svc)

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

	// Health
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
