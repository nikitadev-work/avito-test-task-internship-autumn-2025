package httpadapter

import (
	"encoding/json"
	"net/http"

	"pr-manager-service/internal/usecase"
)

// POST /pullRequest/create
func (h *HTTPHandler) handleCreatePullRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// only for admins
	if _, ok := requireAdmin(w, r); !ok {
		return
	}

	var req pullRequestCreateJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, errorCodeValidation, "invalid json")
		return
	}

	in := usecase.CreatePullRequestInput{
		PullRequestId:   req.PullRequestId,
		PullRequestName: req.PullRequestName,
		AuthorId:        req.AuthorId,
	}

	out, err := h.svc.CreatePullRequest(r.Context(), in)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	resp := pullRequestResponseJSON{
		PR: mapPullRequestDTOToJSON(out.PR),
	}

	writeJSON(w, http.StatusCreated, resp)
}

// POST /pullRequest/merge
func (h *HTTPHandler) handleMergePullRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// only for admins
	if _, ok := requireAdmin(w, r); !ok {
		return
	}

	var req pullRequestIdJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, errorCodeValidation, "invalid json")
		return
	}

	in := usecase.MergePullRequestInput{
		PullRequestId: req.PullRequestId,
	}

	out, err := h.svc.MergePullRequest(r.Context(), in)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	resp := pullRequestResponseJSON{
		PR: mapPullRequestDTOToJSON(out.PR),
	}

	writeJSON(w, http.StatusOK, resp)
}

// POST /pullRequest/reassign
func (h *HTTPHandler) handleReassignReviewer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// only for admins
	if _, ok := requireAdmin(w, r); !ok {
		return
	}

	var req reassignRequestJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, errorCodeValidation, "invalid json")
		return
	}

	in := usecase.ReassignReviewerInput{
		PullRequestId: req.PullRequestId,
		OldUserId:     req.OldUserId,
	}

	out, err := h.svc.ReassignReviewer(r.Context(), in)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	resp := reassignResponseJSON{
		PR:         mapPullRequestDTOToJSON(out.PR),
		ReplacedBy: out.ReplacedBy,
	}

	writeJSON(w, http.StatusOK, resp)
}

// GET /health
func (h *HTTPHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp := healthResponseJSON{
		Status: "ok",
	}
	writeJSON(w, http.StatusOK, resp)
}

func mapPullRequestDTOToJSON(pr usecase.PullRequestDTO) pullRequestJSON {
	return pullRequestJSON{
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		AuthorId:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
	}
}
