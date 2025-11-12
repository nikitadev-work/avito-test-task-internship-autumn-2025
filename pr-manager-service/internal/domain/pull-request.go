package domain

import (
	"time"

	"github.com/google/uuid"
)

var Status int

const (
	Open = iota
	Merged
)

type PulLRequest struct {
	Id                uuid.UUID
	Name              string
	AuthorId          uuid.UUID
	StatusId          int
	NeedMoreReviewers bool
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ReviewerAssignment struct {
	UserId        uuid.UUID
	PullRequestId uuid.UUID
	Slot          int
}
