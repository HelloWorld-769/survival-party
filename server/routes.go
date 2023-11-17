package server

import (
	_ "main/docs"
	"main/server/gateway"
	"main/server/handler"
	gomail "main/server/services/alert_service/Gomail"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {

	//Allowing CORS
	server.engine.Use(gateway.CORSMiddleware())

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.engine.POST("/signup", handler.SignupHandler)
	server.engine.POST("/login", handler.LoginHandler)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.engine.POST("/send-email", gomail.SendEmailOtpService)
}
