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

	query := "select * from user_game_stats where energy <10 "
	var users []model.UserGameStats
	db.QueryExecutor(query, &users)
	for _, user := range users {

		user.Energy++
		err := db.UpdateRecord(&user, user.Energy, "energy").Error
		if err != nil {
			fmt.Println("error updating", err)
			return
		}
	}
}

func EnergyRefillTimer(ctx *gin.Context) {

	timeLeft := EnergyTimer()

	response.ShowResponse("RefillEnergy Hit ", utils.HTTP_OK, utils.SUCCESS, timeLeft, ctx)
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

	// Convert to a human-readable format
	// timeleft := fmt.Sprintf("%d minute(s) and %d second(s) until the next event.", minutesLeft, secondsLeft)
	var timeLeft utils.TimeLeft
	timeLeft.Minutes = minutesLeft
	timeLeft.Seconds = secondsLeft

	return timeLeft

}
