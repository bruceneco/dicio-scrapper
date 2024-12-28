package word

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	ID         uuid.UUID
	Content    string
	Definition string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
