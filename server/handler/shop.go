package handler

import (
	"main/server/request"
	"main/server/response"
	"main/server/services/shop"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// GetStoreHandler Gets the details of the shop
//
// @Summary Gets shop details
// @Description Gets shop details
// @Tags Store
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

// BuyFromStoreHandler Buy the assests from the shop
//
// @Summary Buy things
// @Description Buy the assests from the shop
// @Tags Store
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Param loginDetails body request.BuyStoreRequest true "shop Details"
// @Success 200 {object} response.Success "Sucess"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /buy-store [post]
func BuyFromStoreHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	var input request.BuyStoreRequest
	err := utils.RequestDecoding(ctx, &input)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	shop.BuyFromStoreService(ctx, userId.(string), input)

}
