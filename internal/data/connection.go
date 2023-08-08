package data

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func DBconnect() *sqlx.DB {
	db, err := sqlx.Connect(viper.GetString("DB_DRIVER"), viper.GetString("DB_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
