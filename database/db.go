package database

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	RmsDB *sqlx.DB
)

func ConnectDB(dbUrl string) error {
	conn, err := sqlx.Open("pgx", dbUrl)
	if err != nil {
		return err
	}

	err = conn.Ping()
	if err != nil {
		closeEr := conn.Close()
		if closeEr != nil {
			return closeEr
		}
		return err
	}
	RmsDB = conn
	return nil
}

func ShutdownDatabase() error {
	return RmsDB.Close()
}
