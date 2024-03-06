package handler

import (
	"main/server/request"
	"main/server/response"
	"main/server/services/user"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// @Summary Deduct Amount
// @Description Deducts the specified amount from the user's account
// @Tags Game
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param body body request.DeductAmount true "Deduct Amount Request"
// @Success 200 {object} response.Success
// @Failure 400 {object} response.Success
// @Failure 401 {object} response.Success
// @Router /deduct-amount [put]
func DeductAmountHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	var input request.DeductAmount
	err := utils.RequestDecoding(ctx, &input)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	user.DeductAmountService(ctx, userId.(string), input)
}
