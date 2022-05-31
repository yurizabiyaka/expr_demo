package repository

import (
	"database/sql"
	"expr_demo/internal/app/model"
	"github.com/google/uuid"
)

type DataAccess struct {
	db *sql.DB
}

func NewDataRepo(db *sql.DB) *DataAccess {
	return &DataAccess{db: db}
}

const (
	saveExpr = `
INSERT INTO data(created_at, account, amount_cents, pos, country)
VALUES ($1, $2, $3, $4, $5)
RETURNING id`
)

func (d DataAccess) Save(a model.Auth) (uuid.UUID, error) {
	sqlStatement := saveExpr
	var id uuid.UUID
	err := d.db.QueryRow(sqlStatement, a.CreatedAt, a.Account, a.AmountCents, a.POS, a.CountryMnem).Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
