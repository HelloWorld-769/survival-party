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

	//Auth routes
	server.engine.POST("/users", handler.SignupHandler)
	server.engine.POST("/users/sign_in", handler.LoginHandler)
	server.engine.DELETE("/users/sign_out", gateway.UserAuthorization, handler.SignoutHandler)
	server.engine.POST("/users/social_login", handler.SocialLoginHandler)
	server.engine.PUT("/users/email-verify", handler.VerifyEmailHandler)
	server.engine.POST("/send-otp", handler.SendOtpHandler)
	server.engine.POST("/reset-password", handler.ResetPasswordHandler)
	server.engine.POST("/check_otp", handler.CheckOtpHandler)

	//Player Routes
	server.engine.PUT("/userData", gateway.UserAuthorization, handler.UpdatePlayerInfoHandler)
	server.engine.GET("/get_settings", gateway.UserAuthorization, handler.GetSettingsHandler)
	server.engine.PUT("/update_settings", gateway.UserAuthorization, handler.UpdateSettingsHandler)
	server.engine.GET("/stats", gateway.UserAuthorization, handler.GetPlayerStatsHandler)

	server.engine.GET("/get_level_rewards", gateway.UserAuthorization, handler.GetPlayerLevelRewardsHandler)
	server.engine.POST("/level_reward_collect", gateway.UserAuthorization, handler.PlayerLevelRewardCollectionHandler)

	//Store Routes
	server.engine.GET("/store", gateway.UserAuthorization, handler.GetStoreHandler)
	server.engine.POST("/buy-store", gateway.UserAuthorization, handler.BuyFromStoreHandler)
	server.engine.GET("/popupoffers", handler.GetPopupHandler)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
