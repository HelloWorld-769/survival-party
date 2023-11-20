package handler

import (
	"fmt"
	"main/server/request"
	"main/server/response"
	"main/server/services/alert_service/twilio"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

func TwilioServiceHnadler(ctx *gin.Context) {

	utils.SetHeader(ctx)

	var req request.TwilioSmsRequest

	utils.RequestDecoding(ctx, &req)
	fmt.Println("req", req)

	//validation Check on request body fields
	err := validation.CheckValidation(&req)
	if err != nil {
		response.ShowResponse(err.Error(), 400, "Failure", "", ctx)
		return
	}
	twilio.SendOtpService(ctx, req)

}
