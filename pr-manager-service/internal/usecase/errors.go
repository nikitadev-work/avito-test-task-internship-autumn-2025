package usecase

import "errors"

var (
	ErrTeamNameRequired         = errors.New("team_name is required")
	ErrUserIdRequired           = errors.New("user_id is required")
	ErrPullRequestIdRequired    = errors.New("pull_request_id is required")
	ErrPullRequestNameRequired  = errors.New("pull_request_name is required")
	ErrAuthorIdRequired         = errors.New("author_id is required")
	ErrOldUserIdRequired        = errors.New("old_user_id is required")
	ErrTeamAlreadyExists        = errors.New("team already exists")
	ErrPullRequestAlreadyExists = errors.New("pull request already exists")
	ErrReviewerNotAssigned      = errors.New("reviewer not assigned for this pr")
	ErrNoCandidateInTeam        = errors.New("no review candidates in this team")
	ErrNotFound                 = errors.New("resource not found")
)
