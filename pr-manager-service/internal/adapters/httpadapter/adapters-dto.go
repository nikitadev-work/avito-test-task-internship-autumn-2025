package httpadapter

type teamMemberJSON struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type teamJSON struct {
	TeamName string           `json:"team_name"`
	Members  []teamMemberJSON `json:"members"`
}

type setIsActiveRequestJSON struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type userJSON struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type pullRequestCreateJSON struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

type pullRequestIdJSON struct {
	PullRequestId string `json:"pull_request_id"`
}

type reassignRequestJSON struct {
	PullRequestId string `json:"pull_request_id"`
	OldUserId     string `json:"old_user_id"`
}

type pullRequestJSON struct {
	PullRequestId     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorId          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

type pullRequestShortJSON struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
	Status          string `json:"status"`
}

// ErrorResponse по OpenAPI.

type errorBodyJSON struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponseJSON struct {
	Error errorBodyJSON `json:"error"`
}

// Responses

type createTeamResponseJSON struct {
	Team teamJSON `json:"team"`
}

type getTeamResponseJSON struct {
	Team teamJSON `json:"team"`
}

type setIsActiveResponseJSON struct {
	User userJSON `json:"user"`
}

type userReviewsResponseJSON struct {
	UserId       string                 `json:"user_id"`
	PullRequests []pullRequestShortJSON `json:"pull_requests"`
}

type pullRequestResponseJSON struct {
	PR pullRequestJSON `json:"pr"`
}

type reassignResponseJSON struct {
	PR         pullRequestJSON `json:"pr"`
	ReplacedBy string          `json:"replaced_by"`
}

type healthResponseJSON struct {
	Status string `json:"status"`
}
