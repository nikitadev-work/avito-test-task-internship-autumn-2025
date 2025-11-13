package usecase

import "errors"

func validateCreateTeamInput(in CreateTeamInput) error {
	if in.TeamName == "" {
		return errors.New("team_name is required")
	}
	return nil
}

func validateGetTeamInput(in GetTeamInput) error {
	if in.TeamName == "" {
		return errors.New("team_name is required")
	}
	return nil
}

func validateSetIsActiveInput(in SetIsActiveInput) error {
	if in.UserId == "" {
		return errors.New("user_id is required")
	}
	return nil
}

func validateGetUserReviewsInput(in GetUserReviewsInput) error {
	if in.UserId == "" {
		return errors.New("user_id is required")
	}
	return nil
}

func validateCreatePullRequestInput(in CreatePullRequestInput) error {
	if in.PullRequestId == "" {
		return errors.New("pull_request_id is required")
	}
	if in.PullRequestName == "" {
		return errors.New("pull_request_name is required")
	}
	if in.AuthorId == "" {
		return errors.New("author_id is required")
	}
	return nil
}

func validateMergePullRequestInput(in MergePullRequestInput) error {
	if in.PullRequestId == "" {
		return errors.New("pull_request_id is required")
	}
	return nil
}

func validateReassignReviewerInput(in ReassignReviewerInput) error {
	if in.PullRequestId == "" {
		return errors.New("pull_request_id is required")
	}
	if in.OldUserId == "" {
		return errors.New("old_user_id is required")
	}
	return nil
}
