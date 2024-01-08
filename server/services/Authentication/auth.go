package authentication

import (
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/services/token"
	"main/server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func SignupService(ctx *gin.Context, input *request.SigupRequest) {

	err := utils.IsPassValid(input.User.Password)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	encryptedPassword, err := utils.HashPassword(input.User.Password)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	userRecord := model.User{
		Email:    input.User.Email,
		Password: *encryptedPassword,
		Username: input.User.Username,
		Avatar:   input.User.Avatar,
	}

	err = db.CreateRecord(&userRecord)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	userGameStats := model.UseGameStats{
		UserId:         userRecord.Id,
		MatchesPlayed:  0,
		MatchesWon:     0,
		TotalTimeSpent: time.Now(),
		TotalKills:     0,
	}

	// expirationTime := time.Now().Add(time.Minute * 5)

	err = db.CreateRecord(&userGameStats)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, "", ctx)
		return
	}

	//creating reset token
	// resetClaims := model.Claims{
	// 	Id: userRecord.Id,
	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		ExpiresAt: jwt.NewNumericDate(expirationTime),
	// 	},
	// }
	// tokenString, err := token.GenerateToken(resetClaims)
	// if err != nil {
	// 	// If there is an error in generating the reset token, return an error response.
	// 	response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
	// 	return
	// }

	// link := ctx.Request.Header.Get("Origin") + "/reset-password?token=" + *tokenString

	response.ShowResponse("A mail has been sent to your email, please verify your account", utils.HTTP_OK, utils.SUCCESS, nil, ctx)

}

func LoginService(ctx *gin.Context, input *request.LoginRequest) {

	var user *model.User
	//check if the user exists or not in the database
	if !(db.RecordExist("users", input.User.Email, "email")) {
		response.ShowResponse(utils.USER_NOT_FOUND, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err := db.FindById(&user, input.User.Email, "email")
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	if !utils.CheckPasswordHash(input.User.Password, user.Password) {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	//check if emial is verified or not
	if !user.Emailverified {
		response.ShowResponse("You have to confirm your email address before continuing.", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	accessTokenExpirationTime := time.Now().Add(48 * time.Hour)
	accessTokenClaims := model.Claims{
		Id: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(accessTokenExpirationTime),
		},
	}

	accessToken, err := token.GenerateToken(accessTokenClaims)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	session := model.Session{
		UserId: user.Id,
		Token:  *accessToken,
	}
	err = db.CreateRecord(&session)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, ctx, nil)
		return
	}

	response.ShowResponse(utils.LOGIN_SUCCESS, utils.HTTP_OK, utils.SUCCESS, struct {
		Token string `json:"token"`
	}{Token: "Bearer " + *accessToken}, ctx)

}

func SignoutService(ctx *gin.Context, userId string) {
	var sessionDetails model.Session
	if !db.RecordExist("sessions", userId, "user_id") {
		response.ShowResponse("Session for current user has already been ended", utils.HTTP_NOT_FOUND, utils.FAILURE, nil, ctx)
		return
	}
	err := db.DeleteRecord(&sessionDetails, userId, "user_id")
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	response.ShowResponse(utils.LOGOUT_SUCCESS, utils.HTTP_OK, utils.SUCCESS, nil, ctx)

}
