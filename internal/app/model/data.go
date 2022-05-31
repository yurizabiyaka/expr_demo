package model

import (
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	ID          uuid.UUID `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Account     string    `db:"account"`
	AmountCents uint64    `db:"amount_cents"`
	POS         string    `db:"pos"`
	CountryMnem string    `db:"country"`
}
