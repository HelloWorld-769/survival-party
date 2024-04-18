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
		response.ShowResponse("Incorrect username or password", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	dailygoal.GetDailyGoalsService(ctx, userId.(string))

}

// SkipGoalHandler Skips the daily goals of the player
//
// @Summary Get the daily goals
// @Description Gets the daily goals for given player
// @Tags DailyGoal
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Param loginDetails body request.GoalRequest true "Goal id"
// @Success 200 {object} response.Success "Sucess"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /skip-daily-goal [post]
func SkipGoalHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("Incorrect username or password", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
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

// ClaimDailyGoalHandler Claims the rewards when all daily goals are completed
//
// @Summary Claims the rewards
// @Description Claims the rewards when all daily goals are completed
// @Tags DailyGoal
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Success 200 {object} response.Success "Sucess"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /claim-daily-goal [post]
func ClaimDailyGoalHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("Incorrect username or password", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	dailygoal.ClaimDailyGoalService(ctx, userId.(string))
}
