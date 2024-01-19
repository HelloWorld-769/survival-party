package handler

import (
	"main/server/request"
	"main/server/response"
	"main/server/services/user"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// UpdatePlayerInfoHandler Updates player info like username and avatar
//
// @Summary Updates player info
// @Description Updates player info like username and avatar
// @Tags Player
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Param loginDetails body request.UpdatePlayer true "Player Details"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /user-data [put]
func UpdatePlayerInfoHandler(ctx *gin.Context) {

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	var input request.UpdatePlayer
	err := utils.RequestDecoding(ctx, &input)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	user.UpdatePlayerService(ctx, userId.(string), input)
}

// GetPlayerStatsHandler Gets the player stats
//
// @Summary Get player game stats
// @Description Get the player game stats
// @Tags Player
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /stats [get]
func GetPlayerStatsHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	user.GetPlayerStatsService(ctx, userId.(string))
}

// GetOtherPlayerStatsHandler Gets the player stats
//
// @Summary Get game stats of ther players
// @Description Get game stats of ther players
// @Tags Player
// @Accept json
// @Produce json
// @Param playerId query string true "Id of the other player"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /user-data [get]
func GetOtherPlayerStatsHandler(ctx *gin.Context) {
	userId := ctx.Query("playerId")
	if userId == "" {
		response.ShowResponse("User id not found", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	user.GetPlayerStatsService(ctx, userId)
}
