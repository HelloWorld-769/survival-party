package settings

import (
	"main/server/db"
	"main/server/model"
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

func UpdateSettingsService(ctx *gin.Context, userId string) {

}
