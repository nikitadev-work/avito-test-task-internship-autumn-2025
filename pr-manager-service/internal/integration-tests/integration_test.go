package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

// PRM_BASE_URL - basic URL
func baseURL() string {
	if v := os.Getenv("PRM_BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

func newAdminRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("failed to encode body: %v", err)
		}
	}

	req, err := http.NewRequest(method, baseURL()+path, &buf)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin:integration-admin")

	return req
}

type getTeamResponse struct {
	Team struct {
		TeamName string `json:"team_name"`
		Members  []struct {
			UserID   string `json:"user_id"`
			UserName string `json:"username"`
			IsActive bool   `json:"is_active"`
		} `json:"members"`
	} `json:"team"`
}

type getUserReviewsResponse struct {
	UserID       string `json:"user_id"`
	PullRequests []struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		Status          string `json:"status"`
	} `json:"pull_requests"`
}

// Create team and read it after creation
func TestCreateTeamAndGetTeam_Integration(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	teamName := fmt.Sprintf("integration-team-%d", time.Now().UnixNano())

	// 1) POST /team/add
	createBody := map[string]any{
		"team_name": teamName,
		"members": []map[string]any{
			{
				"user_id":   "u_int_1",
				"username":  "Integration User 1",
				"is_active": true,
			},
			{
				"user_id":   "u_int_2",
				"username":  "Integration User 2",
				"is_active": true,
			},
		},
	}

	req := newAdminRequest(t, http.MethodPost, "/team/add", createBody)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to call /team/add: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	// 2) GET /team/get?team_name=...
	getReq := newAdminRequest(t, http.MethodGet, "/team/get?team_name="+teamName, nil)
	resp2, err := client.Do(getReq)
	if err != nil {
		t.Fatalf("failed to call /team/get: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp2.StatusCode)
	}

	var teamResp getTeamResponse
	if err := json.NewDecoder(resp2.Body).Decode(&teamResp); err != nil {
		t.Fatalf("failed to decode /team/get response: %v", err)
	}

	if teamResp.Team.TeamName != teamName {
		t.Fatalf("expected team_name %q, got %q", teamName, teamResp.Team.TeamName)
	}
	if len(teamResp.Team.Members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(teamResp.Team.Members))
	}

}

// Check creation of PR and getting this PR
func TestCreatePullRequestAndGetUserReviews_Integration(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	teamName := fmt.Sprintf("integration-team-pr-%d", time.Now().UnixNano())

	createTeamBody := map[string]any{
		"team_name": teamName,
		"members": []map[string]any{
			{
				"user_id":   "u_pr_author",
				"username":  "Author",
				"is_active": true,
			},
			{
				"user_id":   "u_pr_reviewer1",
				"username":  "Reviewer1",
				"is_active": true,
			},
			{
				"user_id":   "u_pr_reviewer2",
				"username":  "Reviewer2",
				"is_active": true,
			},
		},
	}

	reqTeam := newAdminRequest(t, http.MethodPost, "/team/add", createTeamBody)
	respTeam, err := client.Do(reqTeam)
	if err != nil {
		t.Fatalf("failed to call /team/add: %v", err)
	}
	respTeam.Body.Close()
	if respTeam.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201 from /team/add, got %d", respTeam.StatusCode)
	}

	prID := fmt.Sprintf("pr-int-%d", time.Now().UnixNano())
	createPRBody := map[string]any{
		"pull_request_id":   prID,
		"pull_request_name": "Integration PR",
		"author_id":         "u_pr_author",
	}

	reqPR := newAdminRequest(t, http.MethodPost, "/pullRequest/create", createPRBody)
	respPR, err := client.Do(reqPR)
	if err != nil {
		t.Fatalf("failed to call /pullRequest/create: %v", err)
	}
	defer respPR.Body.Close()

	if respPR.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201 from /pullRequest/create, got %d", respPR.StatusCode)
	}

	// Check if the endpoint works correctly
	getReq, err := http.NewRequest(http.MethodGet, baseURL()+"/users/getReview?user_id=u_pr_reviewer1", nil)
	if err != nil {
		t.Fatalf("failed to create request for /users/getReview: %v", err)
	}
	getReq.Header.Set("Authorization", "Bearer user:u_pr_reviewer1")

	respGet, err := client.Do(getReq)
	if err != nil {
		t.Fatalf("failed to call /users/getReview: %v", err)
	}
	defer respGet.Body.Close()

	if respGet.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 from /users/getReview, got %d", respGet.StatusCode)
	}

	var reviews getUserReviewsResponse
	if err := json.NewDecoder(respGet.Body).Decode(&reviews); err != nil {
		t.Fatalf("failed to decode /users/getReview response: %v", err)
	}

	if reviews.UserID != "u_pr_reviewer1" {
		t.Fatalf("expected user_id u_pr_reviewer1, got %s", reviews.UserID)
	}
}
