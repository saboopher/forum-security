package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Session struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Token      string
	ExpireTime time.Time
}
