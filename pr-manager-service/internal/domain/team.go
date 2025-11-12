package domain

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
