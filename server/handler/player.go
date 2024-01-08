package handler

import (
	"main/server/request"
	"main/server/response"
	"main/server/services/player"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

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
