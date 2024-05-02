package rewards

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func CollectDailyReward(ctx *gin.Context, userId string, input request.DailyRewardMuti) {

	//get user gamestats data
	var userGameStats model.UserGameStats
	query := "select * from user_game_stats where user_id=?"
	err := db.QueryExecutor(query, &userGameStats, userId)
	if err != nil {
		fmt.Println("here1")
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//get the user data

	user, err := utils.GetUserData(userId)
	if err != nil {
		fmt.Println("here2")

		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	//get daily reward data
	var userRewardData []model.UserDailyRewards
	query = "select * from user_daily_rewards where user_id=? order by day_count asc "
	err = db.QueryExecutor(query, &userRewardData, userId)
	if err != nil {
		fmt.Println("here3")

		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	fmt.Println("user daycount:   ", user.DayCount)
	fmt.Println("user daycount mod 7:   ", (user.DayCount % 7))

	var muliplier int64
	switch input.Type {
	case int(utils.ONE):
		muliplier = 1
	case int(utils.TWO):
		muliplier = 2
	case int(utils.THREE):
		muliplier = 3
	default:
		response.ShowResponse("Invalid Type", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	//check if already claimed
	var gains int64
	rewardDay := (user.DayCount - 1) % 7

	if userRewardData[rewardDay].Status == utils.CLAIMED {
		response.ShowResponse("daily reward already claimed", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	} else {
		//update this userRewardData as claimed true
		userRewardData[rewardDay].Status = utils.CLAIMED

		err = db.UpdateRecord(&userRewardData[rewardDay], userId, "user_id").Error
		if err != nil {
			fmt.Println("here4")

			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

	}

	switch userRewardData[rewardDay].RewardType {
	case 1:
		fmt.Println("Energy")
		userGameStats.Energy += muliplier * userRewardData[rewardDay].Gain
		gains = muliplier * userRewardData[rewardDay].Gain

	case 2:
		fmt.Println("Coins")
		userGameStats.CurrentCoins += muliplier * userRewardData[rewardDay].Gain
		userGameStats.TotalCoins += muliplier * userRewardData[rewardDay].Gain
		gains = muliplier * userRewardData[rewardDay].Gain

	case 3:
		fmt.Println("Gems")
		userGameStats.CurrentGems += muliplier * userRewardData[rewardDay].Gain
		userGameStats.TotalGems += muliplier * userRewardData[rewardDay].Gain
		gains = muliplier * userRewardData[rewardDay].Gain

	case 4:
		fmt.Println("Inventory")
	case 5:
		fmt.Println("Chest")
	default:
		fmt.Println("Invalid")
	}

	//update user game stats with reward data

	err = db.UpdateRecord(&userGameStats, userId, "user_id").Error
	if err != nil {
		fmt.Println("here5")
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	// } else if user.DayCount < 7 {
	// 	if userRewardData[(user.DayCount)-1].Status == utils.CLAIMED {

	// 		response.ShowResponse("daily reward already claimed", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
	// 		return
	// 	} else {
	// 		//update this userRewardData as claimed true
	// 		userRewardData[(user.DayCount % 7)].Status = utils.CLAIMED

	// 		err = db.UpdateRecord(&userRewardData[(user.DayCount)-1], userId, "user_id").Error
	// 		if err != nil {
	// 			fmt.Println("here4")

	// 			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
	// 			return
	// 		}

	// 	}

	// 	switch userRewardData[(user.DayCount)-1].RewardType {
	// 	case 1:
	// 		fmt.Println("Energy")
	// 		userGameStats.Energy += muliplier * userRewardData[(user.DayCount)-1].Gain
	// 		gains = muliplier * userRewardData[(user.DayCount)-1].Gain

	// 	case 2:
	// 		fmt.Println("Coins")
	// 		userGameStats.CurrentCoins += muliplier * userRewardData[(user.DayCount)-1].Gain
	// 		userGameStats.TotalCoins += muliplier * userRewardData[(user.DayCount)-1].Gain
	// 		gains = muliplier * userRewardData[(user.DayCount)-1].Gain

	// 	case 3:
	// 		fmt.Println("Gems")
	// 		userGameStats.CurrentGems += muliplier * userRewardData[(user.DayCount)-1].Gain
	// 		userGameStats.TotalGems += muliplier * userRewardData[(user.DayCount)-1].Gain
	// 		gains = muliplier * userRewardData[(user.DayCount)-1].Gain

	// 	case 4:
	// 		fmt.Println("Inventory")
	// 	case 5:
	// 		fmt.Println("Chest")
	// 	default:
	// 		fmt.Println("Invalid")
	// 	}

	// 	//update user game stats with reward data

	// 	err = db.UpdateRecord(&userGameStats, userId, "user_id").Error
	// 	if err != nil {
	// 		fmt.Println("here5")
	// 		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
	// 		return
	// 	}
	// } else {
	// 	if userRewardData[(user.DayCount)%7].Status == utils.CLAIMED {

	// 		response.ShowResponse("daily reward already claimed", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
	// 		return
	// 	} else {
	// 		//update this userRewardData as claimed true
	// 		userRewardData[(user.DayCount)%7].Status = utils.CLAIMED

	// 		err = db.UpdateRecord(&userRewardData[(user.DayCount%7)], userId, "user_id").Error
	// 		if err != nil {
	// 			fmt.Println("here4")

	// 			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
	// 			return
	// 		}

	// 	}

	// 	switch userRewardData[(user.DayCount % 7)].RewardType {
	// 	case 1:
	// 		fmt.Println("Energy")
	// 		userGameStats.Energy += muliplier * userRewardData[(user.DayCount%7)].Gain
	// 		gains = muliplier * userRewardData[(user.DayCount%7)].Gain

	// 	case 2:
	// 		fmt.Println("Coins")
	// 		userGameStats.CurrentCoins += muliplier * userRewardData[(user.DayCount%7)].Gain
	// 		userGameStats.TotalCoins += muliplier * userRewardData[(user.DayCount%7)].Gain
	// 		gains = muliplier * userRewardData[(user.DayCount%7)-1].Gain

	// 	case 3:
	// 		fmt.Println("Gems")
	// 		userGameStats.CurrentGems += muliplier * userRewardData[(user.DayCount%7)].Gain
	// 		userGameStats.TotalGems += muliplier * userRewardData[(user.DayCount%7)].Gain
	// 		gains = muliplier * userRewardData[(user.DayCount%7)].Gain

	// 	case 4:
	// 		fmt.Println("Inventory")
	// 	case 5:
	// 		fmt.Println("Chest")
	// 	default:
	// 		fmt.Println("Invalid")
	// 	}

	// 	//update user game stats with reward data

	// 	err = db.UpdateRecord(&userGameStats, userId, "user_id").Error
	// 	if err != nil {
	// 		fmt.Println("here5")
	// 		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
	// 		return
	// 	}
	// }

	response.ShowResponse("reward claimed successfully ", utils.HTTP_OK, utils.SUCCESS, struct {
		UserDailyRewards []model.UserDailyRewards `json:"userDailyRewards"`
		Gains            int64                    `json:"gains"`
	}{
		UserDailyRewards: userRewardData,
		Gains:            gains,
	}, ctx)
}

func UpdateDailyRewardsData() {

	//to update the status of the daily reward
	//fetch all the users
	//iterate over them and calculate their daycount (mod 7)
	//check the status of the daily reward (if it is unclaimed ,mark it as Missed)
	//And make the next day reward status from unavailbale to unclaimed(available to claim)

	var users []model.User
	query := "select * from users where email_verified = true"
	err := db.QueryExecutor(query, &users)
	if err != nil {
		fmt.Println("error ", err.Error())
		return
	}

	for _, user := range users {

		var dayCount int
		query := "select day_count from users where email_verified =true and id=?"
		db.QueryExecutor(query, &dayCount, user.Id)

		if dayCount%7 != 1 {

			//other than first of daily reward weekly pack
			var userDailyReward model.UserDailyRewards
			query := "select * from user_daily_rewards where user_id=? and day_count=?"
			err := db.QueryExecutor(query, &userDailyReward, user.Id, dayCount%7-1)
			if err != nil {
				fmt.Println("error ", err.Error())
				return
			}
			//TODO merge the below if condition into sql query
			//make the status missed if still unclaimed
			if userDailyReward.Status == utils.UNCLAIMED {

				//mark as Missed
				userDailyReward.Status = utils.MISSED
				err := db.UpdateRecord(&userDailyReward, user.Id, "user_id").Error
				if err != nil {
					fmt.Println("error ", err.Error())
					return
				}
			}

			//TODO merge the below if condition into sql query
			//make the next day reward status from unavailble to unclaimed
			query = "select * from user_daily_rewards where user_id=? and day_count=?"
			err = db.QueryExecutor(query, &userDailyReward, user.Id, dayCount%7)
			if err != nil {
				fmt.Println("error ", err.Error())
				return
			}
			userDailyReward.Status = utils.UNCLAIMED
			err = db.UpdateRecord(&userDailyReward, user.Id, "user_id").Error
			if err != nil {
				fmt.Println("error ", err.Error())
				return
			}

		}

	}

}

// get user daily reward data
func GetUserDailyRewardData(ctx *gin.Context, userId string) {

	var UserDailyRewardsData []model.UserDailyRewards
	query := "select * from user_daily_rewards where user_id = ? order by day_count"
	err := db.QueryExecutor(query, &UserDailyRewardsData, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("daily reward fetched successfully ", utils.HTTP_OK, utils.SUCCESS, UserDailyRewardsData, ctx)

}

func DailyRewardTimeLeft(ctx *gin.Context) {

	var timeLeft utils.TimeLeft
	hours, mins, seconds := utils.TimeLeftUntilMidnight()
	timeLeft.Hours = hours
	timeLeft.Minutes = mins
	timeLeft.Seconds = seconds

	response.ShowResponse("time left successfully fetched", utils.HTTP_OK, utils.SUCCESS, timeLeft, ctx)

}

func CreateStarterDailyRewards(userId string) error {

	//check if already created
	var count int64
	query := "select count(*) from user_daily_rewards where user_id=?"
	err := db.QueryExecutor(query, &count, userId)
	if err != nil {
		fmt.Println("error in fetching", err.Error())
		return err
	}
	if count == 0 {

		//first day of the user
		var dailyRewards []model.DailyRewards
		query = "select * from daily_rewards"
		err = db.QueryExecutor(query, &dailyRewards)
		if err != nil {

			fmt.Println("error in fetching", err.Error())
			return err
		}
		//create entry of first week  daily reward for this user
		for i := 0; i < 7; i++ {
			//create entry of this daily reward for this user
			var daily_user_reward model.UserDailyRewards
			daily_user_reward.UserId = userId
			daily_user_reward.DayCount = int64(i + 1)
			Multiplier := utils.UserMultipler(userId)

			//find random reward Type
			//append the quantity into reward
			randomInt := rand.Intn(7)
			//DO NOT REMOVE THESE COMMENTS (MAY BE NEEDED IN FUTURE CODE)

			// if randomInt == 5 {
			// 	//inventory
			// 	//set asset name
			// 	daily_user_reward.AssetName = "egg_hat" //(can be random asset in future)
			// 	daily_user_reward.Name = "egg_hat"
			// 	daily_user_reward.Gain = 1
			// } else if randomInt == 4 {
			// 	//Chest
			// 	//set chest level
			// 	randomInt := rand.Intn(6)
			// 	daily_user_reward.ChestType = int64(randomInt)
			// 	daily_user_reward.Gain = 1
			// } else {
			randomIntgain := int(Multiplier) * (rand.Intn(100) + rand.Intn(50))
			daily_user_reward.Gain = int64(randomIntgain)
			// }
			daily_user_reward.Status = utils.UNAVAILABLE
			daily_user_reward.RewardType = int64(randomInt)
			if i == 0 {
				//for the first daly reward
				daily_user_reward.Status = utils.UNCLAIMED
			}

			err = db.CreateRecord(&daily_user_reward)
			if err != nil {
				fmt.Println("error in creating", err.Error())
				return err
			}
		}
	}
	return nil

}

func CreateUserDailyReward() {

	//create user daily reward entry for all the users in the database based on their daycount
	fmt.Println("create user daily reward called!!!!")
	//fetch all the users
	var allUsers []model.User
	query := "select * from users where email_verified =true"
	err := db.QueryExecutor(query, &allUsers)
	if err != nil {
		fmt.Println("error in fetching users query:", err.Error())
	}

	for _, user := range allUsers {

		//create user daily reward entry based on user daycount
		//calculate user daycount

		var dayCount int
		query := "select day_count from users where email_verified = true and id=?"
		err := db.QueryExecutor(query, &dayCount, user.Id)
		if err != nil {
			fmt.Println("error ", err.Error())
			return
		}

		if dayCount%7 == 1 && dayCount >= 7 {

			//delete previous 7 daily reward entries for this user
			query := "delete from user_daily_rewards where user_id =?"
			err := db.QueryExecutor(query, nil, user.Id)
			if err != nil {
				fmt.Println("error in deleting previous daily reward entries", err.Error())
				return
			}
			for i := 0; i < 7; i++ {

				//generate daily reward based on formula
				//formula based on users gameplay_time and users created_at

				Multiplier := utils.UserMultipler(user.Id)

				//create entry of this daily reward for this user
				var daily_user_reward model.UserDailyRewards
				daily_user_reward.DayCount = int64(i + 1)
				daily_user_reward.UserId = user.Id
				//find random reward Type
				//append the quantity into reward
				daily_user_reward.Status = utils.UNAVAILABLE
				if i == 0 {
					daily_user_reward.Status = utils.UNCLAIMED
				}
				randomInt := 1 + rand.Intn(3)
				// if randomInt == 3 {
				// 	//gems
				// 	randomInt := int(Multiplier) * (rand.Intn(10))
				// 	daily_user_reward.Gain = int64(randomInt)

				// } else if randomInt == 4 {
				// 	//inventory
				// 	//set asset name
				// 	daily_user_reward.AssetName = "egg_hat" //(can be random asset in future)
				// 	daily_user_reward.Name = "egg_hat"
				// 	daily_user_reward.Gain = 1
				// } else if randomInt == 5 {
				// 	//Chest
				// 	//set chest level
				// 	randomInt := rand.Intn(6)
				// 	daily_user_reward.ChestType = int64(randomInt)
				// 	daily_user_reward.Gain = 1
				// } else {
				randomIntgain := int(Multiplier) * (rand.Intn(100) + rand.Intn(50))
				daily_user_reward.Gain = int64(randomIntgain)
				// }
				daily_user_reward.RewardType = int64(randomInt)
				err = db.CreateRecord(&daily_user_reward)
				if err != nil {
					fmt.Println("error in creating", err.Error())
					return
				}
			}
		}

	}
}
