package usecase

import (
	"context"
	"database/sql"
	"errors"
)

// Teams

func (s *Service) CreateTeam(ctx context.Context, in CreateTeamInput) (*CreateTeamOutput, error) {
	if err := validateCreateTeamInput(in); err != nil {
		s.logger.Error("create team validation failed", map[string]any{
			"team_name": in.TeamName,
			"error":     err.Error(),
		})
		return nil, err
	}

	s.logger.Info("create team started", map[string]any{
		"team_name":     in.TeamName,
		"members_count": len(in.Members),
	})

	members := mapTeamMembersDTOToDomain(in.Members)

	err := s.teams.CreateTeam(ctx, in.TeamName, members)
	if err != nil {
		s.logger.Error("create team repository error", map[string]any{
			"team_name": in.TeamName,
			"error":     err.Error(),
		})
		return nil, err
	}

	out := &CreateTeamOutput{
		TeamName: in.TeamName,
		Members:  in.Members,
	}

	s.logger.Info("create team completed", map[string]any{
		"team_name":     out.TeamName,
		"members_count": len(out.Members),
	})

	s.metrics.IncTeamCreated()

	return out, nil
}

func (s *Service) GetTeam(ctx context.Context, in GetTeamInput) (*GetTeamOutput, error) {
	if err := validateGetTeamInput(in); err != nil {
		s.logger.Error("get team validation failed", map[string]any{
			"team_name": in.TeamName,
			"error":     err.Error(),
		})
		return nil, err
	}

	s.logger.Info("get team started", map[string]any{
		"team_name": in.TeamName,
	})

	_, members, err := s.teams.GetTeam(ctx, in.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("get team: team not found", map[string]any{
				"team_name": in.TeamName,
				"error":     err.Error(),
			})
			return nil, err
		}

		s.logger.Error("get team repository error", map[string]any{
			"team_name": in.TeamName,
			"error":     err.Error(),
		})
		return nil, err
	}

	outMembers := mapDomainUsersToTeamMembersDTO(members)

	out := &GetTeamOutput{
		TeamName: in.TeamName,
		Members:  outMembers,
	}

	s.logger.Info("get team completed", map[string]any{
		"team_name":     out.TeamName,
		"members_count": len(out.Members),
	})

	return out, nil
}
