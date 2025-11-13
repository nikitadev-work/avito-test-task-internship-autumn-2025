package usecase

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"

	"pr-manager-service/internal/domain"
)

type Service struct {
	teams TeamRepositoryInterface
	users UserRepositoryInterface
	prs   PullRequestRepositoryInterface
}

func NewService(
	teams TeamRepositoryInterface,
	users UserRepositoryInterface,
	prs PullRequestRepositoryInterface,
) *Service {
	return &Service{
		teams: teams,
		users: users,
		prs:   prs,
	}
}

// Teams

func (s *Service) CreateTeam(ctx context.Context, in CreateTeamInput) (*CreateTeamOutput, error) {
	if err := validateCreateTeamInput(in); err != nil {
		return nil, err
	}

	members := mapTeamMembersDTOToDomain(in.Members)

	err := s.teams.CreateTeam(ctx, in.TeamName, members)
	if err != nil {
		return nil, err
	}

	out := &CreateTeamOutput{
		TeamName: in.TeamName,
		Members:  in.Members,
	}

	return out, nil
}

func (s *Service) GetTeam(ctx context.Context, in GetTeamInput) (*GetTeamOutput, error) {
	if err := validateGetTeamInput(in); err != nil {
		return nil, err
	}

	_, members, err := s.teams.GetTeam(ctx, in.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	outMembers := mapDomainUsersToTeamMembersDTO(members)

	out := &GetTeamOutput{
		TeamName: in.TeamName,
		Members:  outMembers,
	}

	return out, nil
}

// Users

func (s *Service) SetIsActive(ctx context.Context, in SetIsActiveInput) (*SetIsActiveOutput, error) {
	if err := validateSetIsActiveInput(in); err != nil {
		return nil, err
	}

	user, teamName, err := s.users.SetIsActive(ctx, in.UserId, in.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	out := mapDomainUserToSetIsActiveOutput(user, teamName)

	return out, nil
}

func (s *Service) GetUserReviews(ctx context.Context, in GetUserReviewsInput) (*GetUserReviewsOutput, error) {
	if err := validateGetUserReviewsInput(in); err != nil {
		return nil, err
	}

	prs, err := s.prs.GetAllPrByUserId(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	out := mapDomainPRsToGetUserReviewsOutput(in.UserId, prs)

	return out, nil
}

// Pull requests

func (s *Service) CreatePullRequest(ctx context.Context, in CreatePullRequestInput) (*CreatePullRequestOutput, error) {
	if err := validateCreatePullRequestInput(in); err != nil {
		return nil, err
	}

	// Check if the author exists
	author, err := s.users.GetUser(ctx, in.AuthorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	// Get author's team name
	teamName, err := s.users.GetTeamName(ctx, author.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	// Get active members from this team
	candidates, err := s.prs.GetActiveTeamMembers(ctx, teamName)
	if err != nil {
		return nil, err
	}

	// Exclude the author from candidates
	filtered := make([]domain.User, 0, len(candidates))
	for _, u := range candidates {
		if u.UserId == author.UserId {
			continue
		}
		filtered = append(filtered, u)
	}

	// Assign up to 2 reviewers
	assigned := make([]string, 0, 2)
	for i := 0; i < len(filtered) && i < 2; i++ {
		assigned = append(assigned, filtered[i].UserId)
	}

	pr := mapCreatePRInputToDomain(in, assigned)

	err = s.prs.CreatePullRequest(ctx, pr)
	if err != nil {
		return nil, err
	}

	out := &CreatePullRequestOutput{
		PR: mapDomainPRToDTO(pr),
	}

	return out, nil
}

func (s *Service) MergePullRequest(ctx context.Context, in MergePullRequestInput) (*MergePullRequestOutput, error) {
	if err := validateMergePullRequestInput(in); err != nil {
		return nil, err
	}

	pr, err := s.prs.MergePullRequest(ctx, in.PullRequestId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	out := &MergePullRequestOutput{
		PR: mapDomainPRToDTO(pr),
	}

	return out, nil
}

func (s *Service) ReassignReviewer(ctx context.Context, in ReassignReviewerInput) (*ReassignReviewerOutput, error) {
	if err := validateReassignReviewerInput(in); err != nil {
		return nil, err
	}

	// Get PR
	pr, err := s.prs.GetPullRequest(ctx, in.PullRequestId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	// Check if the PR is already merged
	if statusString(pr.StatusId) == "MERGED" {
		return nil, domain.ErrEditMergedPR
	}

	// Check if the old reviewer are correct
	foundOld := false
	for _, r := range pr.AssignedReviewers {
		if r == in.OldUserId {
			foundOld = true
			break
		}
	}
	if !foundOld {
		return nil, errors.New("reviewer is not assigned to this PR")
	}

	// Get team of the old reviewer
	teamName, err := s.users.GetTeamName(ctx, in.OldUserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	// Get only active members
	candidates, err := s.prs.GetActiveTeamMembers(ctx, teamName)
	if err != nil {
		return nil, err
	}

	// Filter the candidates: exclude author and already assigned members
	candidateIds := make([]string, 0)
	for _, u := range candidates {
		if u.UserId == in.OldUserId {
			continue
		}
		alreadyAssigned := false
		for _, r := range pr.AssignedReviewers {
			if r == u.UserId {
				alreadyAssigned = true
				break
			}
		}
		if alreadyAssigned {
			continue
		}
		candidateIds = append(candidateIds, u.UserId)
	}

	if len(candidateIds) == 0 {
		return nil, domain.ErrNoAvailableCandidates
	}

	// Choose new reviewer
	newReviewerId := candidateIds[rand.Intn(len(candidateIds))]

	// Assign new reviewer
	err = s.prs.ReplaceReviewer(ctx, in.PullRequestId, in.OldUserId, newReviewerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	// Get the edited PR
	updatedPr, err := s.prs.GetPullRequest(ctx, in.PullRequestId)
	if err != nil {
		return nil, err
	}

	out := &ReassignReviewerOutput{
		PR:         mapDomainPRToDTO(updatedPr),
		ReplacedBy: newReviewerId,
	}

	return out, nil
}
