package usecase

func validateCreateTeamInput(in CreateTeamInput) error {
	if in.TeamName == "" {
		return ErrTeamNameRequired
	}
	return nil
}

func validateGetTeamInput(in GetTeamInput) error {
	if in.TeamName == "" {
		return ErrTeamNameRequired
	}
	return nil
}

func validateSetIsActiveInput(in SetIsActiveInput) error {
	if in.UserId == "" {
		return ErrUserIdRequired
	}
	return nil
}

func validateGetUserReviewsInput(in GetUserReviewsInput) error {
	if in.UserId == "" {
		return ErrUserIdRequired
	}
	return nil
}

func validateCreatePullRequestInput(in CreatePullRequestInput) error {
	if in.PullRequestId == "" {
		return ErrPullRequestIdRequired
	}
	if in.PullRequestName == "" {
		return ErrPullRequestNameRequired
	}
	if in.AuthorId == "" {
		return ErrAuthorIdRequired
	}
	return nil
}

func validateMergePullRequestInput(in MergePullRequestInput) error {
	if in.PullRequestId == "" {
		return ErrPullRequestIdRequired
	}
	return nil
}

func validateReassignReviewerInput(in ReassignReviewerInput) error {
	if in.PullRequestId == "" {
		return ErrPullRequestIdRequired
	}
	if in.OldUserId == "" {
		return ErrOldUserIdRequired
	}
	return nil
}
