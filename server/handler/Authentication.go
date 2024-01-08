package handler

import (
	"fmt"
	"main/server/request"
	"main/server/response"
	authentication "main/server/services/Authentication"
	"main/server/services/alert_service/Gomail"
	"main/server/utils"
	"main/server/validation"

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

func SendOtpHandler(ctx *gin.Context) {

	utils.SetHeader(ctx)

	var req request.EmailRequest
	fmt.Println("request", ctx.Request.Body)

	utils.RequestDecoding(ctx, &req)
	fmt.Println("req", req)

	//validation Check on request body fields
	err := validation.CheckValidation(&req)
	if err != nil {
		response.ShowResponse(err.Error(), 400, "Failure", "", ctx)
		return
	}

	//call the service

	Gomail.SendEmailOtpService(ctx, req)

}

func ResetPasswordHandler(ctx *gin.Context) {

	utils.SetHeader(ctx)

	var req request.RestPasswordRequest
	fmt.Println("request", ctx.Request.Body)

	utils.RequestDecoding(ctx, &req)
	fmt.Println("req", req)

	//validation Check on request body fields
	err := validation.CheckValidation(&req)
	if err != nil {
		response.ShowResponse(err.Error(), 400, "Failure", "", ctx)
		return
	}

	//call the service
	// Gomail.SendEmailOtpService(ctx, req)

}

func SignoutHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	authentication.SignoutService(ctx, userId.(string))

}
