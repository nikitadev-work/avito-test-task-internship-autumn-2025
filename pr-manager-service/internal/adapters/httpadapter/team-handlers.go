package httpadapter

import (
	"encoding/json"
	"net/http"

	"pr-manager-service/internal/usecase"
)

// POST /team/add
func (h *HTTPHandler) handleCreateTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req teamJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, errorCodeValidation, "invalid json")
		return
	}

	in := usecase.CreateTeamInput{
		TeamName: req.TeamName,
		Members:  make([]usecase.TeamMemberDTO, 0, len(req.Members)),
	}
	for _, m := range req.Members {
		in.Members = append(in.Members, usecase.TeamMemberDTO{
			UserId:   m.UserId,
			UserName: m.Username,
			IsActive: m.IsActive,
		})
	}

	out, err := h.svc.CreateTeam(r.Context(), in)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	resp := createTeamResponseJSON{
		Team: mapCreateTeamOutputToJSON(out),
	}

	writeJSON(w, http.StatusCreated, resp)
}

// GET /team/get?team_name=...
func (h *HTTPHandler) handleGetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// for admins or users
	if _, ok := requireAnyAuth(w, r); !ok {
		return
	}

	teamName := r.URL.Query().Get("team_name")
	in := usecase.GetTeamInput{TeamName: teamName}

	out, err := h.svc.GetTeam(r.Context(), in)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	resp := getTeamResponseJSON{
		Team: mapGetTeamOutputToJSON(out),
	}

	writeJSON(w, http.StatusOK, resp)
}

func mapCreateTeamOutputToJSON(out *usecase.CreateTeamOutput) teamJSON {
	members := make([]teamMemberJSON, 0, len(out.Members))
	for _, m := range out.Members {
		members = append(members, teamMemberJSON{
			UserId:   m.UserId,
			Username: m.UserName,
			IsActive: m.IsActive,
		})
	}
	return teamJSON{
		TeamName: out.TeamName,
		Members:  members,
	}
}

func mapGetTeamOutputToJSON(out *usecase.GetTeamOutput) teamJSON {
	members := make([]teamMemberJSON, 0, len(out.Members))
	for _, m := range out.Members {
		members = append(members, teamMemberJSON{
			UserId:   m.UserId,
			Username: m.UserName,
			IsActive: m.IsActive,
		})
	}
	return teamJSON{
		TeamName: out.TeamName,
		Members:  members,
	}
}
