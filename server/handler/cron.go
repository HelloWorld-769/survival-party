package handler

import (
	"fmt"
	"main/server/services/rewards"
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
		if formattedTime == "18:25" {

			//create User daily rewards (available to claim)
			rewards.CreateUserDailyReward()
		}

	})

	c.Start()

}
