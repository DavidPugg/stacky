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
		fmt.Sprintf("%s://%s", viper.GetString("DB_DRIVER"), viper.GetString("DB_URL")),
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
