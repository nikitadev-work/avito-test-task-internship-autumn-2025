package usecase

import (
	"context"
	"testing"

	"pr-manager-service/internal/domain"
)

type mockUserRepo struct {
	getUserResp     *domain.User
	getUserErr      error
	getTeamNameResp string
	getTeamNameErr  error
}

func (m *mockUserRepo) GetUser(ctx context.Context, userId string) (*domain.User, error) {
	return m.getUserResp, m.getUserErr
}

func (m *mockUserRepo) SetIsActive(ctx context.Context, userId string, isActive bool) (*domain.User, string, error) {
	panic("not used in this test")
}

func (m *mockUserRepo) GetTeamName(ctx context.Context, userId string) (string, error) {
	return m.getTeamNameResp, m.getTeamNameErr
}

type mockPRRepo struct {
	createCalled bool
	createdPR    *domain.PullRequest

	getPRResp *domain.PullRequest
	getPRErr  error
}

func (m *mockPRRepo) CreatePullRequest(ctx context.Context, pr *domain.PullRequest) error {
	m.createCalled = true
	m.createdPR = pr
	return nil
}

func (m *mockPRRepo) GetPullRequest(ctx context.Context, prId string) (*domain.PullRequest, error) {
	return m.getPRResp, m.getPRErr
}

func (m *mockPRRepo) MergePullRequest(ctx context.Context, prId string) (*domain.PullRequest, error) {
	panic("not used in this test")
}

func (m *mockPRRepo) GetAllPrByUserId(ctx context.Context, userId string) ([]domain.PullRequest, error) {
	panic("not used in this test")
}

func (m *mockPRRepo) ReplaceReviewer(ctx context.Context, prId, oldUserId, newUserId string) error {
	panic("not used in this test")
}

func (m *mockPRRepo) GetActiveTeamMembers(ctx context.Context, teamName string) ([]domain.User, error) {
	return []domain.User{
		{UserId: "u1", UserName: "Alice", IsActive: true},
		{UserId: "u2", UserName: "Bob", IsActive: true},
		{UserId: "u3", UserName: "Charlie", IsActive: true},
	}, nil
}

type noopLogger struct{}

func (l *noopLogger) Debug(string, map[string]any) {}
func (l *noopLogger) Info(string, map[string]any)  {}
func (l *noopLogger) Warn(string, map[string]any)  {}
func (l *noopLogger) Error(string, map[string]any) {}

type dummyMetrics struct{}

func (m *dummyMetrics) IncTeamCreated()           {}
func (m *dummyMetrics) IncUserActivated()         {}
func (m *dummyMetrics) IncUserDeactivated()       {}
func (m *dummyMetrics) IncPullRequestCreated()    {}
func (m *dummyMetrics) IncPullRequestMerged()     {}
func (m *dummyMetrics) IncPullRequestReassigned() {}

func TestCreatePullRequest_AssignsReviewers(t *testing.T) {
	ctx := context.Background()

	userRepo := &mockUserRepo{
		getUserResp: &domain.User{
			UserId:   "u1",
			UserName: "Alice",
			IsActive: true,
		},
		getTeamNameResp: "payments",
	}
	prRepo := &mockPRRepo{}
	svc := &Service{
		teams:   nil,
		users:   userRepo,
		prs:     prRepo,
		logger:  &noopLogger{},
		metrics: &dummyMetrics{},
	}

	in := CreatePullRequestInput{
		PullRequestId:   "pr-1001",
		PullRequestName: "Add search",
		AuthorId:        "u1",
	}

	out, err := svc.CreatePullRequest(ctx, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !prRepo.createCalled {
		t.Fatalf("expected CreatePullRequest to be called on repository")
	}

	if len(out.PR.AssignedReviewers) != 2 {
		t.Fatalf("expected 2 assigned reviewers, got %d", len(out.PR.AssignedReviewers))
	}

	for _, r := range out.PR.AssignedReviewers {
		if r == "u1" {
			t.Fatalf("author should not be in assigned reviewers")
		}
	}
}

func TestReassignReviewer_MergedPR_ReturnsError(t *testing.T) {
	ctx := context.Background()

	prRepo := &mockPRRepo{
		getPRResp: &domain.PullRequest{
			PullRequestId:     "pr-1001",
			PullRequestName:   "Add search",
			AuthorId:          "u1",
			StatusId:          2,
			AssignedReviewers: []string{"u2", "u3"},
		},
	}
	userRepo := &mockUserRepo{}
	svc := &Service{
		users:   userRepo,
		prs:     prRepo,
		logger:  &noopLogger{},
		metrics: &dummyMetrics{},
	}

	in := ReassignReviewerInput{
		PullRequestId: "pr-1001",
		OldUserId:     "u2",
	}

	_, err := svc.ReassignReviewer(ctx, in)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != domain.ErrEditMergedPR {
		t.Fatalf("expected ErrEditMergedPR, got %v", err)
	}
}
