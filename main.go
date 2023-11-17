package main

import (
	"log"
	"main/server"
	"main/server/db"
	"main/server/services/alert_service/twilio"
	"main/server/socket"
	"os"

	"github.com/joho/godotenv"
)

// @title Gin Demo App
// @version 1.0
// @description This is a demo version of Gin app.
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connection := db.InitDB()
	db.Transfer(connection)
	twilio.TwilioInit()
	socketServer := socket.SocketInit()
	defer socketServer.Close()
	app := server.NewServer(connection)
	server.ConfigureRoutes(app)

	if err := app.Run(os.Getenv("PORT")); err != nil {
		log.Print(err)
	}
}
