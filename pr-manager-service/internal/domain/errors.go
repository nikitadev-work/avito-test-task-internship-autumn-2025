package domain

import "errors"

var (
	ErrEditMergedPR          = errors.New("cannot edit pull request: it is already merged")
	ErrMoreThanTwoReviewers  = errors.New("cannot add new reviewer: there are two reviewers")
	ErrAssignInactiveUser    = errors.New("cannot assign new reviewer to pr: user is inactive")
	ErrNoAvailableCandidates = errors.New("cannot assign new reviwer to pr: there is no available candidates")
	ErrPrAlreadyExists       = errors.New("cannot create new pr: pr already exists")
	ErrTeamAlreadyExists     = errors.New("cannot create new pr: pr already exists")
)
