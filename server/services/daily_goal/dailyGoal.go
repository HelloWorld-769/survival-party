package dailygoal

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/utils"
	"math/rand"
	"time"
)

func DailyGoalGeneration() {
	noOfGoalsMin := 4
	noOfGoalsMax := 7
	noOfGoals := rand.Intn(noOfGoalsMax-noOfGoalsMin+1) + noOfGoalsMin

	var data []struct {
		Id    string
		level int64
	}
	query := "SELECT id,level FROM users WHERE emailverified =true"
	err := db.QueryExecutor(query, &data)
	if err != nil {
		fmt.Println("Error in gettign the users from the database")
		return
	}

	for _, it := range data {

		var temp []model.UserDailyGoals

		for i := 0; i < noOfGoals; i++ {
			rand.Seed(time.Now().UnixNano())
			//selecting a random goal type
			randGoalType := rand.Intn(6) + 1

			fmt.Println("aRandom Goal type:", randGoalType)

			var record model.UserDailyGoals
			record.UserId = it.Id
			record.GoalType = int64(randGoalType)

			//selecting the random currency type
			currencyType := []int{utils.C_ADS, utils.C_GEMS}
			randCurrency := currencyType[rand.Intn(len(currencyType))]
			record.CurrencyType = int64(randCurrency)
			if randCurrency == utils.C_ADS {
				record.Price = 1
			} else {
				record.Price = int64(rand.Intn(100-50+1) + 50)
			}

			if randGoalType == int(utils.PLAYERS_KILLED) || randGoalType == int(utils.ZOMBIES_KILLED) {
				min := 20
				max := 80

				lowerRange := min + (int(it.level-1) * ((max - min) / utils.TOTAL_LEVELS))
				upperRange := min + (int(it.level) * ((max - min) / utils.TOTAL_LEVELS))

				kills := generateRandomNumber(int(it.level), lowerRange, upperRange)

				baseCoins := 20
				baseGems := 2
				record.Coins = int64(baseCoins) * kills
				record.Gems = int64(baseGems) * kills
				record.TotalProgress = kills

			} else if randGoalType == int(utils.MINI_GAMES_PLAYED) {
				min := 3
				max := 5
				baseCoins := 25
				baseGems := 8

				gamPlay := generateRandomNumber(int(it.level), min, max)
				record.Coins = int64(baseCoins) * gamPlay
				record.Gems = int64(baseGems) * gamPlay

			} else if randGoalType == int(utils.BECAME_ZOMBIE) {
				min := 2
				max := 6

				baseCoins := 20
				baseGems := 2

				gamPlay := generateRandomNumber(int(it.level), min, max)
				record.Coins = int64(baseCoins) * gamPlay
				record.Gems = int64(baseGems) * gamPlay

			} else if randGoalType == int(utils.ESCAPE_SURVIVOR) {

				min := 3
				max := 10
				baseCoins := 50
				baseGems := 15

				gamPlay := generateRandomNumber(int(it.level), min, max)
				record.Coins = int64(baseCoins) * gamPlay
				record.Gems = int64(baseGems) * gamPlay
			} else if randGoalType == int(utils.COMPLETED_TASKS) {
				min := 2
				max := 3
				baseCoins := 250
				baseGems := 25

				gamPlay := generateRandomNumber(int(it.level), min, max)
				record.Coins = int64(baseCoins) * gamPlay
				record.Gems = int64(baseGems) * gamPlay
			}

			temp = append(temp, record)
		}

		err = db.CreateRecord(&temp)
		if err != nil {
			fmt.Println("Error in creting the entry in db.")
			return
		}
	}

	fmt.Println("Sucessfully generated daily goals for all the users")

}

func DeleteAllGoals() {
	query := "DELETE FROM user_daily_goals"
	err := db.RawExecutor(query)
	if err != nil {
		fmt.Println("Error in deleting the records from the tabale")
		return
	}
}

func generateRandomNumber(seed, min, max int) int64 {
	rand.Seed(int64(seed)) // Seed the random number generator

	// Generate a random number within the desired range
	randomValue := rand.Intn(max-min+1) + min

	return int64(randomValue)
}
