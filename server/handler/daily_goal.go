package handler

import (
	"main/server/request"
	"main/server/response"
	dailygoal "main/server/services/daily_goal"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// GetDailyGoalsHandler Gets the daily goals for given player
//
// @Summary Get the daily goals
// @Description Gets the daily goals for given player
// @Tags DailyGoal
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Success 200 {object} response.Success "Sucess"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /get-daily-goals [get]
func GetDailyGoalsHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	dailygoal.GetDailyGoalsService(ctx, userId.(string))

}

func SkipGoalHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	var input request.GoalRequest
	err := utils.RequestDecoding(ctx, &input)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	dailygoal.SkipGoalService(ctx, userId.(string), input)

}

func ClaimDailyGoalHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	var input request.GoalRequest
	err := utils.RequestDecoding(ctx, &input)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	dailygoal.ClaimDailyGoalService(ctx, userId.(string), input)
}
