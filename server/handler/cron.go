package handler

import (
	"fmt"
	dailygoal "main/server/services/daily_goal"
	"main/server/services/rewards"
	"main/server/services/shop"
	"main/server/services/user"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func StartCron() {
	c := cron.New()
	fmt.Println("Current time is:", time.Now())
	//Concept
	//Run cron at every minute--[change it to hour  in future]
	//check the next reward time is same as current time. If yes give the rewards
	c.AddFunc("*/1 * * * *", func() {
		fmt.Println(".....................Cron hit..................................")

		currentTime := time.Now()
		// Format the time to HH:MM:SS
		formattedTime := currentTime.Format("15:04")
		fmt.Println("formatted time is:", formattedTime)

		if formattedTime == "05:17" {

			user.UpdateDayCount()

			rewards.UpdateDailyRewardsData()
			//create User daily rewards (available to claim)
			rewards.CreateUserDailyReward()
			shop.GiveNewSpecialOffer()
		}

		if formattedTime == "14:20" {

			dailygoal.DailyGoalGeneration(false, nil)

		}
		if isEvenMinutes(formattedTime) {

			shop.RefillEnergy()
		}

		// dailygoal.DailyGoalGeneration()

		// if formattedTime == "17:48" {

		// 	dailygoal.DeleteAllGoals()
		// }

	})

	c.Start()

}

func isEvenMinutes(timeString string) bool {
	// Split the time string into hours and minutes
	parts := strings.Split(timeString, ":")
	if len(parts) != 2 {
		// Invalid time format
		return false
	}

	// Extract the minutes part
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		// Invalid minutes
		return false
	}

	// Check if minutes is even
	return minutes%2 == 0
}
