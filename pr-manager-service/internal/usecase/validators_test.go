package usecase

import (
	"errors"
	"testing"
)

func TestValidateCreateTeamInput(t *testing.T) {
	tests := []struct {
		name    string
		in      CreateTeamInput
		wantErr error
	}{
		{
			name:    "ok",
			in:      CreateTeamInput{TeamName: "backend"},
			wantErr: nil,
		},
		{
			name:    "empty team name",
			in:      CreateTeamInput{TeamName: ""},
			wantErr: ErrTeamNameRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateTeamInput(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateGetTeamInput(t *testing.T) {
	tests := []struct {
		name    string
		in      GetTeamInput
		wantErr error
	}{
		{
			name:    "ok",
			in:      GetTeamInput{TeamName: "backend"},
			wantErr: nil,
		},
		{
			name:    "empty team name",
			in:      GetTeamInput{TeamName: ""},
			wantErr: ErrTeamNameRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGetTeamInput(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateSetIsActiveInput(t *testing.T) {
	tests := []struct {
		name    string
		in      SetIsActiveInput
		wantErr error
	}{
		{
			name:    "ok",
			in:      SetIsActiveInput{UserId: "u1", IsActive: true},
			wantErr: nil,
		},
		{
			name:    "empty user id",
			in:      SetIsActiveInput{UserId: "", IsActive: true},
			wantErr: ErrUserIdRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSetIsActiveInput(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateGetUserReviewsInput(t *testing.T) {
	tests := []struct {
		name    string
		in      GetUserReviewsInput
		wantErr error
	}{
		{
			name:    "ok",
			in:      GetUserReviewsInput{UserId: "u1"},
			wantErr: nil,
		},
		{
			name:    "empty user id",
			in:      GetUserReviewsInput{UserId: ""},
			wantErr: ErrUserIdRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGetUserReviewsInput(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateCreatePullRequestInput(t *testing.T) {
	tests := []struct {
		name    string
		in      CreatePullRequestInput
		wantErr error
	}{
		{
			name: "ok",
			in: CreatePullRequestInput{
				PullRequestId:   "pr-1",
				PullRequestName: "Add search",
				AuthorId:        "u1",
			},
			wantErr: nil,
		},
		{
			name: "empty id",
			in: CreatePullRequestInput{
				PullRequestId:   "",
				PullRequestName: "Add search",
				AuthorId:        "u1",
			},
			wantErr: ErrPullRequestIdRequired,
		},
		{
			name: "empty name",
			in: CreatePullRequestInput{
				PullRequestId:   "pr-1",
				PullRequestName: "",
				AuthorId:        "u1",
			},
			wantErr: ErrPullRequestNameRequired,
		},
		{
			name: "empty author",
			in: CreatePullRequestInput{
				PullRequestId:   "pr-1",
				PullRequestName: "Add search",
				AuthorId:        "",
			},
			wantErr: ErrAuthorIdRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreatePullRequestInput(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateMergePullRequestInput(t *testing.T) {
	tests := []struct {
		name    string
		in      MergePullRequestInput
		wantErr error
	}{
		{
			name:    "ok",
			in:      MergePullRequestInput{PullRequestId: "pr-1"},
			wantErr: nil,
		},
		{
			name:    "empty id",
			in:      MergePullRequestInput{},
			wantErr: ErrPullRequestIdRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMergePullRequestInput(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateReassignReviewerInput(t *testing.T) {
	tests := []struct {
		name    string
		in      ReassignReviewerInput
		wantErr error
	}{
		{
			name:    "ok",
			in:      ReassignReviewerInput{PullRequestId: "pr-1", OldUserId: "u2"},
			wantErr: nil,
		},
		{
			name:    "empty pr id",
			in:      ReassignReviewerInput{PullRequestId: "", OldUserId: "u2"},
			wantErr: ErrPullRequestIdRequired,
		},
		{
			name:    "empty old user id",
			in:      ReassignReviewerInput{PullRequestId: "pr-1", OldUserId: ""},
			wantErr: ErrOldUserIdRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReassignReviewerInput(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}
		})
	}
}
