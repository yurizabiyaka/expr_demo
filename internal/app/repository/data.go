package repository

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"

	"expr_demo/internal/pkg/model"

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
	err := d.db.QueryRow(sqlStatement, a.Created_at, a.Account, a.Amount_cents, a.Pos, a.Country).Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (d DataAccess) GetStringsFromData(column, where, having string, eventValues []interface{}) ([]string, error) {
	sqlStatement := getSqlStatement(column, where, having)
	rows, err := d.db.Query(sqlStatement, eventValues...)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("cannot select %s from data", column))
	}
	defer rows.Close()

	var strValues []string
	for rows.Next() {
		var str string
		if err = rows.Scan(&str); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("cannot scan %s from data", column))
		}
		strValues = append(strValues, str)
	}

	return strValues, nil
}

func getSqlStatement(column, where, having string) string {
	var whereDef, groupDef, havingDef string
	if where != "" {
		whereDef = "WHERE " + where
	}
	if having != "" { // having supposes group by
		groupDef = "GROUP BY " + column
		havingDef = "HAVING " + having
	}

	sqlStatement := fmt.Sprintf("SELECT %s FROM data %s %s %s", column, whereDef, groupDef, havingDef)

	return sqlStatement
}
