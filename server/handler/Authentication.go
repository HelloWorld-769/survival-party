package handler

import (
	"main/server/request"
	"main/server/response"
	authentication "main/server/services/Authentication"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func SignupHandler(ctx *gin.Context) {
	var input request.SigupRequest
	err := utils.RequestDecoding(ctx, &input)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	//call the service with the inputRequest credentials
	authentication.SignupService(ctx, &input)

}

func LoginHandler(ctx *gin.Context) {

	var input request.LoginRequest
	err := utils.RequestDecoding(ctx, &input)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	authentication.LoginService(ctx, &input)

}
