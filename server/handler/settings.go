package handler

import (
	"main/server/response"
	"main/server/services/settings"
	"main/server/utils"

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
