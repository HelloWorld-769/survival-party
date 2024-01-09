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

	response.ShowResponse(utils.PLAYER_UPDATE_SUCCESS, utils.HTTP_OK, utils.SUCCESS, input, ctx)
}

func GetPlayerStatsService(ctx *gin.Context, userId string) {

	// var playerResponse struct {
	// 	Username        string    `json:"username"`
	// 	Avatar          int64     `json:"avatar"`
	// 	UserId          string    `json:"userId"`
	// 	XP              int64     `json:"xp"`
	// 	Level           int64     `json:"level"`
	// 	Coins           int64     `json:"coins"`
	// 	Badges          string    `json:"bage"`
	// 	Gems            int64     `json:"gems"`
	// 	Energy          int64     `json:"energy"`
	// 	TotalCoins      int64     `json:"totalCoins"`
	// 	CurrentCoins    int64     `json:"currentCoins"`
	// 	TotalGems       int64     `json:"totalGems"`
	// 	CurrentGems     int64     `json:"currentGems"`
	// 	CurrentTrophies int64     `json:"currentTrophies"`
	// 	HighestTrophies int64     `json:"highestTrophies"`
	// 	MatchesPlayed   int64     `json:"matchesPlayed"`
	// 	MatchesWon      int64     `json:"matchesWon"`
	// 	TotalTimeSpent  time.Time `json:"timeSpent"`
	// 	TotalKills      int64     `json:"totalKills"`
	// 	CreatedAt       time.Time `json:"created_at"`
	// 	UpdatedAt       time.Time `json:"updated_at"`
	// 	DeletedAt       gorm.DeletedAt
	// }

	var playerResponse model.UserGameStats
	query := `
	SELECT u.username, u.avatar, ARRAY_AGG(ugs.badges) AS badges
	FROM users u
	JOIN user_game_stats ugs ON ugs.user_id = u.id
	WHERE u.id = ?
	GROUP BY u.id, u.username, u.avatar;`

	err := db.QueryExecutor(query, &playerResponse, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, playerResponse, ctx)

}
