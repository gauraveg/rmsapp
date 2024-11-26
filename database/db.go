package database

import (
	"context"

	"github.com/gauraveg/rmsapp/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
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

func WithTxn(ctx context.Context, loggers *logger.ZapLogger, fn TxFn) error {
	tx, err := RmsDB.Beginx()
	if err != nil {
		loggers.ErrorWithContext(ctx, "Cannot begin database transaction")
		return err
	}

	defer func() {
		if err != nil {
			loggers.InfoWithContext(ctx, "Rolling back database transaction")
			err := tx.Rollback()
			if err != nil {
				loggers.ErrorWithContext(ctx, "Cannot rollback database transaction")
				return
			}
		} else {
			loggers.InfoWithContext(ctx, "Committing database transaction")
			err := tx.Commit()
			if err != nil {
				loggers.ErrorWithContext(ctx, "Cannot commit database transaction")
				return
			}
		}
	}()

	loggers.InfoWithContext(ctx, "Starting database transaction")
	err = fn(tx)
	return err
}

func ShutdownDatabase() error {
	return RmsDB.Close()
}
