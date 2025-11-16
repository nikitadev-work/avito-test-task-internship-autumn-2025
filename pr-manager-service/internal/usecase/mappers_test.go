package usecase

import (
	"testing"

	"pr-manager-service/internal/domain"
)

func TestMapTeamMembersDTOToDomain(t *testing.T) {
	in := []TeamMemberDTO{
		{UserId: "u1", UserName: "Alice", IsActive: true},
		{UserId: "u2", UserName: "Bob", IsActive: false},
	}

	out := mapTeamMembersDTOToDomain(in)

	if len(out) != len(in) {
		t.Fatalf("expected %d users, got %d", len(in), len(out))
	}

	for i, u := range out {
		if u.UserId != in[i].UserId || u.UserName != in[i].UserName || u.IsActive != in[i].IsActive {
			t.Fatalf("mismatch at %d: got %+v, want %+v", i, u, in[i])
		}
	}
}

func TestMapDomainUsersToTeamMembersDTO(t *testing.T) {
	in := []domain.User{
		{UserId: "u1", UserName: "Alice", IsActive: true},
		{UserId: "u2", UserName: "Bob", IsActive: false},
	}

	out := mapDomainUsersToTeamMembersDTO(in)

	if len(out) != len(in) {
		t.Fatalf("expected %d members, got %d", len(in), len(out))
	}

	for i, m := range out {
		if m.UserId != in[i].UserId || m.UserName != in[i].UserName || m.IsActive != in[i].IsActive {
			t.Fatalf("mismatch at %d: got %+v, want %+v", i, m, in[i])
		}
	}
}

func TestMapDomainUserToSetIsActiveOutput(t *testing.T) {
	user := &domain.User{
		UserId:   "u1",
		UserName: "Alice",
		IsActive: true,
	}
	teamName := "payments"

	out := mapDomainUserToSetIsActiveOutput(user, teamName)

	if out.UserId != user.UserId || out.UserName != user.UserName || out.TeamName != teamName || out.IsActive != user.IsActive {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestMapDomainPRsToGetUserReviewsOutput(t *testing.T) {
	prs := []domain.PullRequest{
		{PullRequestId: "pr-1", PullRequestName: "A", AuthorId: "u1", StatusId: 1},
		{PullRequestId: "pr-2", PullRequestName: "B", AuthorId: "u2", StatusId: 2},
	}

	out := mapDomainPRsToGetUserReviewsOutput("uX", prs)

	if out.UserId != "uX" {
		t.Fatalf("expected userId uX, got %s", out.UserId)
	}
	if len(out.PullRequests) != 2 {
		t.Fatalf("expected 2 PRs, got %d", len(out.PullRequests))
	}
	if out.PullRequests[0].Status != "OPEN" || out.PullRequests[1].Status != "MERGED" {
		t.Fatalf("unexpected statuses: %+v", out.PullRequests)
	}
}

func TestMapCreatePRInputToDomain(t *testing.T) {
	in := CreatePullRequestInput{
		PullRequestId:   "pr-1",
		PullRequestName: "Add search",
		AuthorId:        "u1",
	}
	assigned := []string{"u2", "u3"}

	pr := mapCreatePRInputToDomain(in, assigned)

	if pr.PullRequestId != in.PullRequestId ||
		pr.PullRequestName != in.PullRequestName ||
		pr.AuthorId != in.AuthorId {
		t.Fatalf("unexpected mapped fields: %+v", pr)
	}
	if pr.StatusId != 1 {
		t.Fatalf("expected StatusId=1, got %d", pr.StatusId)
	}
	if len(pr.AssignedReviewers) != len(assigned) {
		t.Fatalf("expected %d assigned reviewers, got %d", len(assigned), len(pr.AssignedReviewers))
	}
}

func TestMapDomainPRToDTO(t *testing.T) {
	pr := &domain.PullRequest{
		PullRequestId:     "pr-1",
		PullRequestName:   "Add search",
		AuthorId:          "u1",
		StatusId:          2,
		AssignedReviewers: []string{"u2"},
	}

	dto := mapDomainPRToDTO(pr)

	if dto.PullRequestId != pr.PullRequestId ||
		dto.PullRequestName != pr.PullRequestName ||
		dto.AuthorId != pr.AuthorId {
		t.Fatalf("unexpected mapped fields: %+v", dto)
	}
	if dto.Status != "MERGED" {
		t.Fatalf("expected status MERGED, got %s", dto.Status)
	}
	if len(dto.AssignedReviewers) != 1 || dto.AssignedReviewers[0] != "u2" {
		t.Fatalf("unexpected reviewers: %+v", dto.AssignedReviewers)
	}
}

func TestStatusString(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want string
	}{
		{"open", 1, "OPEN"},
		{"merged", 2, "MERGED"},
		{"default", 42, "OPEN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := statusString(tt.in)
			if got != tt.want {
				t.Fatalf("statusString(%d) = %s, want %s", tt.in, got, tt.want)
			}
		})
	}
}
