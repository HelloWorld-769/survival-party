package main

import (
	"log"
	"main/server"
	"main/server/db"
	"main/server/handler"
	"main/server/services"
	"main/server/services/alert_service/twilio"
	"main/server/services/rewards"
	"main/server/socket"
	"os"

	"github.com/joho/godotenv"
)

// @title Survival Party
// @version 1.0
// @description This is the api documentation of survival party game
// @BasePath /api/v1/
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connection := db.InitDB()
	db.Transfer(connection)
	twilio.TwilioInit(os.Getenv("TWILIO_AUTH_TOKEN"))
	socketServer := socket.SocketInit()
	defer socketServer.Close()
	app := server.NewServer(connection)
	server.ConfigureRoutes(app)
	handler.AddDummyDataHandler()
	handler.StartCron()
	rewards.AddPlayerLevel()
	services.AddDummyUsers()

	if err := app.Run(os.Getenv("PORT")); err != nil {
		log.Print(err)
	}
}
