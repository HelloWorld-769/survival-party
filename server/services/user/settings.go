package user

import (
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func GetSettingsService(ctx *gin.Context, userId string) {

	var userSettings model.UserSettings

	query := "select * from user_settings where user_id = ?"
	err := db.QueryExecutor(query, &userSettings, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("settings fetched successfully", utils.HTTP_OK, utils.SUCCESS, userSettings, ctx)

}

func UpdateSettingsService(ctx *gin.Context, userId string, req request.UpdatePlayerSettingsRequest) {

	var userSettings model.UserSettings

	err := db.FindById(&userSettings, userId, "user_id")
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	if req.Sound != userSettings.Sound && req.Sound != 0 {
		userSettings.Sound = req.Sound
	}
	if req.Music != userSettings.Music && req.Music != 0 {
		userSettings.Music = req.Music
	}
	if req.JoystickSize != userSettings.JoystickSize && req.JoystickSize != 0 {
		userSettings.JoystickSize = req.JoystickSize
	}
	if req.Vibration != userSettings.Vibration {
		userSettings.Vibration = req.Vibration
	}
	if req.Language != userSettings.Language && req.Language != "" {
		userSettings.Language = req.Language
	}
	if req.FriendRequests != userSettings.FriendRequests {
		userSettings.FriendRequests = req.FriendRequests
	}
	if req.Notifications != userSettings.Notifications {
		userSettings.Notifications = req.Notifications
	}
	if req.VoicePack != userSettings.VoicePack {
		userSettings.VoicePack = req.VoicePack
	}

	updateFields := map[string]interface{}{
		"sound":           userSettings.Sound,
		"music":           userSettings.Music,
		"joystick_size":   userSettings.JoystickSize,
		"vibration":       userSettings.Vibration,
		"language":        userSettings.Language,
		"friend_requests": userSettings.FriendRequests,
		"notifications":   userSettings.Notifications,
		"voice_pack":      userSettings.VoicePack,
	}
	err = db.UpdateZeroVals(model.UserSettings{}, "user_id", userId, updateFields)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("settings Updated successfully", utils.HTTP_OK, utils.SUCCESS, userSettings, ctx)

}
