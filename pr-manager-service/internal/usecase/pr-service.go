package usecase

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"

	"pr-manager-service/internal/domain"
)

// Pull requests

func (s *Service) CreatePullRequest(ctx context.Context, in CreatePullRequestInput) (*CreatePullRequestOutput, error) {
	if err := validateCreatePullRequestInput(in); err != nil {
		s.logger.Error("create pull request validation failed", map[string]any{
			"pull_request_id":   in.PullRequestId,
			"pull_request_name": in.PullRequestName,
			"author_id":         in.AuthorId,
			"error":             err.Error(),
		})
		return nil, err
	}

	s.logger.Info("create pull request started", map[string]any{
		"pull_request_id":   in.PullRequestId,
		"pull_request_name": in.PullRequestName,
		"author_id":         in.AuthorId,
	})

	// Check if the author exists
	author, err := s.users.GetUser(ctx, in.AuthorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("create pull request: author not found", map[string]any{
				"pull_request_id": in.PullRequestId,
				"author_id":       in.AuthorId,
				"error":           err.Error(),
			})
			return nil, err
		}

		s.logger.Error("create pull request: get author repository error", map[string]any{
			"pull_request_id": in.PullRequestId,
			"author_id":       in.AuthorId,
			"error":           err.Error(),
		})
		return nil, err
	}

	// Get author's team name
	teamName, err := s.users.GetTeamName(ctx, author.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("create pull request: author has no team", map[string]any{
				"pull_request_id": in.PullRequestId,
				"author_id":       in.AuthorId,
				"error":           err.Error(),
			})
			return nil, err
		}

		s.logger.Error("create pull request: get author team repository error", map[string]any{
			"pull_request_id": in.PullRequestId,
			"author_id":       in.AuthorId,
			"error":           err.Error(),
		})
		return nil, err
	}

	// Get active members from this team
	candidates, err := s.prs.GetActiveTeamMembers(ctx, teamName)
	if err != nil {
		s.logger.Error("create pull request: get active team members error", map[string]any{
			"pull_request_id": in.PullRequestId,
			"author_id":       in.AuthorId,
			"team_name":       teamName,
			"error":           err.Error(),
		})
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
		s.logger.Error("create pull request repository error", map[string]any{
			"pull_request_id": in.PullRequestId,
			"author_id":       in.AuthorId,
			"error":           err.Error(),
		})
		return nil, err
	}

	out := &CreatePullRequestOutput{
		PR: mapDomainPRToDTO(pr),
	}

	s.logger.Info("create pull request completed", map[string]any{
		"pull_request_id":    out.PR.PullRequestId,
		"author_id":          out.PR.AuthorId,
		"assigned_reviewers": out.PR.AssignedReviewers,
	})

	return out, nil
}

func (s *Service) MergePullRequest(ctx context.Context, in MergePullRequestInput) (*MergePullRequestOutput, error) {
	if err := validateMergePullRequestInput(in); err != nil {
		s.logger.Error("merge pull request validation failed", map[string]any{
			"pull_request_id": in.PullRequestId,
			"error":           err.Error(),
		})
		return nil, err
	}

	s.logger.Info("merge pull request started", map[string]any{
		"pull_request_id": in.PullRequestId,
	})

	pr, err := s.prs.MergePullRequest(ctx, in.PullRequestId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("merge pull request: pr not found", map[string]any{
				"pull_request_id": in.PullRequestId,
				"error":           err.Error(),
			})
			return nil, err
		}

		s.logger.Error("merge pull request repository error", map[string]any{
			"pull_request_id": in.PullRequestId,
			"error":           err.Error(),
		})
		return nil, err
	}

	out := &MergePullRequestOutput{
		PR: mapDomainPRToDTO(pr),
	}

	s.logger.Info("merge pull request completed", map[string]any{
		"pull_request_id": out.PR.PullRequestId,
		"status":          out.PR.Status,
	})

	return out, nil
}

func (s *Service) ReassignReviewer(ctx context.Context, in ReassignReviewerInput) (*ReassignReviewerOutput, error) {
	if err := validateReassignReviewerInput(in); err != nil {
		s.logger.Error("reassign reviewer validation failed", map[string]any{
			"pull_request_id": in.PullRequestId,
			"old_user_id":     in.OldUserId,
			"error":           err.Error(),
		})
		return nil, err
	}

	s.logger.Info("reassign reviewer started", map[string]any{
		"pull_request_id": in.PullRequestId,
		"old_user_id":     in.OldUserId,
	})

	// Get PR
	pr, err := s.prs.GetPullRequest(ctx, in.PullRequestId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("reassign reviewer: pr not found", map[string]any{
				"pull_request_id": in.PullRequestId,
				"old_user_id":     in.OldUserId,
				"error":           err.Error(),
			})
			return nil, err
		}

		s.logger.Error("reassign reviewer: get pr repository error", map[string]any{
			"pull_request_id": in.PullRequestId,
			"old_user_id":     in.OldUserId,
			"error":           err.Error(),
		})
		return nil, err
	}

	// Check if the PR is already merged
	if statusString(pr.StatusId) == "MERGED" {
		s.logger.Warn("reassign reviewer: pr already merged", map[string]any{
			"pull_request_id": in.PullRequestId,
			"old_user_id":     in.OldUserId,
		})
		return nil, domain.ErrEditMergedPR
	}

	// Check if the old reviewer is correct
	foundOld := false
	for _, r := range pr.AssignedReviewers {
		if r == in.OldUserId {
			foundOld = true
			break
		}
	}
	if !foundOld {
		s.logger.Warn("reassign reviewer: old reviewer is not assigned", map[string]any{
			"pull_request_id": in.PullRequestId,
			"old_user_id":     in.OldUserId,
		})
		return nil, ErrReviewerNotAssigned
	}

	// Get team of the old reviewer
	teamName, err := s.users.GetTeamName(ctx, in.OldUserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("reassign reviewer: reviewer has no team", map[string]any{
				"pull_request_id": in.PullRequestId,
				"old_user_id":     in.OldUserId,
				"error":           err.Error(),
			})
			return nil, err
		}

		s.logger.Error("reassign reviewer: get team name repository error", map[string]any{
			"pull_request_id": in.PullRequestId,
			"old_user_id":     in.OldUserId,
			"error":           err.Error(),
		})
		return nil, err
	}

	// Get only active members
	candidates, err := s.prs.GetActiveTeamMembers(ctx, teamName)
	if err != nil {
		s.logger.Error("reassign reviewer: get active team members error", map[string]any{
			"pull_request_id": in.PullRequestId,
			"old_user_id":     in.OldUserId,
			"team_name":       teamName,
			"error":           err.Error(),
		})
		return nil, err
	}

	// Filter the candidates: exclude old reviewer and already assigned members
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
		s.logger.Warn("reassign reviewer: no available candidates", map[string]any{
			"pull_request_id": in.PullRequestId,
			"old_user_id":     in.OldUserId,
			"team_name":       teamName,
		})
		return nil, domain.ErrNoAvailableCandidates
	}

	// Choose new reviewer
	newReviewerId := candidateIds[rand.Intn(len(candidateIds))]

	// Assign new reviewer
	err = s.prs.ReplaceReviewer(ctx, in.PullRequestId, in.OldUserId, newReviewerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("reassign reviewer: old reviewer not found in db for this pr", map[string]any{
				"pull_request_id": in.PullRequestId,
				"old_user_id":     in.OldUserId,
				"new_user_id":     newReviewerId,
				"error":           err.Error(),
			})
			return nil, err
		}

		s.logger.Error("reassign reviewer: replace reviewer repository error", map[string]any{
			"pull_request_id": in.PullRequestId,
			"old_user_id":     in.OldUserId,
			"new_user_id":     newReviewerId,
			"error":           err.Error(),
		})
		return nil, err
	}

	// Get the edited PR
	updatedPr, err := s.prs.GetPullRequest(ctx, in.PullRequestId)
	if err != nil {
		s.logger.Error("reassign reviewer: get updated pr repository error", map[string]any{
			"pull_request_id": in.PullRequestId,
			"old_user_id":     in.OldUserId,
			"new_user_id":     newReviewerId,
			"error":           err.Error(),
		})
		return nil, err
	}

	out := &ReassignReviewerOutput{
		PR:         mapDomainPRToDTO(updatedPr),
		ReplacedBy: newReviewerId,
	}

	s.logger.Info("reassign reviewer completed", map[string]any{
		"pull_request_id": out.PR.PullRequestId,
		"old_user_id":     in.OldUserId,
		"new_user_id":     out.ReplacedBy,
	})

	return out, nil
}
