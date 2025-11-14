package httpadapter

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"pr-manager-service/internal/domain"
	"pr-manager-service/internal/usecase"
)

// Constant strings for errors
const (
	errorCodeTeamExists  = "TEAM_EXISTS"
	errorCodePrExists    = "PR_EXISTS"
	errorCodePrMerged    = "PR_MERGED"
	errorCodeNotAssigned = "NOT_ASSIGNED"
	errorCodeNoCandidate = "NO_CANDIDATE"
	errorCodeNotFound    = "NOT_FOUND"
	errorCodeValidation  = "VALIDATION"
	errorCodeInternal    = "INTERNAL_ERROR"
)

// Write JSON to http response
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Return error for http request
func writeError(w http.ResponseWriter, status int, code, message string) {
	resp := errorResponseJSON{
		Error: errorBodyJSON{
			Code:    code,
			Message: message,
		},
	}
	writeJSON(w, status, resp)
}

func writeMappedError(w http.ResponseWriter, err error) {
	// Validation errors
	if errors.Is(err, usecase.ErrTeamNameRequired) ||
		errors.Is(err, usecase.ErrUserIdRequired) ||
		errors.Is(err, usecase.ErrPullRequestIdRequired) ||
		errors.Is(err, usecase.ErrPullRequestNameRequired) ||
		errors.Is(err, usecase.ErrAuthorIdRequired) ||
		errors.Is(err, usecase.ErrOldUserIdRequired) {
		writeError(w, http.StatusBadRequest, errorCodeValidation, err.Error())
		return
	}

	// TEAM_EXISTS
	if errors.Is(err, usecase.ErrTeamAlreadyExists) {
		writeError(w, http.StatusBadRequest, errorCodeTeamExists, err.Error())
		return
	}

	// PR_EXISTS
	if errors.Is(err, usecase.ErrPullRequestAlreadyExists) {
		writeError(w, http.StatusConflict, errorCodePrExists, err.Error())
		return
	}

	// PR_MERGED
	if errors.Is(err, domain.ErrEditMergedPR) {
		writeError(w, http.StatusConflict, errorCodePrMerged, err.Error())
		return
	}

	// NOT_ASSIGNED
	if errors.Is(err, usecase.ErrReviewerNotAssigned) {
		writeError(w, http.StatusConflict, errorCodeNotAssigned, err.Error())
		return
	}

	// NO_CANDIDATE
	if errors.Is(err, usecase.ErrNoCandidateInTeam) ||
		errors.Is(err, domain.ErrNoAvailableCandidates) {
		writeError(w, http.StatusConflict, errorCodeNoCandidate, err.Error())
		return
	}

	// NOT_FOUND
	if errors.Is(err, sql.ErrNoRows) ||
		errors.Is(err, usecase.ErrNotFound) {
		writeError(w, http.StatusNotFound, errorCodeNotFound, err.Error())
		return
	}

	// Other - internal
	writeError(w, http.StatusInternalServerError, errorCodeInternal, "internal error")
}
