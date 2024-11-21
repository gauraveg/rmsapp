package main

import (
	"log"
	"os"

	"github.com/gauraveg/rmsapp/database"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	//fmt.Println("Hello World")
	err := godotenv.Load()
	if err != nil {
		return
	}
	// todo to explain how to get those value from env config and store it in a struct
	dbUrl := os.Getenv("DB_URL")
	port := os.Getenv("PORT")

	err = database.ConnectDB(dbUrl)
	if err != nil {
		log.Printf("Failed to connect to database with error: %+v", err)
		zap.L().Error("Failed to connect to database", zap.Error(err))
		return
	}
	log.Printf("Db connection successful!")

	srv := RmsRouters()
	log.Printf("Server has started at PORT %v", port)
	err = srv.Run(port)
	if err != nil {
		log.Printf("Failed to run server. Error: %v", err)
		return
	}

	err = database.ShutdownDatabase()
	if err != nil {
		log.Printf("failed to close database connection. Error: %v", err)
		return
	}
}
