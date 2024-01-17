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

// @Summary Sign Up
// @Description Perform signup and sends email for verification
// @Tags Authentication
// @Accept json
// @Produce json
// @Param guestLoginRequest body request.SigupRequest true "Signup Request"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /users/sign-up [post]
func SignupHandler(ctx *gin.Context) {
	var input request.SigupRequest
	err := utils.RequestDecoding(ctx, &input)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = input.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	//call the service with the inputRequest credentials
	authentication.SignupService(ctx, &input)

}

// LoginService handles user login and token generation.
//
// @Summary User Login
// @Description Perform Users login and generate access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginDetails body request.LoginRequest true "Login Details"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /users/sign-in [post]
func LoginHandler(ctx *gin.Context) {

	var input request.LoginRequest
	err := utils.RequestDecoding(ctx, &input)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = input.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	authentication.LoginService(ctx, &input)

}

// SendOtpHanlder sends the otp on the register email
//
// @Summary Sends OTP
// @Description Sends the otp on the register email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginDetails body request.EmailRequest true "Email Details"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /send-otp [post]
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

// ResetPasswordHandler Resets the password of the user
//
// @Summary Resets the password
// @Description Resets the password of the user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginDetails body request.RestPasswordRequest true "Email Details"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /reset-password [post]
func ResetPasswordHandler(ctx *gin.Context) {

	utils.SetHeader(ctx)

	var req request.RestPasswordRequest
	fmt.Println("request", ctx.Request.Body)

	utils.RequestDecoding(ctx, &req)
	fmt.Println("req", req)
	err := utils.IsPassValid(req.Password)
	if err != nil {
		response.ShowResponse(err.Error(), 400, "Failure", "", ctx)
		return
	}

	//validation Check on request body fields
	err = validation.CheckValidation(&req)
	if err != nil {
		response.ShowResponse(err.Error(), 400, "Failure", "", ctx)
		return
	}

	//call the service
	authentication.ResetPasswordService(ctx, req)
}

// CheckOtpHandler Verifies the otp sent on email
//
// @Summary Verifies OTP
// @Description Verifies the otp sent on email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginDetails body request.OtpRequest true "Email Details"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /check_otp [post]
func CheckOtpHandler(ctx *gin.Context) {

	utils.SetHeader(ctx)

	var req request.OtpRequest
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
	authentication.CheckOtpService(ctx, req)

}

// @Summary Logout Player
// @Description	Logs out a player
// @Accept			json
// @Produce		json
// @Param Authorization header string true "Player Access Token"
// @Success		200	{object}	response.Success
// @Failure		400	{object}	response.Success
// @Failure		401	{object}	response.Success
// @Tags			Authentication
// @Router			/users/sign_out [delete]
func SignoutHandler(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	authentication.SignoutService(ctx, userId.(string))

}

// Verifies the email of the user
//
// @Summary User email verification
// @Description Perform Email verifictaion
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /users/email-verify [get]
func VerifyEmailHandler(ctx *gin.Context) {

	var tokenString string
	if ctx.Query("token") != "" {
		tokenString = ctx.Query("token")
	}

	authentication.VerifyEmail(ctx, tokenString)

}

// Social handles user social login and token generation.
//
// @Summary User Login
// @Description Perform Users social login and generate access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginDetails body request.SocialLoginReq true "Login Details"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /users/social_login [post]
func SocialLoginHandler(ctx *gin.Context) {
	var input request.SocialLoginReq
	err := utils.RequestDecoding(ctx, &input)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return

	}

	err = input.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	authentication.SocialLoginService(ctx, &input)

}

func ChangePasswordHandler(ctx *gin.Context) {

	var inputPassword string
	err := utils.RequestDecoding(ctx, &inputPassword)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	//validation check
	if inputPassword == "" {

		response.ShowResponse("password cannot be blank", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	userId, exists := ctx.Get("user_id")
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	authentication.ChangePasswordService(ctx, userId.(string), inputPassword)
}
