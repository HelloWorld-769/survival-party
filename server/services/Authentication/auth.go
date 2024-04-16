package authentication

import (
	"errors"
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/services/alert_service/Gomail"
	dailygoal "main/server/services/daily_goal"
	"main/server/services/rewards"
	"main/server/services/token"
	"main/server/utils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func SignupService(ctx *gin.Context, input request.SigupRequest) {

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
		Email:             input.User.Email,
		Password:          *encryptedPassword,
		Username:          strings.ToLower(input.User.Username),
		Avatar:            input.User.Avatar,
		EmailVerifiedAt:   time.Now(),
		UsernameUpdatedAt: time.Now(),
		DayCount:          1,
	}

	tx := db.BeginTransaction()
	err = tx.Create(&userRecord).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			response.ShowResponse("Credentials should be unique", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	userGameStats := model.UserGameStats{
		UserId:         userRecord.Id,
		CurrentCoins:   10000,
		Energy:         10,
		CurrentGems:    10000,
		TotalTimeSpent: 0,
		TotalKills:     0,
	}

	userSettings := model.UserSettings{
		UserId:         userRecord.Id,
		Sound:          0.5,
		Music:          0.5,
		JoystickSize:   0.5,
		Vibration:      false,
		VoicePack:      false,
		Notifications:  false,
		FriendRequests: false,
		Language:       "english",
	}

	err = tx.Create(&userSettings).Error
	if err != nil {
		tx.Rollback()
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	err = tx.Create(&userGameStats).Error
	if err != nil {
		tx.Rollback()
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	var specailOfferId string
	query := "SELECT id FROM special_offers order by created_at ASC limit 1"
	err = db.QueryExecutor(query, &specailOfferId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//Giving the starter pack to the user after signup
	//For 7 days starter pack will be valid
	userStartPack := model.UserSpecialOffer{
		SpecialOfferId: specailOfferId,
		UserId:         userRecord.Id,
		Purchased:      false,
	}

	err = tx.Create(&userStartPack).Error
	if err != nil {
		tx.Rollback()
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//sending verification email to the user
	resetClaims := model.Claims{
		Id: userRecord.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	tokenString, err := token.GenerateToken(resetClaims)
	if err != nil {
		tx.Rollback()
		// If there is an error in generating the reset token, return an error response.
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	// fmt.Println("hvasdas", ctx.Request.Header.Get("Origin"))
	link := "http://192.180.2.109:" + os.Getenv("PORT") + "/api/v1/users/email-verify?token=" + tokenString

	fmt.Println("link is", link)

	err = Gomail.SendEmailService(link, userRecord.Email)
	if err != nil {
		tx.Rollback()
		// If there is an error in generating the reset token, return an error response.
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("A mail has been sent to your email, please verify your account", utils.HTTP_OK, utils.SUCCESS, nil, ctx)

}

func VerifyEmail(ctx *gin.Context, tokenString string) {

	//Decoding the token
	claims, err := token.DecodeToken(tokenString)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	fmt.Println("claims:", claims)
	err = claims.Valid()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	//check if the email is already verifed or not
	var emailStatus bool
	query := "SELECT email_verified FROM users WHERE id=?"

	err = db.QueryExecutor(query, &emailStatus, claims.Id)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	if emailStatus {
		response.ShowResponse("Email already verified", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	query = "UPDATE users SET email_verified=true WHERE id=?"
	err = db.RawExecutor(query, claims.Id)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	var specailOfferId string
	query = "SELECT id FROM special_offers order by created_at ASC limit 1"
	err = db.QueryExecutor(query, &specailOfferId)
	if err != nil {
		return
	}

	//Giving the starter pack to the user after signup
	//For 7 days starter pack will be valid
	userStartPack := model.UserSpecialOffer{
		SpecialOfferId: specailOfferId,
		UserId:         claims.Id,
		Purchased:      false,
	}

	err = db.CreateRecord(&userStartPack)
	if err != nil {
		return
	}

	err = rewards.CreateStarterDailyRewards(claims.Id)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//generating daily goals for user
	go dailygoal.DailyGoalGeneration(true, &claims.Id)
	//generating level reward for user
	go rewards.GenerateLevelReward(claims.Id)

	response.ShowResponse("Email verified succesfully", utils.HTTP_OK, utils.SUCCESS, nil, ctx)

}

func LoginService(ctx *gin.Context, input *request.LoginRequest) {

	var user *model.User

	input.User.Credential = strings.ToLower(input.User.Credential)
	//Login using username and email
	if utils.IsEmail(input.User.Credential) {
		err := db.FindById(&user, input.User.Credential, "email")
		if err != nil {
			// If the player doesn't exist, return an error response.
			response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}

	} else {
		err := db.FindById(&user, input.User.Credential, "username")
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
	if !user.EmailVerified {
		response.ShowResponse("You have to confirm your email address before continuing.", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	if user.DayCount == 0 {

		user.DayCount = 1
		err := db.UpdateRecord(&user, user.Id, "id").Error
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
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

	if db.RecordExist("sessions", user.Id, "user_id") {
		//update the record
		query := "UPDATE sessions SET token=? WHERE user_id=?"
		err := db.RawExecutor(query, accessToken, user.Id)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

	} else {
		session := model.Session{
			UserId: user.Id,
			Token:  accessToken,
		}
		err = db.CreateRecord(&session)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
	}

	response.ShowResponse(utils.LOGIN_SUCCESS, utils.HTTP_OK, utils.SUCCESS, struct {
		Token  string `json:"token"`
		UserId string `json:"userId"`
	}{Token: "Bearer " + accessToken,
		UserId: user.Id,
	}, ctx)

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

	var userId string
	var accessToken string
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
			Email:             input.Email,
			EmailVerified:     true,
			Password:          "",
			Username:          strings.ToLower("Suvival_Party_" + strconv.Itoa(count)),
			Avatar:            input.Avatar,
			SocialId:          input.Uid,
			EmailVerifiedAt:   time.Now(),
			UsernameUpdatedAt: time.Now(),
			DayCount:          1,
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
		userId = userRecord.Id

		userGameStats := model.UserGameStats{
			UserId:         userRecord.Id,
			TotalTimeSpent: 0,
			Energy:         10,
			CurrentCoins:   10000,
			CurrentGems:    10000,
			// Badges:         []int64{},
			TotalKills: 0,
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
			Token:  accessToken,
		}
		err = db.CreateRecord(&session)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		var specailOfferId string
		query = "SELECT id FROM special_offers order by created_at ASC limit 1"
		err = db.QueryExecutor(query, &specailOfferId)
		if err != nil {
			return
		}

		//Giving the starter pack to the user after signup
		//For 7 days starter pack will be valid
		userStartPack := model.UserSpecialOffer{
			SpecialOfferId: specailOfferId,
			UserId:         userRecord.Id,
			Purchased:      false,
		}

		err = db.CreateRecord(&userStartPack)
		if err != nil {
			return
		}

		err = rewards.CreateStarterDailyRewards(userRecord.Id)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
		//generating daily goals for user
		go dailygoal.DailyGoalGeneration(true, &userRecord.Id)
		//generating level reward for user
		go rewards.GenerateLevelReward(userRecord.Id)

	} else {
		//user is trying to log in in using social login
		var user *model.User
		if input.Email != "" {
			query := "SELECT * FROM users WHERE email=? "
			err := db.QueryExecutor(query, &user, input.Email)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					response.ShowResponse("User not found", utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}
				response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
				return
			}
		} else if input.Uid != "" {
			query := "SELECT * FROM users WHERE  social_id=?"
			err := db.QueryExecutor(query, &user, input.Uid)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					response.ShowResponse("User not found", utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}
				response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
				return
			}
		}
		userId = user.Id

		if user.DayCount == 0 {
			user.DayCount = 1
			err := db.UpdateRecord(&user, user.Id, "id").Error
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
				return
			}
		}

		accessTokenExpirationTime := time.Now().Add(48 * time.Hour)
		accessTokenClaims := model.Claims{
			Id: user.Id,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(accessTokenExpirationTime),
			},
		}

		fmt.Println("accessTokenClaims", accessTokenClaims)
		// accessToken = "asdbjasbd"
		var err error
		accessToken, err = token.GenerateToken(accessTokenClaims)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}

		fmt.Println("accessToken", accessToken)
		if db.RecordExist("sessions", user.Id, "user_id") {
			//update the record
			query := "UPDATE sessions SET token=? WHERE user_id=?"
			err := db.RawExecutor(query, accessToken, user.Id)
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
				return
			}

		} else {
			session := model.Session{
				UserId: user.Id,
				Token:  accessToken,
			}
			err = db.CreateRecord(&session)
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
				return
			}
		}

		if user.DayCount == 0 {

			user.DayCount = 1
			err := db.UpdateRecord(&user, user.Id, "id").Error
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
				return
			}
		}

	}

	fmt.Println("accessToken", accessToken)
	response.ShowResponse(utils.LOGIN_SUCCESS, utils.HTTP_OK, utils.SUCCESS, struct {
		Token  string `json:"token"`
		UserId string `json:"userId"`
	}{Token: "Bearer " + accessToken,
		UserId: userId,
	}, ctx)
}
func CheckOtpService(ctx *gin.Context, req request.OtpRequest) {

	//check otp from restSession table corresponding to user email
	var usersRestSession model.ResetSession
	query := "select * from reset_sessions where user_email=?"
	err := db.QueryExecutor(query, &usersRestSession, req.Email)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	if usersRestSession.Otp == req.Otp {
		response.ShowResponse("OTP correct", utils.HTTP_OK, utils.SUCCESS, nil, ctx)
		return
	}

	response.ShowResponse("OTP Incorrect", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)

}

func ResetPasswordService(ctx *gin.Context, req request.RestPasswordRequest) {

	//hash the pasword before updating the password in the table

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	var user model.User
	user.Password = *passwordHash
	err = db.UpdateRecord(&user, req.Email, "email").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("Password Updated Successfully", utils.HTTP_OK, utils.SUCCESS, nil, ctx)

	// query := "UPDATE users SET password =? WHERE email = ?;"
	// db.QueryExecutor(query, &user, passwordHash, req.Email)

}

func ChangePasswordService(ctx *gin.Context, userId string, password string) {

	//for the logged in user
	//check the current password matches with the inputPssword

	var user model.User
	err := db.FindById(&user, userId, "user_id")
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	if utils.CheckPasswordHash(password, user.Password) {

		//update or change the password
		newHashPassword, err := utils.HashPassword(password)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return

		}
		user.Password = *newHashPassword

		err = db.UpdateRecord(&user, userId, "user_id").Error
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		response.ShowResponse("password changed successfully", utils.HTTP_OK, utils.SUCCESS, nil, ctx)
	}

}
