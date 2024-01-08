package handler

import (
	"fmt"
	"main/server/request"
	"main/server/response"
	"main/server/services/settings"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

func GetSettingsHandler(ctx *gin.Context) {

	utils.SetHeader(ctx)

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//call the service
	settings.GetSettingsService(ctx, userId.(string))
}

func UpdateSettingsHandler(ctx *gin.Context) {

	utils.SetHeader(ctx)

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	var req request.UpdatePlayerSettingsRequest
	// fmt.Println("request", ctx.Request.Body)

	utils.RequestDecoding(ctx, &req)
	fmt.Println("req", req)

	//validation Check on request body fields
	err := validation.CheckValidation(&req)
	if err != nil {
		response.ShowResponse(err.Error(), 400, "Failure", "", ctx)
		return
	}

	//call the service
	settings.UpdateSettingsService(ctx, userId.(string), req)
}
