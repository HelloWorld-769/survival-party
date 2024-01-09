package handler

import (
	"main/server/request"
	"main/server/response"
	"main/server/services/player"
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
// @Param loginDetails body request.UpdatePlayer true "Player Details"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /userData [post]
func UpdatePlayerInfoHandler(ctx *gin.Context) {

	userId, exists := ctx.Get("user_id")
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

	player.UpdatePlayerService(ctx, userId.(string), input)
}

func GetPlayerStatsHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	player.GetPlayerStatsService(ctx, userId.(string))
}
