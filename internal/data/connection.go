package data

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func DBconnect() *sqlx.DB {
	db, err := sqlx.Connect("mysql", strings.Split(viper.GetString("DB_URL"), "://")[1])
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
