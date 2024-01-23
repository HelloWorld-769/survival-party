package shop

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func RefillEnergy() {
	query := "select * from user_game_stats where energy < 10"
	var users []model.UserGameStats
	db.QueryExecutor(query, &users)

	// Create a map to store column updates
	updates := make(map[string]interface{})

	// Iterate over users and update their energy in the map
	for _, user := range users {
		updates["energy"] = user.Energy + 1

		query = "update user_game_stats set energy=energy+1 where user_id=?"
		err := db.QueryExecutor(query, nil, user.UserId)

		if err != nil {
			fmt.Println("error updating", err)
			return
		}
	}
}

// EnergyRefillTimer Gives the time for energy renewal
//
// @Summary Get the time left for energy refill
// @Description Get the time left for energy refill
// @Tags Energy
// @Accept json
// @Produce json
// @Success 200 {object} response.Success "Sucess"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /energy-refill-timer [get]
func EnergyRefillTimer(ctx *gin.Context) {

	timeLeft := EnergyTimer()

	response.ShowResponse("RefillEnergy Timer ", utils.HTTP_OK, utils.SUCCESS, timeLeft, ctx)
}

func EnergyTimer() utils.TimeLeft {
	now := time.Now()
	minutes := now.Minute()
	seconds := now.Second()

	// Find how many minutes have passed since the last event
	minutesSinceLastEvent := minutes % 2

	// Calculate how many minutes and seconds are left until the next event
	minutesLeft := 2 - minutesSinceLastEvent - 1
	secondsLeft := 60 - seconds

	// timeleft := fmt.Sprintf("%d minute(s) and %d second(s) until the next event.", minutesLeft, secondsLeft)
	var timeLeft utils.TimeLeft
	timeLeft.Minutes = minutesLeft
	timeLeft.Seconds = secondsLeft

	return timeLeft

}
