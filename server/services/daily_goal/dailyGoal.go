package dailygoal

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DailyGoalGeneration(isNew bool, userId *string) error {
	noOfGoalsMin := 3
	noOfGoalsMax := 5
	rand.Seed(time.Now().UnixNano())
	noOfGoals := rand.Intn(noOfGoalsMax-noOfGoalsMin+1) + noOfGoalsMin

	fmt.Println("No.OfGoals", noOfGoals)

	var data []struct {
		Id    string
		Level int64
	}
	if isNew {
		query := "SELECT id,level FROM users WHERE id=?"
		err := db.QueryExecutor(query, &data, *userId)

		if err != nil {
			fmt.Println("Error in getting the users from the database")
			return err
		}

	} else {
		query := "SELECT id,level FROM users WHERE email_verified =true"
		err := db.QueryExecutor(query, &data)

		if err != nil {
			fmt.Println("Error in getting the users from the database")
			return err
		}
	}

	fmt.Println("Data is ", data)

	// Iterate over each user in the data slice
	for _, it := range data {

		// Initialize a slice to store the generated UserDailyGoals records
		var temp []model.UserDailyGoals

		mp := make(map[int]int)

		//generatig a unique number to give to all rewrds so that they can be linked to major rwrd which is given on completion of all goals
		rewardId := uuid.New().String()

		//this code is to avoid random number genrator to generate same number twice
		usedValues := make(map[int]bool)
		for i := 0; i < noOfGoals; i++ {
			rand.Seed(time.Now().UnixNano())

			// Generate a unique random goal type
			var randGoalType int
			for {
				randGoalType = rand.Intn(5) + 1
				if !usedValues[randGoalType] {
					usedValues[randGoalType] = true
					break
				}
			}

			// Add the random goal type to the map
			mp[i] = randGoalType
		}
		fmt.Println("Map is ", mp)

		// Initialize a variable to store the sum of progress for all goals
		var sum int64

		// Iterate over the generated goal types and create UserDailyGoals records
		for _, val := range mp {

			rand.Seed(time.Now().UnixNano())
			var record model.UserDailyGoals
			record.UserId = it.Id
			record.GoalType = int64(val)

			//selecting the random currency type
			currencyType := []int{utils.C_ADS, utils.C_GEMS}
			randCurrency := currencyType[rand.Intn(len(currencyType))]
			record.CurrencyType = int64(randCurrency)
			if randCurrency == utils.C_ADS {
				record.Price = 1
			} else {
				rand.Seed(time.Now().UnixNano())
				record.Price = utils.RoundToNearestMultiple(int64(rand.Intn(100-50+1)+50), 10)
			}

			if val == int(utils.PLAYERS_KILLED) {
				min := 20
				max := 80

				lowerRange := min + (int(it.Level-1) * ((max - min) / utils.TOTAL_LEVELS))
				upperRange := min + (int(it.Level) * ((max - min) / utils.TOTAL_LEVELS))

				fmt.Println("lower range and upper range ", lowerRange, upperRange)

				kills := generateRandomNumber(int(it.Level), lowerRange, upperRange)

				baseCoins := 20
				baseGems := 2
				record.Coins = int64(baseCoins) * kills
				record.Gems = int64(baseGems) * kills
				record.TotalProgress = kills

				sum += kills

			} else if val == int(utils.MINI_GAMES_PLAYED) {
				min := 3
				max := 5
				baseCoins := 25
				baseGems := 8

				fmt.Println("Generating mini game played goal")

				gamPlay := generateRandomNumber(int(it.Level), min, max)
				record.Coins = int64(baseCoins) * gamPlay
				record.Gems = int64(baseGems) * gamPlay
				record.TotalProgress = gamPlay
				sum += gamPlay

			} else if val == int(utils.ZOMBIES_KILLED) {

				min := 30
				max := 90

				fmt.Println("Generating zombies killed goal")
				lowerRange := min + (int(it.Level-1) * ((max - min) / utils.TOTAL_LEVELS))
				upperRange := min + (int(it.Level) * ((max - min) / utils.TOTAL_LEVELS))

				kills := generateRandomNumber(int(it.Level), lowerRange, upperRange)

				baseCoins := 20
				baseGems := 2
				record.Coins = int64(baseCoins) * kills
				record.Gems = int64(baseGems) * kills
				record.TotalProgress = kills
				sum += kills

			} else if val == int(utils.BECAME_ZOMBIE) {
				min := 2
				max := 6

				baseCoins := 20
				baseGems := 2

				fmt.Println("Generating became zombie goal")

				gamPlay := generateRandomNumber(int(it.Level), min, max)
				record.Coins = int64(baseCoins) * gamPlay
				record.Gems = int64(baseGems) * gamPlay
				record.TotalProgress = gamPlay
				sum += gamPlay

			} else if val == int(utils.ESCAPE_SURVIVOR) {

				min := 3
				max := 10
				baseCoins := 50
				baseGems := 15
				fmt.Println("Generating escape survivor goal")

				gamPlay := generateRandomNumber(int(it.Level), min, max)
				record.Coins = int64(baseCoins) * gamPlay
				record.Gems = int64(baseGems) * gamPlay
				record.TotalProgress = gamPlay
				sum += gamPlay

			}
			record.DailyRewardId = rewardId

			temp = append(temp, record)
		}

		// fmt.Println("Sum is ", sum)

		MaxCoins := 1000
		Maxgems := 100
		MaxEnergy := 8

		// Create a DailyGoalRewards record with aggregated values for coins, gems, and energy
		baseCoins := float64(MaxCoins / 200)
		baseGems := float64(float64(Maxgems) / 200)
		baseEnergy := float64(float64(MaxEnergy) / 200)

		// fmt.Println("base coins is ", baseCoins)
		// fmt.Println("base gems is ", baseGems)
		// fmt.Println("base energy is", baseEnergy)
		record := model.DailyGoalRewards{
			Id:     rewardId,
			Coins:  (utils.RoundToNearestMultiple(int64((baseCoins)*float64(sum)), 10)),
			Gems:   (utils.RoundToNearestMultiple((int64(baseGems * float64(sum))), 10)),
			Energy: int64(baseEnergy * float64(sum)),
		}

		// Set Chest to 1 if the user's level is greater than 5
		if it.Level > 5 {
			record.Chest = 1
		}

		err := db.CreateRecord(&record)
		if err != nil {
			fmt.Println("Error in creting the entry in db.")
			return err
		}

		err = db.CreateRecord(&temp)
		if err != nil {
			fmt.Println("Error in creting the entry in db.")
			return err
		}
	}

	fmt.Println("Sucessfully generated daily goals for all the users")

	return nil

}

func DeleteAllGoals() {
	query := "DELETE FROM user_daily_goals"
	err := db.RawExecutor(query)
	if err != nil {
		fmt.Println("Error in deleting the records from the tabale")
		return
	}

	query = "DELETE FROM daily_goal_rewards"
	err = db.RawExecutor(query)
	if err != nil {
		fmt.Println("Error in deleting the records from the tabale")
		return
	}
}

func generateRandomNumber(seed, min, max int) int64 {
	rand.Seed(int64(seed)) // Seed the random number generator

	// Generate a random number within the desired range
	randomValue := rand.Intn(max-min+1) + min

	return utils.RoundToNearestMultiple(int64(randomValue), 10)
}

type RewardResponse struct {
	GoalsData []struct {
		Id              string                    `json:"id"`
		GoalType        int64                     `json:"goalType"`
		CurrentProgress int64                     `json:"currentProgress"`
		TotalProgress   int64                     `json:"totalProgress"`
		CurrencyType    int64                     `json:"currencyType"`
		Price           int64                     `json:"price"`
		RewardData      []response.RewardResponse `json:"rewardData"`
	} `json:"goalsData"`
	Timer struct {
		Hour   int64 `json:"hour"`
		Minute int64 `json:"minute"`
	} `json:"timer"`
	Claimed bool                      `json:"claimed"`
	Rewards []response.RewardResponse `json:"rewards"`
}

func GetDailyGoalsService(ctx *gin.Context, userId string) {

	var singleDailyGoal []model.UserDailyGoals

	query := "SELECT * FROM user_daily_goals WHERE user_id=?"
	err := db.QueryExecutor(query, &singleDailyGoal, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	var totalReward model.DailyGoalRewards
	query = "SELECT * FROM daily_goal_rewards WHERE id=?"
	err = db.QueryExecutor(query, &totalReward, singleDailyGoal[0].DailyRewardId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	res := RewardResponse{
		Rewards: []response.RewardResponse{
			{
				RewardType: utils.Coins,
				Quantity:   totalReward.Coins,
			}, {
				RewardType: utils.Gems,
				Quantity:   totalReward.Gems,
			}, {
				RewardType: utils.Energy,
				Quantity:   totalReward.Energy,
			},
		},
	}
	res.Timer = struct {
		Hour   int64 "json:\"hour\""
		Minute int64 "json:\"minute\""
	}{
		Hour:   int64(totalReward.CreatedAt.Add(24 * time.Hour).Sub(time.Now()).Hours()),
		Minute: int64(totalReward.CreatedAt.Add(24*time.Hour).Sub(time.Now()).Minutes()) % 60,
	}

	// if totalReward.Energy != 0 {
	// 	res.Rewards = append(res.Rewards, response.RewardResponse{RewardType: utils.Chest,
	// 		Quantity: totalReward.Chest,
	// 	})
	// }
	res.Claimed = totalReward.Claimed

	for _, data := range singleDailyGoal {
		res.GoalsData = append(res.GoalsData, struct {
			Id              string                    "json:\"id\""
			GoalType        int64                     "json:\"goalType\""
			CurrentProgress int64                     "json:\"currentProgress\""
			TotalProgress   int64                     "json:\"totalProgress\""
			CurrencyType    int64                     `json:"currencyType"`
			Price           int64                     `json:"price"`
			RewardData      []response.RewardResponse "json:\"rewardData\""
		}{
			Id:              data.Id,
			GoalType:        data.GoalType,
			CurrentProgress: data.CurrentProgress,
			TotalProgress:   data.TotalProgress,
			CurrencyType:    data.CurrencyType,
			Price:           data.Price,
			RewardData: []response.RewardResponse{
				{
					RewardType: utils.Coins,
					Quantity:   data.Coins,
				}, {
					RewardType: utils.Gems,
					Quantity:   data.Gems,
				},
			},
		})

	}

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, res, ctx)

}

func SkipGoalService(ctx *gin.Context, userId string, input request.GoalRequest) {

	var userDetails model.UserGameStats
	query := "SELECT * from user_game_stats WHERE user_id=?"
	err := db.QueryExecutor(query, &userDetails, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	var dailyGoal model.UserDailyGoals
	query = "SELECT * FROM user_daily_goals WHERE user_id=? AND id=?"
	err = db.QueryExecutor(query, &dailyGoal, userId, input.Id)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	dailyGoal.CurrentProgress = dailyGoal.TotalProgress

	userDetails.TotalCoins += dailyGoal.Coins
	userDetails.TotalGems += dailyGoal.Gems

	userDetails.CurrentCoins += dailyGoal.Coins
	userDetails.CurrentGems += dailyGoal.Gems

	if dailyGoal.CurrencyType == utils.C_GEMS {
		if userDetails.CurrentGems < dailyGoal.Price {
			response.ShowResponse("Not enough gems", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}
		userDetails.CurrentGems -= dailyGoal.Price
	}

	err = db.UpdateRecord(&dailyGoal, input.Id, "id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	err = db.UpdateRecord(&userDetails, userId, "user_id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	res := []response.RewardResponse{
		{
			RewardType: utils.Coins,
			Quantity:   dailyGoal.Coins,
		}, {
			RewardType: utils.Gems,
			Quantity:   dailyGoal.Gems,
		},
	}

	response.ShowResponse(utils.SUCCESS, utils.HTTP_OK, utils.SUCCESS, res, ctx)

}

func ClaimDailyGoalService(ctx *gin.Context, userId string) {

	var completed bool
	query := ` SELECT
		(SELECT COUNT(*) FROM user_daily_goals WHERE user_id = ? AND total_progress = current_progress) =
		(SELECT COUNT(*) FROM user_daily_goals WHERE user_id = ?) AS result;`
	err := db.QueryExecutor(query, &completed, userId, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	if !completed {
		response.ShowResponse("Tasks not completed", utils.HTTP_BAD_REQUEST, utils.SUCCESS, nil, ctx)
		return
	}

	var userDetails model.UserGameStats
	query = "SELECT * from user_game_stats WHERE user_id=?"
	err = db.QueryExecutor(query, &userDetails, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	var id string
	query = "SELECT daily_reward_id from user_daily_goals WHERE user_id=? limit 1"
	err = db.QueryExecutor(query, &id, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	var dailyGoalReward model.DailyGoalRewards
	query = "SELECT * FROM daily_goal_rewards WHERE id=?"
	err = db.QueryExecutor(query, &dailyGoalReward, id)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	if dailyGoalReward.Claimed {
		response.ShowResponse("Reward already claimed", utils.HTTP_BAD_REQUEST, utils.SUCCESS, nil, ctx)
		return

	}

	userDetails.TotalCoins += dailyGoalReward.Coins
	userDetails.TotalGems += dailyGoalReward.Gems

	userDetails.CurrentCoins += dailyGoalReward.Coins
	userDetails.CurrentGems += dailyGoalReward.Gems
	userDetails.Energy += dailyGoalReward.Energy

	err = db.UpdateRecord(&userDetails, userId, "user_id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	query = "UPDATE daily_goal_rewards SET claimed=true WHERE id=?"
	err = db.RawExecutor(query, id)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	res := []response.RewardResponse{
		{
			RewardType: utils.Coins,
			Quantity:   dailyGoalReward.Coins,
		}, {
			RewardType: utils.Gems,
			Quantity:   dailyGoalReward.Gems,
		},
		{
			RewardType: utils.Energy,
			Quantity:   dailyGoalReward.Energy,
		},
		// {
		// 	RewardType: utils.Chest,
		// 	Quantity:   dailyGoalReward.Chest,
		// },
	}

	response.ShowResponse(utils.SUCCESS, utils.HTTP_OK, utils.SUCCESS, res, ctx)

}

func UpdateGoalService(ctx *gin.Context, userId string, input request.UpdateGoalReq) {

	var userGameStats model.UserGameStats
	query := "SELECT * FROM user_game_stats where user_id=?"
	err := db.QueryExecutor(query, &userGameStats, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	var userDailyGoals []model.UserDailyGoals
	query = "SELECT * FROM user_daily_goals WHERe user_id=?"
	err = db.QueryExecutor(query, &userDailyGoals, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	var res []struct {
		Id              string                    `json:"id"`
		GoalType        int64                     `json:"goalType"`
		CurrentProgress int64                     `json:"currentProgress"`
		TotalProgress   int64                     `json:"totalProgress"`
		CurrencyType    int64                     `json:"currencyType"`
		Price           int64                     `json:"price"`
		RewardData      []response.RewardResponse `json:"rewardData"`
	}

	completeCount := 0
	for _, goal := range userDailyGoals {
		switch goal.GoalType {
		case int64(utils.PLAYERS_KILLED):
			{
				goal.CurrentProgress += input.KillAsZomb
				if goal.CurrentProgress >= goal.TotalProgress {
					completeCount++
					goal.CurrentProgress = goal.TotalProgress
					res = append(res, struct {
						Id              string                    "json:\"id\""
						GoalType        int64                     "json:\"goalType\""
						CurrentProgress int64                     "json:\"currentProgress\""
						TotalProgress   int64                     "json:\"totalProgress\""
						CurrencyType    int64                     "json:\"currencyType\""
						Price           int64                     "json:\"price\""
						RewardData      []response.RewardResponse "json:\"rewardData\""
					}{
						Id:              goal.Id,
						GoalType:        goal.GoalType,
						CurrentProgress: goal.CurrentProgress,
						TotalProgress:   goal.TotalProgress,
						CurrencyType:    goal.CurrencyType,
						Price:           goal.Price,
						RewardData: []response.RewardResponse{
							{
								RewardType: utils.Coins,
								Quantity:   goal.Coins,
							}, {
								RewardType: utils.Gems,
								Quantity:   goal.Gems,
							},
						},
					})
				}
			}
		case int64(utils.BECAME_ZOMBIE):
			{
				goal.CurrentProgress++
				if goal.CurrentProgress >= goal.TotalProgress {
					completeCount++
					goal.CurrentProgress = goal.TotalProgress
					res = append(res, struct {
						Id              string                    "json:\"id\""
						GoalType        int64                     "json:\"goalType\""
						CurrentProgress int64                     "json:\"currentProgress\""
						TotalProgress   int64                     "json:\"totalProgress\""
						CurrencyType    int64                     "json:\"currencyType\""
						Price           int64                     "json:\"price\""
						RewardData      []response.RewardResponse "json:\"rewardData\""
					}{
						Id:              goal.Id,
						GoalType:        goal.GoalType,
						CurrentProgress: goal.CurrentProgress,
						TotalProgress:   goal.TotalProgress,
						CurrencyType:    goal.CurrencyType,
						Price:           goal.Price,
						RewardData: []response.RewardResponse{
							{
								RewardType: utils.Coins,
								Quantity:   goal.Coins,
							}, {
								RewardType: utils.Gems,
								Quantity:   goal.Gems,
							},
						},
					})
				}
			}
		case int64(utils.MINI_GAMES_PLAYED):
			{
				goal.CurrentProgress++
				if goal.CurrentProgress >= goal.TotalProgress {
					completeCount++
					goal.CurrentProgress = goal.TotalProgress
					res = append(res, struct {
						Id              string                    "json:\"id\""
						GoalType        int64                     "json:\"goalType\""
						CurrentProgress int64                     "json:\"currentProgress\""
						TotalProgress   int64                     "json:\"totalProgress\""
						CurrencyType    int64                     "json:\"currencyType\""
						Price           int64                     "json:\"price\""
						RewardData      []response.RewardResponse "json:\"rewardData\""
					}{
						Id:              goal.Id,
						GoalType:        goal.GoalType,
						CurrentProgress: goal.CurrentProgress,
						TotalProgress:   goal.TotalProgress,
						CurrencyType:    goal.CurrencyType,
						Price:           goal.Price,
						RewardData: []response.RewardResponse{
							{
								RewardType: utils.Coins,
								Quantity:   goal.Coins,
							}, {
								RewardType: utils.Gems,
								Quantity:   goal.Gems,
							},
						},
					})
				}
			}
		case int64(utils.ZOMBIES_KILLED):
			{

				goal.CurrentProgress += input.KillAsSur
				if goal.CurrentProgress >= goal.TotalProgress {
					completeCount++
					goal.CurrentProgress = goal.TotalProgress
					res = append(res, struct {
						Id              string                    "json:\"id\""
						GoalType        int64                     "json:\"goalType\""
						CurrentProgress int64                     "json:\"currentProgress\""
						TotalProgress   int64                     "json:\"totalProgress\""
						CurrencyType    int64                     "json:\"currencyType\""
						Price           int64                     "json:\"price\""
						RewardData      []response.RewardResponse "json:\"rewardData\""
					}{
						Id:              goal.Id,
						GoalType:        goal.GoalType,
						CurrentProgress: goal.CurrentProgress,
						TotalProgress:   goal.TotalProgress,
						CurrencyType:    goal.CurrencyType,
						Price:           goal.Price,
						RewardData: []response.RewardResponse{
							{
								RewardType: utils.Coins,
								Quantity:   goal.Coins,
							}, {
								RewardType: utils.Gems,
								Quantity:   goal.Gems,
							},
						},
					})
				}
			}
		}

	}

	response.ShowResponse("Goal completed", utils.HTTP_OK, utils.SUCCESS, res, ctx)

}
