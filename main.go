package main

import (
	"os"

	"github.com/gauraveg/rmsapp/database"
	"github.com/gauraveg/rmsapp/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	//fmt.Println("Hello World")
	logger := logger.Init()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()
	zap.L().Info("Logger Initiated")
	godotenv.Load() // todo to explain how to get those value from env config and store it in a struct
	dbUrl := os.Getenv("DB_URL")
	port := os.Getenv("PORT")

	err := database.ConnectDB(dbUrl)
	if err != nil {
		//log.Printf("Failed to connect to database with error: %+v", err)
		zap.L().Error("Failed to connect to database", zap.Error(err))
		return
	}
	zap.L().Info("Db connection successful!")

	srv := RmsRouters()
	zap.L().Info("Server has started", zap.String("PORT", port))
	err = srv.Run(port)
	if err != nil {
		zap.L().Error("Failed to run server", zap.Error(err))
		return
	}

	err = database.ShutdownDatabase()
	if err != nil {
		zap.L().Error("failed to close database connection", zap.Error(err))
		return
	}
}
