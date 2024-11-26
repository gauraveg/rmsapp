package main

import (
	"github.com/gauraveg/rmsapp/logger"
	"os"

	"github.com/gauraveg/rmsapp/database"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	//fmt.Println("Hello World")
	loggers := logger.LogWrapperInit()
	err := godotenv.Load()
	if err != nil {
		return
	}
	// todo to explain how to get those value from env config and store it in a struct
	dbUrl := os.Getenv("DB_URL")
	port := os.Getenv("PORT")

	err = database.ConnectDB(dbUrl)
	if err != nil {
		loggers.Info(loggers.Logger, "Failed to connect to database with error: %+v", err)
		loggers.Error(loggers.Logger, "Failed to connect to database", zap.Error(err))
		return
	}
	loggers.Info(loggers.Logger, "Db connection successful!")

	srv := RmsRouters(loggers)
	loggers.Info(loggers.Logger, "Server has started at PORT %v", port)
	err = srv.Run(port)
	if err != nil {
		loggers.Error(loggers.Logger, "Failed to run server. Error: %v", err)
		return
	}

	err = database.ShutdownDatabase()
	if err != nil {
		loggers.Error(loggers.Logger, "failed to close database connection. Error: %v", err)
		return
	}
}
