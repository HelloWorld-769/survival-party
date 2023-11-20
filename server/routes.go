package server

import (
	_ "main/docs"
	gomail "main/server/services/alert_service/Gomail"

	"github.com/gin-gonic/gin"
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
	server.engine.POST("/send-sms", handler.TwilioServiceHnadler)

	server.engine.GET("/ping", Pong)

}

// For server Testing(acknowledgement)
func Pong(ctx *gin.Context) {

	msg := "Pong"
	ctx.Writer.Write([]byte(msg))
}
