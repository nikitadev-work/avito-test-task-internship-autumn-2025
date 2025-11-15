package usecase

import (
	"context"
	"database/sql"
	"errors"
)

// Users

func (s *Service) SetIsActive(ctx context.Context, in SetIsActiveInput) (*SetIsActiveOutput, error) {
	if err := validateSetIsActiveInput(in); err != nil {
		s.logger.Error("set is_active validation failed", map[string]any{
			"user_id": in.UserId,
			"error":   err.Error(),
		})
		return nil, err
	}

	s.logger.Info("set is_active started", map[string]any{
		"user_id":   in.UserId,
		"is_active": in.IsActive,
	})

	user, teamName, err := s.users.SetIsActive(ctx, in.UserId, in.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("set is_active: user not found", map[string]any{
				"user_id": in.UserId,
				"error":   err.Error(),
			})
			return nil, err
		}

		s.logger.Error("set is_active repository error", map[string]any{
			"user_id": in.UserId,
			"error":   err.Error(),
		})
		return nil, err
	}

	out := mapDomainUserToSetIsActiveOutput(user, teamName)

	s.logger.Info("set is_active completed", map[string]any{
		"user_id":   out.UserId,
		"is_active": out.IsActive,
		"team_name": out.TeamName,
	})

	if in.IsActive {
		s.metrics.IncUserActivated()
	} else {
		s.metrics.IncUserDeactivated()
	}

	return out, nil
}

func (s *Service) GetUserReviews(ctx context.Context, in GetUserReviewsInput) (*GetUserReviewsOutput, error) {
	if err := validateGetUserReviewsInput(in); err != nil {
		s.logger.Error("get user reviews validation failed", map[string]any{
			"user_id": in.UserId,
			"error":   err.Error(),
		})
		return nil, err
	}

	s.logger.Info("get user reviews started", map[string]any{
		"user_id": in.UserId,
	})

	prs, err := s.prs.GetAllPrByUserId(ctx, in.UserId)
	if err != nil {
		s.logger.Error("get user reviews repository error", map[string]any{
			"user_id": in.UserId,
			"error":   err.Error(),
		})
		return nil, err
	}

	out := mapDomainPRsToGetUserReviewsOutput(in.UserId, prs)

	s.logger.Info("get user reviews completed", map[string]any{
		"user_id":  out.UserId,
		"pr_count": len(out.PullRequests),
	})

	return out, nil
}
