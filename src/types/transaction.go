package types

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID                uuid.UUID
	SenderAccountID   uuid.UUID
	ReceiverAccountID uuid.UUID
	Amount            float64
	CreatedAt         time.Time
}
