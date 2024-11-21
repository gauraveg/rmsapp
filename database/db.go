package database

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	RmsDB *sqlx.DB
)

type TxFn func(*sqlx.Tx) error

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

func WithTxn(logger *zap.Logger, fn TxFn) error {
	tx, err := RmsDB.Beginx()
	if err != nil {
		logger.Error("Cannot begin database transaction")
		return err
	}

	defer func() {
		if err != nil {
			logger.Error("Rolling back database transaction")
			err := tx.Rollback()
			if err != nil {
				logger.Error("Cannot rollback database transaction")
				return
			}
		} else {
			logger.Error("Committing database transaction")
			err := tx.Commit()
			if err != nil {
				logger.Error("Cannot commit database transaction")
				return
			}
		}
	}()

	logger.Info("Starting database transaction")
	err = fn(tx)
	return err
}

func ShutdownDatabase() error {
	return RmsDB.Close()
}
