package authentication

import (
	"errors"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/services/alert_service/Gomail"
	"main/server/services/token"
	"main/server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
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

	userGameStats := model.UserGameStats{
		UserId:         userRecord.Id,
		MatchesPlayed:  0,
		MatchesWon:     0,
		TotalTimeSpent: time.Now(),
		TotalKills:     0,
	}

	userSettings := model.UserSettings{
		UserId:         userRecord.Id,
		Sound:          1,
		Music:          1,
		Vibration:      false,
		VoicePack:      false,
		Notifications:  false,
		FriendRequests: false,
		Language:       "english",
	}

	err = db.CreateRecord(&userSettings)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	err = db.CreateRecord(&userGameStats)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//sending verification emial to the user
	resetClaims := model.Claims{
		Id: userRecord.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	tokenString, err := token.GenerateToken(resetClaims)
	if err != nil {
		// If there is an error in generating the reset token, return an error response.
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	link := ctx.Request.Header.Get("Origin") + "/email-verify?token=" + *tokenString

	Gomail.SendEmailService(ctx, link, userRecord.Email)

	response.ShowResponse("A mail has been sent to your email, please verify your account", utils.HTTP_OK, utils.SUCCESS, nil, ctx)

}

func VerifyEmail(ctx *gin.Context, userId string) {
	//check if the email is already verifed or not
	var emailStatus bool
	query := "SELECT emailverified FROM users WHERE id=?"

	err := db.QueryExecutor(query, &emailStatus, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	if emailStatus {
		response.ShowResponse("Email already verified", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)

		return
	}

}

func LoginService(ctx *gin.Context, input *request.LoginRequest) {

	var user *model.User

	//Login using username and email
	if utils.IsEmail(input.User.Email) {
		input.User.Email = strings.ToLower(input.User.Email)
		err := db.FindById(&user, input.User.Email, "email")
		if err != nil {
			// If the player doesn't exist, return an error response.
			response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}

	} else {
		err := db.FindById(&user, input.User.Email, "username")
		if err != nil {
			// If the player doesn't exist, return an error response.
			response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}
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
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
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

func SocialLoginService(ctx *gin.Context, input *request.SocialLoginReq) {

	var accessToken *string
	//if there is no entry in db then user is doing signup with social login
	if !db.RecordExist("users", input.Email, "email") {
		var count int
		query := "SELECT count(*) FROM users"
		err := db.QueryExecutor(query, &count)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
		//give a random userNmae to that user
		userRecord := model.User{
			Email:         input.Email,
			Emailverified: false,
			Password:      "",
			Username:      "Suvival_Party_" + strconv.Itoa(count),
			Avatar:        input.Avatar,
			SocialId:      input.Uid,
		}

		err = db.CreateRecord(&userRecord)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		userSettings := model.UserSettings{
			UserId:         userRecord.Id,
			Sound:          1,
			Music:          1,
			Vibration:      false,
			VoicePack:      false,
			Notifications:  false,
			FriendRequests: false,
			Language:       "english",
		}

		err = db.CreateRecord(&userSettings)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		userGameStats := model.UserGameStats{
			UserId:         userRecord.Id,
			MatchesPlayed:  0,
			MatchesWon:     0,
			TotalTimeSpent: time.Now(),
			TotalKills:     0,
		}

		err = db.CreateRecord(&userGameStats)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		accessTokenExpirationTime := time.Now().Add(48 * time.Hour)
		accessTokenClaims := model.Claims{
			Id: userRecord.Id,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(accessTokenExpirationTime),
			},
		}

		accessToken, err = token.GenerateToken(accessTokenClaims)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}

		session := model.Session{
			UserId: userRecord.Id,
			Token:  *accessToken,
		}
		err = db.CreateRecord(&session)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

	} else {
		//user is trying to log in in using social login
		var user *model.User
		query := "SELECT * FROM users WHERE email=? AND social_id=?"
		err := db.QueryExecutor(query, &user, input.Email, input.Uid)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ShowResponse("User not found", utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
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

		accessToken, err = token.GenerateToken(accessTokenClaims)
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
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

	}

	response.ShowResponse(utils.LOGIN_SUCCESS, utils.HTTP_OK, utils.SUCCESS, struct {
		Token string `json:"token"`
	}{Token: "Bearer " + *accessToken}, ctx)

}
