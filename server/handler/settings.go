package handler

import (
	"fmt"
	"main/server/request"
	"main/server/response"
	"main/server/services/user"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

// GetSettingsHandler Gets the current settings of that player
//
// @Summary Gets the settings
// @Description Gets the current settings of that player
// @Tags Player
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Success 200 {object} response.Success "Sucess"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /get_settings [get]
func GetSettingsHandler(ctx *gin.Context) {

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	//call the service
	user.GetSettingsService(ctx, userId.(string))
}

// UpdateSettingsHandler Updates the settings of that player
//
// @Summary Updates setting
// @Description Updates the game settings of that player
// @Tags Player
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Param loginDetails body request.PlayerLevelRewardCollectRequest true "Player Details"
// @Success 200 {object} response.Success "Sucess"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /update_settings [put]
func UpdateSettingsHandler(ctx *gin.Context) {

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	var req request.UpdatePlayerSettingsRequest
	// fmt.Println("request", ctx.Request.Body)

	utils.RequestDecoding(ctx, &req)
	fmt.Println("req", req)

	//validation Check on request body fields
	err := validation.CheckValidation(&req)
	if err != nil {
		response.ShowResponse(err.Error(), 400, utils.FAILURE, nil, ctx)
		return
	}

	//call the service
	user.UpdateSettingsService(ctx, userId.(string), req)
}
