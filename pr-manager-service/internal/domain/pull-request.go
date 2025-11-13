package domain

type PullRequest struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	StatusId          int
	AssignedReviewers []string
}

