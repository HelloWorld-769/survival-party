package user

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	var dbResposne []struct {
		Username        string    `json:"username"`
		Avatar          int64     `json:"avatar"`
		UserId          string    `json:"userId"`
		XP              int64     `json:"xp"`
		Level           int64     `json:"level"`
		Energy          int64     `json:"energy"`
		TotalCoins      int64     `json:"totalCoins"`
		CurrentCoins    int64     `json:"currentCoins"`
		TotalGems       int64     `json:"totalGems"`
		CurrentGems     int64     `json:"currentGems"`
		CurrentTrophies int64     `json:"currentTrophies"`
		HighestTrophies int64     `json:"highestTrophies"`
		MatchesWon      int64     `json:"matchesWon"`
		MatchesLost     int64     `json:"matchesLost"`
		TotalTimeSpent  int64     `json:"timeSpent"`
		TotalKills      int64     `json:"totalKills"`
		Badge           int64     `json:"badge"`
		CreatedAt       time.Time `json:"created_at"`
		UpdatedAt       time.Time `json:"updated_at"`
		DeletedAt       gorm.DeletedAt
	}

	// var playerResponse model.UserGameStats
	query := `
		SELECT u.username,u.avatar, ugs.*,ub.badge
		FROM users u 
		JOIN user_game_stats ugs ON ugs.user_id=u.id
		LEFT JOIN user_badges ub ON ub.user_id = u.id
		WHERE u.id=?`

	err := db.QueryExecutor(query, &dbResposne, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	type resp struct {
		Username        string  `json:"username"`
		Avatar          int64   `json:"avatar"`
		UserId          string  `json:"userId"`
		XP              int64   `json:"xp"`
		Level           int64   `json:"level"`
		Energy          int64   `json:"energy"`
		TotalCoins      int64   `json:"totalCoins"`
		CurrentCoins    int64   `json:"currentCoins"`
		TotalGems       int64   `json:"totalGems"`
		CurrentGems     int64   `json:"currentGems"`
		CurrentTrophies int64   `json:"currentTrophies"`
		HighestTrophies int64   `json:"highestTrophies"`
		MatchesPlayed   int64   `json:"matchesPlayed"`
		MatchesWon      int64   `json:"matchesWon"`
		MatchesLost     int64   `json:"matchesLost"`
		TotalTimeSpent  int64   `json:"timeSpent"`
		TotalKills      int64   `json:"totalKills"`
		Badges          []int64 `json:"badges"`
	}

	var playerResponse resp
	for _, data := range dbResposne {
		playerResponse = resp{
			Username:        data.Username,
			Avatar:          data.Avatar,
			UserId:          data.UserId,
			XP:              data.XP,
			Level:           data.Level,
			Energy:          data.Energy,
			TotalCoins:      data.TotalCoins,
			CurrentCoins:    data.CurrentCoins,
			TotalGems:       data.TotalGems,
			CurrentGems:     data.CurrentGems,
			CurrentTrophies: data.CurrentTrophies,
			HighestTrophies: data.HighestTrophies,
			MatchesPlayed:   data.MatchesWon + data.MatchesLost,
			MatchesWon:      data.MatchesWon,
			MatchesLost:     data.MatchesLost,
			TotalTimeSpent:  data.TotalTimeSpent,
			TotalKills:      data.TotalKills,
		}
		if data.Badge != 0 {
			playerResponse.Badges = append(playerResponse.Badges, data.Badge)
		} else {
			playerResponse.Badges = []int64{}
		}
	}

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, playerResponse, ctx)

}

func UpdateDayCount() {

	var users []model.User
	query := "select * from users where email_verified =true"
	err := db.QueryExecutor(query, &users)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}

	for _, user := range users {

		user.DayCount++
		fmt.Println("daycount updated", user.DayCount)
		err := db.UpdateRecord(&user, user.Id, "id").Error
		if err != nil {
			fmt.Println("error:", err.Error())
			return
		}
	}

}
