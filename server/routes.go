package server

import (
	_ "main/docs"
	"main/server/gateway"
	"main/server/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {

	//Allowing CORS
	server.engine.Use(gateway.CORSMiddleware())

	//Auth routes
	server.engine.POST("/users", handler.SignupHandler)
	server.engine.POST("/users/sign_in", handler.LoginHandler)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.engine.POST("/send-otp", handler.SendOtpHandler)
	server.engine.POST("/reset-password", handler.ResetPasswordHandler)

	// server.engine.POST("/send-email", gomail.SendEmailOtpService)

}

// For server Testing(acknowledgement)
func Pong(ctx *gin.Context) {

	msg := "Pong"
	ctx.Writer.Write([]byte(msg))
}
