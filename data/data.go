package data

import "github.com/jmoiron/sqlx"

type Data struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *Data {
	return &Data{
		DB: db,
	}
}
