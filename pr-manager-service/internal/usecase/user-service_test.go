package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"pr-manager-service/internal/domain"
)

type userRepoMockForUserService struct {
	setUserResp   *domain.User
	setTeamName   string
	setErr        error
}

func (m *userRepoMockForUserService) GetUser(ctx context.Context, userId string) (*domain.User, error) {
	panic("not used in these tests")
}

func (m *userRepoMockForUserService) SetIsActive(ctx context.Context, userId string, isActive bool) (*domain.User, string, error) {
	return m.setUserResp, m.setTeamName, m.setErr
}

func (m *userRepoMockForUserService) GetTeamName(ctx context.Context, userId string) (string, error) {
	panic("not used in these tests")
}

type prRepoMockForUserService struct {
	getAllResp []domain.PullRequest
	getAllErr  error
}

func (m *prRepoMockForUserService) CreatePullRequest(ctx context.Context, pr *domain.PullRequest) error {
	panic("not used")
}

func (m *prRepoMockForUserService) GetPullRequest(ctx context.Context, prId string) (*domain.PullRequest, error) {
	panic("not used")
}

func (m *prRepoMockForUserService) MergePullRequest(ctx context.Context, prId string) (*domain.PullRequest, error) {
	panic("not used")
}

func (m *prRepoMockForUserService) GetAllPrByUserId(ctx context.Context, userId string) ([]domain.PullRequest, error) {
	return m.getAllResp, m.getAllErr
}

func (m *prRepoMockForUserService) ReplaceReviewer(ctx context.Context, prId, oldUserId, newUserId string) error {
	panic("not used")
}

func (m *prRepoMockForUserService) GetActiveTeamMembers(ctx context.Context, teamName string) ([]domain.User, error) {
	panic("not used")
}

type metricsMock struct {
	teamCreated           int
	userActivated         int
	userDeactivated       int
	prCreated             int
	prMerged              int
	prReassigned          int
}

func (m *metricsMock) IncTeamCreated()           { m.teamCreated++ }
func (m *metricsMock) IncUserActivated()         { m.userActivated++ }
func (m *metricsMock) IncUserDeactivated()       { m.userDeactivated++ }
func (m *metricsMock) IncPullRequestCreated()    { m.prCreated++ }
func (m *metricsMock) IncPullRequestMerged()     { m.prMerged++ }
func (m *metricsMock) IncPullRequestReassigned() { m.prReassigned++ }

func TestSetIsActive_TableDriven(t *testing.T) {
	ctx := context.Background()

	user := &domain.User{
		UserId:   "u1",
		UserName: "Alice",
		IsActive: true,
	}

	tests := []struct {
		name          string
		input         SetIsActiveInput
		repoUser      *domain.User
		repoTeam      string
		repoErr       error
		wantErr       error
		wantActivated int
		wantDeact     int
	}{
		{
			name:    "validation error - empty user id",
			input:   SetIsActiveInput{UserId: "", IsActive: true},
			wantErr: ErrUserIdRequired,
		},
		{
			name:    "user not found (sql.ErrNoRows)",
			input:   SetIsActiveInput{UserId: "u1", IsActive: true},
			repoErr: sql.ErrNoRows,
			wantErr: sql.ErrNoRows,
		},
		{
			name:    "repository error",
			input:   SetIsActiveInput{UserId: "u1", IsActive: true},
			repoErr: errors.New("db error"),
			wantErr: errors.New("db error"),
		},
		{
			name:          "activate user",
			input:         SetIsActiveInput{UserId: "u1", IsActive: true},
			repoUser:      user,
			repoTeam:      "payments",
			wantErr:       nil,
			wantActivated: 1,
			wantDeact:     0,
		},
		{
			name:          "deactivate user",
			input:         SetIsActiveInput{UserId: "u1", IsActive: false},
			repoUser:      user,
			repoTeam:      "payments",
			wantErr:       nil,
			wantActivated: 0,
			wantDeact:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &userRepoMockForUserService{
				setUserResp: tt.repoUser,
				setTeamName: tt.repoTeam,
				setErr:      tt.repoErr,
			}
			metrics := &metricsMock{}
			svc := &Service{
				teams:   nil,
				users:   userRepo,
				prs:     nil,
				logger:  &noopLogger{},
				metrics: metrics,
			}

			out, err := svc.SetIsActive(ctx, tt.input)

			if (err != nil) != (tt.wantErr != nil) {
				t.Fatalf("unexpected error: got=%v, wantErrExist=%v", err, tt.wantErr != nil)
			}
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected err %v, got %v", tt.wantErr, err)
				}
			}

			if tt.wantErr == nil && out == nil {
				t.Fatalf("expected non-nil output on success")
			}

			if metrics.userActivated != tt.wantActivated {
				t.Fatalf("expected userActivated=%d, got %d", tt.wantActivated, metrics.userActivated)
			}
			if metrics.userDeactivated != tt.wantDeact {
				t.Fatalf("expected userDeactivated=%d, got %d", tt.wantDeact, metrics.userDeactivated)
			}
		})
	}
}

func TestGetUserReviews_TableDriven(t *testing.T) {
	ctx := context.Background()

	prs := []domain.PullRequest{
		{PullRequestId: "pr-1", PullRequestName: "A", AuthorId: "u1", StatusId: 1},
		{PullRequestId: "pr-2", PullRequestName: "B", AuthorId: "u2", StatusId: 2},
	}

	tests := []struct {
		name     string
		input    GetUserReviewsInput
		repoPRs  []domain.PullRequest
		repoErr  error
		wantErr  error
		wantNil  bool
		wantLen  int
	}{
		{
			name:    "validation error - empty user id",
			input:   GetUserReviewsInput{UserId: ""},
			wantErr: ErrUserIdRequired,
			wantNil: true,
		},
		{
			name:    "repository error",
			input:   GetUserReviewsInput{UserId: "u1"},
			repoErr: errors.New("db error"),
			wantErr: errors.New("db error"),
			wantNil: true,
		},
		{
			name:    "ok",
			input:   GetUserReviewsInput{UserId: "u1"},
			repoPRs: prs,
			wantErr: nil,
			wantNil: false,
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prRepo := &prRepoMockForUserService{
				getAllResp: tt.repoPRs,
				getAllErr:  tt.repoErr,
			}
			svc := &Service{
				teams:   nil,
				users:   &userRepoMockForUserService{},
				prs:     prRepo,
				logger:  &noopLogger{},
				metrics: &dummyMetrics{},
			}

			out, err := svc.GetUserReviews(ctx, tt.input)

			if (err != nil) != (tt.wantErr != nil) {
				t.Fatalf("unexpected error: got=%v, wantErrExist=%v", err, tt.wantErr != nil)
			}
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected err %v, got %v", tt.wantErr, err)
				}
			}

			if tt.wantNil && out != nil {
				t.Fatalf("expected nil output, got %+v", out)
			}
			if !tt.wantNil {
				if out == nil {
					t.Fatalf("expected non-nil output")
				}
				if len(out.PullRequests) != tt.wantLen {
					t.Fatalf("expected %d PRs, got %d", tt.wantLen, len(out.PullRequests))
				}
			}
		})
	}
}
