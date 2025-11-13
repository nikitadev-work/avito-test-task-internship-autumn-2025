package usecase

import "pr-manager-service/internal/domain"

// Teams

func mapTeamMembersDTOToDomain(members []TeamMemberDTO) []domain.User {
	result := make([]domain.User, 0, len(members))
	for _, m := range members {
		result = append(result, domain.User{
			UserId:   m.UserId,
			UserName: m.UserName,
			IsActive: m.IsActive,
		})
	}
	return result
}

func mapDomainUsersToTeamMembersDTO(users []domain.User) []TeamMemberDTO {
	result := make([]TeamMemberDTO, 0, len(users))
	for _, u := range users {
		result = append(result, TeamMemberDTO{
			UserId:   u.UserId,
			UserName: u.UserName,
			IsActive: u.IsActive,
		})
	}
	return result
}

// Users

func mapDomainUserToSetIsActiveOutput(user *domain.User, teamName string) *SetIsActiveOutput {
	return &SetIsActiveOutput{
		UserId:   user.UserId,
		UserName: user.UserName,
		TeamName: teamName,
		IsActive: user.IsActive,
	}
}

func mapDomainPRsToGetUserReviewsOutput(userId string, prs []domain.PullRequest) *GetUserReviewsOutput {
	result := make([]PullRequestShortDTO, 0, len(prs))
	for _, pr := range prs {
		result = append(result, PullRequestShortDTO{
			PullRequestId:   pr.PullRequestId,
			PullRequestName: pr.PullRequestName,
			AuthorId:        pr.AuthorId,
			Status:          statusString(pr.StatusId),
		})
	}
	return &GetUserReviewsOutput{
		UserId:       userId,
		PullRequests: result,
	}
}

// Pull requests

func mapCreatePRInputToDomain(in CreatePullRequestInput, assigned []string) *domain.PullRequest {
	return &domain.PullRequest{
		PullRequestId:     in.PullRequestId,
		PullRequestName:   in.PullRequestName,
		AuthorId:          in.AuthorId,
		StatusId:          1, // 1 - OPEN, 2 - MERGED
		AssignedReviewers: assigned,
	}
}

func mapDomainPRToDTO(pr *domain.PullRequest) PullRequestDTO {
	return PullRequestDTO{
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		AuthorId:          pr.AuthorId,
		Status:            statusString(pr.StatusId),
		AssignedReviewers: pr.AssignedReviewers,
	}
}

// Other

func statusString(statusId int) string {
	switch statusId {
	case 1:
		return "OPEN"
	case 2:
		return "MERGED"
	default:
		return "OPEN"
	}
}
