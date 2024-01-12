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
	server.engine.POST("/api/v1/users", handler.SignupHandler)
	server.engine.POST("/api/v1/users/sign-in", handler.LoginHandler)
	server.engine.DELETE("/api/v1/users/sign-out", gateway.UserAuthorization, handler.SignoutHandler)
	server.engine.POST("/api/v1/users/social-login", handler.SocialLoginHandler)
	server.engine.PUT("/api/v1/users/email-verify", handler.VerifyEmailHandler)
	server.engine.POST("/api/v1/send-otp", handler.SendOtpHandler)
	server.engine.POST("/api/v1/reset-password", handler.ResetPasswordHandler)
	server.engine.POST("/api/v1/check-otp", handler.CheckOtpHandler)

	//Player Routes
<<<<<<< HEAD
	server.engine.PUT("/api/v1/userData", gateway.UserAuthorization, handler.UpdatePlayerInfoHandler)
	server.engine.GET("/api/v1/get-settings", gateway.UserAuthorization, handler.GetSettingsHandler)
	server.engine.PUT("/api/v1/update-settings", gateway.UserAuthorization, handler.UpdateSettingsHandler)
	server.engine.GET("/api/v1/stats", gateway.UserAuthorization, handler.GetPlayerStatsHandler)
	server.engine.GET("/api/v1/store", gateway.UserAuthorization, handler.GetStoreHandler)
	server.engine.POST("/api/v1/buy-store", gateway.UserAuthorization, handler.BuyFromStoreHandler)
=======
	server.engine.PUT("/user-data", gateway.UserAuthorization, handler.UpdatePlayerInfoHandler)
	server.engine.GET("/get-settings", gateway.UserAuthorization, handler.GetSettingsHandler)
	server.engine.PUT("/update-settings", gateway.UserAuthorization, handler.UpdateSettingsHandler)
	server.engine.GET("/stats", gateway.UserAuthorization, handler.GetPlayerStatsHandler)
<<<<<<< HEAD
	server.engine.GET("/store", gateway.UserAuthorization, handler.GetStoreHandler)
<<<<<<< HEAD
>>>>>>> a23bf06 (starter pack fix)
=======
	server.engine.POST("/buy-store", gateway.UserAuthorization, handler.BuyFromStoreHandler)
>>>>>>> 813b225 (Feat:added shop buy handler and starter pack)
=======
>>>>>>> c4acd51 (feat: added pop up offer route)

	server.engine.GET("/api/v1/get-level-rewards", gateway.UserAuthorization, handler.GetPlayerLevelRewardsHandler)
	server.engine.POST("/api/v1/level-reward-collect", gateway.UserAuthorization, handler.PlayerLevelRewardCollectionHandler)

	//Daily Goals
	server.engine.GET("/api/v1/get-daily-goals", gateway.UserAuthorization, handler.GetDailyGoalsHandler)
	server.engine.POST("/api/v1/skip-daily-goal", gateway.UserAuthorization, handler.SkipGoalHandler)
	server.engine.POST("/api/v1/claim-daily-goal", gateway.UserAuthorization, handler.ClaimDailyGoalHandler)

	//Store Routes
	server.engine.GET("/store", gateway.UserAuthorization, handler.GetStoreHandler)
	server.engine.POST("/buy-store", gateway.UserAuthorization, handler.BuyFromStoreHandler)
	server.engine.GET("/popupoffers", handler.GetPopupHandler)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//daily rewards
	server.engine.PUT("/collect-daily-rewards", gateway.UserAuthorization, handler.CollectDailyRewardHandler)
	server.engine.GET("/daily-rewards", gateway.UserAuthorization, handler.GetUserDailyRewardDataHandler)
	//Store Routes
	server.engine.GET("/store", gateway.UserAuthorization, handler.GetStoreHandler)
	server.engine.POST("/buy-store", gateway.UserAuthorization, handler.BuyFromStoreHandler)
	server.engine.GET("/popupoffers", handler.GetPopupHandler)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
