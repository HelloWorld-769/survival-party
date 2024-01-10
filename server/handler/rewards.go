package handler

import (
	"main/server/request"
	"main/server/response"
	"main/server/services/rewards"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func GetPlayerLevelRewardsHandler(ctx *gin.Context) {

	utils.SetHeader(ctx)

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	rewards.GetPlayerLevelRewards(ctx, userId.(string))
}

func PlayerLevelRewardCollectionHandler(ctx *gin.Context) {

	utils.SetHeader(ctx)

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	var req request.PlayerLevelRewardCollectRequest
	err := utils.RequestDecoding(ctx, &req)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	rewards.PlayerLevelRewardCollect(ctx, userId.(string), req)

}
