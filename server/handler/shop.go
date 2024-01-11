package handler

import (
	"main/server/response"
	"main/server/services/shop"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// GetStoreHandler Gets the details of the shop
//
// @Summary Gets shop details
// @Description Gets shop details
// @Tags Player
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Success 200 {object} response.Success "Sucess"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /store [get]
func GetStoreHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	shop.GetStoreService(ctx, userId.(string))

}
