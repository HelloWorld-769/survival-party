package server

import (
	_ "main/docs"
	"main/server/gateway"
	"main/server/handler"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {

	//Allowing CORS
	server.engine.Use(gateway.CORSMiddleware())
	//swagger route
	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Auth routes
	server.engine.POST("/users/sign-up", handler.SignupHandler)
	server.engine.POST("/users/sign-in", handler.LoginHandler)
	server.engine.DELETE("/users/sign-out", gateway.UserAuthorization, handler.SignoutHandler)
	server.engine.POST("/users/social-login", handler.SocialLoginHandler)
	server.engine.GET("/users/email-verify", handler.VerifyEmailHandler)
	server.engine.POST("/send-otp", handler.SendOtpHandler)
	server.engine.POST("/reset-password", handler.ResetPasswordHandler)
	server.engine.POST("/check-otp", handler.CheckOtpHandler)

	//Player Routes
	server.engine.PUT("/user-data", gateway.UserAuthorization, handler.UpdatePlayerInfoHandler)
	server.engine.GET("/get-settings", gateway.UserAuthorization, handler.GetSettingsHandler)
	server.engine.PUT("/update-settings", gateway.UserAuthorization, handler.UpdateSettingsHandler)
	server.engine.GET("/stats", gateway.UserAuthorization, handler.GetPlayerStatsHandler)
	server.engine.GET("/store", gateway.UserAuthorization, handler.GetStoreHandler)

	//Level rewards
	server.engine.GET("/get-level-rewards", gateway.UserAuthorization, handler.GetPlayerLevelRewardsHandler)
	server.engine.POST("/level-reward-collect", gateway.UserAuthorization, handler.PlayerLevelRewardCollectionHandler)

	//Store Routes
	server.engine.GET("/store", gateway.UserAuthorization, handler.GetStoreHandler)
	server.engine.POST("/buy-store", gateway.UserAuthorization, handler.BuyFromStoreHandler)
	server.engine.GET("/popupoffers", handler.GetPopupHandler)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//daily rewards
	server.engine.PUT("/collect-daily-rewards", gateway.UserAuthorization, handler.CollectDailyRewardHandler)
	server.engine.GET("/daily-rewards", gateway.UserAuthorization, handler.GetUserDailyRewardDataHandler)
}
