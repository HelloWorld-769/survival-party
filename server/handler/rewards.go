package handler

import (
	"main/server/request"
	"main/server/response"
	"main/server/services/rewards"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// GetPlayerLevelRewardsHandler Gets the rewards according to player level
//
// @Summary Gets reward list
// @Description Gets the rewards according to player level
// @Tags Player
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /get-level-rewards [get]
func GetPlayerLevelRewardsHandler(ctx *gin.Context) {

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	rewards.GetPlayerLevelRewards(ctx, userId.(string))
}

// PlayerLevelRewardCollectionHandler Collects the reward for that level
//
// @Summary Collects Reward
// @Description Collects the reward for a level of that player
// @Tags Player
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Param loginDetails body request.PlayerLevelRewardCollectRequest true "Player Details"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /level-reward-collect [post]
func PlayerLevelRewardCollectionHandler(ctx *gin.Context) {

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
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

func CollectDailyRewardHandler(ctx *gin.Context) {

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	rewards.CollectDailyReward(ctx, userId.(string))
}

func GetUserDailyRewardDataHandler(ctx *gin.Context) {

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	rewards.GetUserDailyRewardData(ctx, userId.(string))

}
