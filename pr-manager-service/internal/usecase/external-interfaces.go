package usecase

import (
	"context"

	"pr-manager-service/internal/domain"
)

type TeamRepositoryInterface interface {
	CreateTeam(ctx context.Context, teamName string, members []domain.User) error
	GetTeam(ctx context.Context, teamName string) (*domain.Team, []domain.User, error)
}

type UserRepositoryInterface interface {
	GetUser(ctx context.Context, userId string) (*domain.User, error)
	SetIsActive(ctx context.Context, userId string, isActive bool) (*domain.User, string, error)
	GetTeamName(ctx context.Context, userId string) (string, error)
}

type PullRequestRepositoryInterface interface {
	CreatePullRequest(ctx context.Context, pr *domain.PullRequest) error
	GetPullRequest(ctx context.Context, prId string) (*domain.PullRequest, error)
	MergePullRequest(ctx context.Context, prId string) (*domain.PullRequest, error)
	GetAllPrByUserId(ctx context.Context, userId string) ([]domain.PullRequest, error)
	ReplaceReviewer(ctx context.Context, prId, oldUserId, newUserId string) error
	GetActiveTeamMembers(ctx context.Context, teamName string) ([]domain.User, error)
}

type LoggerInterface interface {
	Debug(msg string, params map[string]any)
	Info(msg string, params map[string]any)
	Warn(msg string, params map[string]any)
	Error(msg string, params map[string]any)
}

type MetricsInterface interface {
	IncTeamCreated()
	IncUserActivated()
	IncUserDeactivated()
	IncPullRequestCreated()
	IncPullRequestMerged()
	IncPullRequestReassigned()
}
