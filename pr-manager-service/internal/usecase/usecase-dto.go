package usecase

// Teams

type TeamMemberDTO struct {
	UserId   string
	UserName string
	IsActive bool
}

type CreateTeamInput struct {
	TeamName string
	Members  []TeamMemberDTO
}

type CreateTeamOutput struct {
	TeamName string
	Members  []TeamMemberDTO
}

type GetTeamInput struct {
	TeamName string
}

type GetTeamOutput struct {
	TeamName string
	Members  []TeamMemberDTO
}

// Users

type SetIsActiveInput struct {
	UserId   string
	IsActive bool
}

type SetIsActiveOutput struct {
	UserId   string
	UserName string
	TeamName string
	IsActive bool
}

type GetUserReviewsInput struct {
	UserId string
}

type PullRequestShortDTO struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
	Status          string
}

type GetUserReviewsOutput struct {
	UserId       string
	PullRequests []PullRequestShortDTO
}

// Pull requests

type CreatePullRequestInput struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
}

type PullRequestDTO struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	Status            string
	AssignedReviewers []string
}

type CreatePullRequestOutput struct {
	PR PullRequestDTO
}

type MergePullRequestInput struct {
	PullRequestId string
}

type MergePullRequestOutput struct {
	PR PullRequestDTO
}

type ReassignReviewerInput struct {
	PullRequestId string
	OldUserId     string
}

type ReassignReviewerOutput struct {
	PR         PullRequestDTO
	ReplacedBy string
}
