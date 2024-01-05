package authentication

import (
	"fmt"
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

func AlreadyExists(data string) bool {

	return db.RecordExist("users", data, "email")

}

func SignupService(ctx *gin.Context, input *request.SigupRequest) {

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
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, "", ctx)
		return
	}

	userGameStats := model.UseGameStats{
		UserId:         userRecord.Id,
		MatchesPlayed:  0,
		MatchesWon:     0,
		TotalTimeSpent: time.Now(),
		TotalKills:     0,
	}

	err = db.CreateRecord(&userGameStats)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, "", ctx)
		return
	}

	response.ShowResponse(utils.SIGNUP_SUCCESS, utils.HTTP_OK, utils.SUCCESS, "", ctx)

}

func LoginService(ctx *gin.Context, input *request.LoginRequest) {

	var user *model.User
	//check if the user exists or not in the database
	if !(db.RecordExist("users", input.User.Email, "email")) {
		response.ShowResponse(utils.USER_NOT_FOUND, utils.HTTP_BAD_REQUEST, utils.FAILURE, "", ctx)
		return
	}

	err := db.FindById(&user, input.User.Email, "email")
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	if !utils.CheckPasswordHash(input.User.Password, user.Password) {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, "", ctx)
		return
	}

	accessTokenExpirationTime := time.Now().Add(48 * time.Hour)
	accessTokenClaims := model.Claims{
		Id:   user.Id,
		Role: "player",
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
	fmt.Println("accessToken", accessToken)

}
