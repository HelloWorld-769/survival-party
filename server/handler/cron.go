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

		currentTime := time.Now().UTC()
		fmt.Println("Current time is:", currentTime)
		// Format the time to HH:MM:SS
		formattedTime := currentTime.Format("15:04")
		fmt.Println("formatted time is:", formattedTime)

		if formattedTime == "23:59" {

			fmt.Println("Cron is working")

			user.UpdateDayCount()

			rewards.UpdateDailyRewardsData()
			//create User daily rewards (available to claim)
			rewards.CreateUserDailyReward()
			shop.GiveNewSpecialOffer()

			//Daily goal generation
			dailygoal.DeleteAllGoals()
			dailygoal.DailyGoalGeneration(false, nil)
		}
		if isEvenMinutes(formattedTime) {

			shop.RefillEnergy()
		}

		// dailygoal.DailyGoalGeneration()

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
