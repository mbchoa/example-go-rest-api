package main

import (
	"log"
	"os"

	"github.com/mbchoa/example-go-rest-api/3-rest-database/controllers"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	}

	server := controllers.Server{}

	server.ConnectDb(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"))
	server.StartServer(":8080")
}
