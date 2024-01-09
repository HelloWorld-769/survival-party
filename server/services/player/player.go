package player

import (
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func UpdatePlayerService(ctx *gin.Context, userId string, input request.UpdatePlayer) {

	var user model.User
	query := "SELECT * from users WHERE id=?"
	err := db.QueryExecutor(query, &user, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Avatar != 0 {
		user.Avatar = input.Avatar
	}

	//update the record
	err = db.UpdateRecord(&user, userId, "id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("User data updated successfully", utils.HTTP_OK, utils.SUCCESS, input, ctx)
}

func GetPlayerStatsService(ctx *gin.Context, userId string) {

}
