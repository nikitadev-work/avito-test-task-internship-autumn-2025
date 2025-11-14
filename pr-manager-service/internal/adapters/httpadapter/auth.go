package httpadapter

import (
	"errors"
	"net/http"
	"strings"
)

type authInfo struct {
	UserId  string
	IsAdmin bool
}

var (
	errNoAuthHeader   = errors.New("missing Authorization header")
	errInvalidAuthFmt = errors.New("invalid Authorization header format")
)

// Authorization: Bearer admin:<user_id>
// Authorization: Bearer user:<user_id>
func parseAuthHeader(r *http.Request) (*authInfo, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return nil, errNoAuthHeader
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, errInvalidAuthFmt
	}

	token := parts[1]
	tokenParts := strings.SplitN(token, ":", 2)
	if len(tokenParts) != 2 {
		return nil, errInvalidAuthFmt
	}

	role := tokenParts[0]
	userId := tokenParts[1]
	if userId == "" {
		return nil, errInvalidAuthFmt
	}

	info := &authInfo{
		UserId: userId,
	}

	switch role {
	case "admin":
		info.IsAdmin = true
	case "user":
		info.IsAdmin = false
	default:
		return nil, errInvalidAuthFmt
	}

	return info, nil
}

func requireAdmin(w http.ResponseWriter, r *http.Request) (*authInfo, bool) {
	info, err := parseAuthHeader(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, errorCodeNotFound, "admin token required")
		return nil, false
	}

	if !info.IsAdmin {
		writeError(w, http.StatusUnauthorized, errorCodeNotFound, "admin token required")
		return nil, false
	}

	return info, true
}

func requireAnyAuth(w http.ResponseWriter, r *http.Request) (*authInfo, bool) {
	info, err := parseAuthHeader(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, errorCodeNotFound, "auth token required")
		return nil, false
	}
	return info, true
}
