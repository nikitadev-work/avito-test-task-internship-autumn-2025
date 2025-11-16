package domain

// PullRequest represents a domain pull request entity
type PullRequest struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	StatusId          int
	AssignedReviewers []string
}
