package model

import (
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	ID           uuid.UUID `db:"id"`
	Created_at   time.Time `db:"created_at"`
	Account      string    `db:"account"`
	Amount_cents uint64    `db:"amount_cents"`
	Pos          string    `db:"pos"`
	Country      string    `db:"country"`
}
