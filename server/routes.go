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
	server.engine.DELETE("/users/sign_out", gateway.UserAuthorization, handler.SignoutHandler)
	server.engine.POST("/users/social_login", handler.SocialLoginHandler)
	server.engine.PUT("/users/email-verify", handler.VerifyEmailHandler)

	//Player Routes
	server.engine.PUT("/userData", gateway.AdminAuthorization, handler.UpdatePlayerInfoHandler)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.engine.POST("/send_otp", handler.SendOtpHandler)
	server.engine.POST("/check_otp", handler.CheckOtpHandler)

	server.engine.POST("/reset_password", handler.ResetPasswordHandler)

	server.engine.GET("/get_settings", gateway.UserAuthorization, handler.GetSettingsHandler)
	server.engine.PUT("/update_settings", gateway.UserAuthorization, handler.UpdateSettingsHandler)

	// server.engine.POST("/send-email", gomail.SendEmailOtpService)

}

// For server Testing(acknowledgement)
func Pong(ctx *gin.Context) {

	msg := "Pong"
	ctx.Writer.Write([]byte(msg))
}
