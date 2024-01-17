package server

import (
	_ "main/docs"
	"main/server/gateway"
	"main/server/handler"
	"main/server/services/rewards"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {

	//Allowing CORS
	server.engine.Use(gateway.CORSMiddleware())
	//swagger route
	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Auth routes
	server.engine.POST("/api/v1/users/sign-up", handler.SignupHandler)
	server.engine.POST("/api/v1/users/sign-in", handler.LoginHandler)
	server.engine.DELETE("/api/v1/users/sign-out", gateway.UserAuthorization, handler.SignoutHandler)
	server.engine.POST("/api/v1/users/social-login", handler.SocialLoginHandler)
	server.engine.GET("/api/v1/users/email-verify", handler.VerifyEmailHandler)
	server.engine.POST("/api/v1/send-otp", handler.SendOtpHandler)
	server.engine.POST("/api/v1/reset-password", handler.ResetPasswordHandler)
	server.engine.POST("/api/v1/check-otp", handler.CheckOtpHandler)

	//Player Routes

	server.engine.PUT("/api/v1/user-data", gateway.UserAuthorization, handler.UpdatePlayerInfoHandler)
	server.engine.GET("/api/v1/get-settings", gateway.UserAuthorization, handler.GetSettingsHandler)
	server.engine.PUT("/api/v1/update-settings", gateway.UserAuthorization, handler.UpdateSettingsHandler)
	server.engine.GET("/api/v1/stats", gateway.UserAuthorization, handler.GetPlayerStatsHandler)

	server.engine.GET("/api/v1/get-level-rewards", gateway.UserAuthorization, handler.GetPlayerLevelRewardsHandler)
	server.engine.POST("/api/v1/level-reward-collect", gateway.UserAuthorization, handler.PlayerLevelRewardCollectionHandler)

	//Daily Goals
	server.engine.GET("/api/v1/get-daily-goals", gateway.UserAuthorization, handler.GetDailyGoalsHandler)
	server.engine.POST("/api/v1/skip-daily-goal", gateway.UserAuthorization, handler.SkipGoalHandler)
	server.engine.POST("/api/v1/claim-daily-goal", gateway.UserAuthorization, handler.ClaimDailyGoalHandler)

	//Store Routes
	server.engine.GET("/api/v1/store", gateway.UserAuthorization, handler.GetStoreHandler)
	server.engine.POST("/api/v1/buy-store", gateway.UserAuthorization, handler.BuyFromStoreHandler)
	server.engine.GET("/api/v1/popupoffers", handler.GetPopupHandler)

	//daily rewards
	server.engine.PUT("/api/v1/collect-daily-rewards", gateway.UserAuthorization, handler.CollectDailyRewardHandler)
	server.engine.GET("/api/v1/daily-rewards", gateway.UserAuthorization, handler.GetUserDailyRewardDataHandler)
	server.engine.GET("/api/v1/timeLeft", rewards.DailyRewardTimeLeft)
}
