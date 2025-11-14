package httpadapter

import (
	"encoding/json"
	"net/http"

	"pr-manager-service/internal/usecase"
)

// POST /users/setIsActive
func (h *HTTPHandler) handleSetIsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// only for admins
	if _, ok := requireAdmin(w, r); !ok {
		return
	}

	var req setIsActiveRequestJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, errorCodeValidation, "invalid json")
		return
	}

	in := usecase.SetIsActiveInput{
		UserId:   req.UserId,
		IsActive: req.IsActive,
	}

	out, err := h.svc.SetIsActive(r.Context(), in)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	resp := setIsActiveResponseJSON{
		User: userJSON{
			UserId:   out.UserId,
			Username: out.UserName,
			TeamName: out.TeamName,
			IsActive: out.IsActive,
		},
	}

	writeJSON(w, http.StatusOK, resp)
}

// GET /users/getReview?user_id=...
func (h *HTTPHandler) handleGetUserReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// for users or admins
	auth, ok := requireAnyAuth(w, r)
	if !ok {
		return
	}

	userId := r.URL.Query().Get("user_id")

	// only admin can get reviews for other users
	if !auth.IsAdmin && auth.UserId != userId {
		writeError(w, http.StatusUnauthorized, errorCodeNotFound, "forbidden for this user_id")
		return
	}

	in := usecase.GetUserReviewsInput{UserId: userId}

	out, err := h.svc.GetUserReviews(r.Context(), in)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	resp := userReviewsResponseJSON{
		UserId:       out.UserId,
		PullRequests: mapPullRequestShortDTOsToJSON(out.PullRequests),
	}

	writeJSON(w, http.StatusOK, resp)
}

func mapPullRequestShortDTOsToJSON(in []usecase.PullRequestShortDTO) []pullRequestShortJSON {
	result := make([]pullRequestShortJSON, 0, len(in))
	for _, p := range in {
		result = append(result, pullRequestShortJSON{
			PullRequestId:   p.PullRequestId,
			PullRequestName: p.PullRequestName,
			AuthorId:        p.AuthorId,
			Status:          p.Status,
		})
	}
	return result
}
