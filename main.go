package main

import (
	"log"
	"os"

	"github.com/gauraveg/rmsapp/database"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	//fmt.Println("Hello World")
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	port := os.Getenv("PORT")

	err := database.ConnectDB(dbUrl)
	if err != nil {
		log.Printf("Failed to connect to database with error: %+v", err)
		return
	}
	log.Println("Db connection successful!")

	srv := RmsRouters()
	log.Printf("Server has started on port: %v", port)
	err = srv.Run(port)
	if err != nil {
		log.Printf("Failed to run server with error: %+v", err)
		return
	}

	err = database.ShutdownDatabase()
	if err != nil {
		log.Printf("failed to close database connection \n %v", err)
		return
	}
}
