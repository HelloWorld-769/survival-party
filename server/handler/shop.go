package handler

import (
	"main/server/response"
	"main/server/services/shop"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func GetStoreHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	shop.GetStoreService(ctx, userId.(string))

}
