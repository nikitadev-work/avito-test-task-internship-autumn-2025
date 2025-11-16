package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"pr-manager-service/internal/domain"
)

type mockTeamRepo struct {
	createErr    error
	createCalled bool

	getTeamRespTeam  *domain.Team
	getTeamRespUsers []domain.User
	getTeamErr       error
}

func (m *mockTeamRepo) CreateTeam(ctx context.Context, teamName string, members []domain.User) error {
	m.createCalled = true
	return m.createErr
}

func (m *mockTeamRepo) GetTeam(ctx context.Context, teamName string) (*domain.Team, []domain.User, error) {
	return m.getTeamRespTeam, m.getTeamRespUsers, m.getTeamErr
}

func TestCreateTeam_TableDriven(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     CreateTeamInput
		repoErr   error
		wantErr   error
		wantCalls bool
	}{
		{
			name: "ok",
			input: CreateTeamInput{
				TeamName: "backend",
				Members: []TeamMemberDTO{
					{UserId: "u1", UserName: "Alice", IsActive: true},
				},
			},
			repoErr:   nil,
			wantErr:   nil,
			wantCalls: true,
		},
		{
			name: "validation error - empty team name",
			input: CreateTeamInput{
				TeamName: "",
			},
			repoErr:   nil,
			wantErr:   ErrTeamNameRequired,
			wantCalls: false,
		},
		{
			name: "repository error",
			input: CreateTeamInput{
				TeamName: "backend",
			},
			repoErr:   errors.New("db error"),
			wantErr:   errors.New("db error"),
			wantCalls: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teamRepo := &mockTeamRepo{createErr: tt.repoErr}
			svc := &Service{
				teams:   teamRepo,
				users:   nil,
				prs:     nil,
				logger:  &noopLogger{},
				metrics: &dummyMetrics{},
			}

			out, err := svc.CreateTeam(ctx, tt.input)

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
			if teamRepo.createCalled != tt.wantCalls {
				t.Fatalf("expected createCalled=%v, got %v", tt.wantCalls, teamRepo.createCalled)
			}
		})
	}
}

func TestGetTeam_TableDriven(t *testing.T) {
	ctx := context.Background()

	members := []domain.User{
		{UserId: "u1", UserName: "Alice", IsActive: true},
		{UserId: "u2", UserName: "Bob", IsActive: false},
	}

	tests := []struct {
		name     string
		input    GetTeamInput
		repoErr  error
		repoRows []domain.User
		wantErr  error
		wantNil  bool
	}{
		{
			name: "ok",
			input: GetTeamInput{
				TeamName: "backend",
			},
			repoRows: members,
			repoErr:  nil,
			wantErr:  nil,
			wantNil:  false,
		},
		{
			name: "validation error - empty team name",
			input: GetTeamInput{
				TeamName: "",
			},
			repoErr: nil,
			wantErr: ErrTeamNameRequired,
			wantNil: true,
		},
		{
			name: "team not found",
			input: GetTeamInput{
				TeamName: "backend",
			},
			repoErr: sql.ErrNoRows,
			wantErr: sql.ErrNoRows,
			wantNil: true,
		},
		{
			name: "repository error",
			input: GetTeamInput{
				TeamName: "backend",
			},
			repoErr: errors.New("db error"),
			wantErr: errors.New("db error"),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teamRepo := &mockTeamRepo{
				getTeamRespUsers: tt.repoRows,
				getTeamErr:       tt.repoErr,
			}
			svc := &Service{
				teams:   teamRepo,
				users:   nil,
				prs:     nil,
				logger:  &noopLogger{},
				metrics: &dummyMetrics{},
			}

			out, err := svc.GetTeam(ctx, tt.input)

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
			if !tt.wantNil && out == nil {
				t.Fatalf("expected non-nil output")
			}
		})
	}
}
