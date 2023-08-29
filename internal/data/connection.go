package data

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func DBconnect() *sqlx.DB {
	db, err := sqlx.Connect(
		viper.GetString("DB_DRIVER"),
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			viper.GetString("DB_HOST"),
			viper.GetString("DB_PORT"),
			viper.GetString("DB_USER"),
			viper.GetString("DB_PASSWORD"),
			viper.GetString("DB_NAME"),
			viper.GetString("DB_SSLMODE"),
		),
	)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
