package rewards

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
)

func GetPlayerLevelRewards(ctx *gin.Context, userId string) {

	var userReward []model.UserLevelRewards
	query := "select * from user_level_rewards where user_id=?"
	err := db.QueryExecutor(query, &userReward, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	response.ShowResponse("playerRewards Get success", utils.HTTP_OK, utils.SUCCESS, userReward, ctx)
}

func PlayerLevelRewardCollect(ctx *gin.Context, userId string, req request.PlayerLevelRewardCollectRequest) {

	//check whether user has enough level to collect the levelReward
	//user level reward dummy for testing
	var userLevelRecord model.UserLevelRewards
	userData, err := utils.GetUserGameStatsData(userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	if userData.Level < req.Level {
		//user is not allowed to collect this reward
		// fmt.Println("userdata level", userData.User.Level)
		// fmt.Println("req.Level", req.Level)
		// fmt.Println("user level", userData.Level)
		response.ShowResponse("Not enough user level ", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	//find the lvl_reward corresponding to the input lvl
	var LevelReward model.LevelRewards
	query := "select * from level_rewards where level_required=?"
	err = db.QueryExecutor(query, &LevelReward, req.Level)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//add user rewards to gamestats of user
	switch LevelReward.RewardType {
	case utils.Coins:
		fmt.Println("Coins")
		userData.CurrentCoins += LevelReward.Quantity
		userData.TotalCoins += LevelReward.Quantity

	case utils.Energy:
		fmt.Println("Energy")
		userData.Energy += LevelReward.Quantity
	case utils.Gems:
		fmt.Println("Gems")
		userData.CurrentGems += LevelReward.Quantity
		userData.TotalGems += LevelReward.Quantity

	default:
		fmt.Println("nothing for reward")
	}
	err = db.UpdateRecord(&userData, userId, "user_id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	//update the claimed field in user level reward
	query = "select * from user_level_rewards where user_id=?"
	err = db.QueryExecutor(query, &userLevelRecord, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	userLevelRecord.Claimed = true
	err = db.UpdateRecord(&userLevelRecord, userId, "user_id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("reward collected successfully", utils.HTTP_OK, utils.SUCCESS, nil, ctx)

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
		var dailyReward model.DailyRewards
		if user.EmailVerified {

			dayCount := utils.CalculateDays(user.CreatedAt) + 1
			query := "select * from daily_rewards where day_count=?"
			err := db.QueryExecutor(query, &dailyReward, dayCount)

		var dayCount int
		query := "select day_count from users where emailverified =true and id=?"
		err := db.QueryExecutor(query, &dayCount, user.Id)
		if err != nil {
			fmt.Println("error ", err.Error())
			return
		}
		fmt.Println("user ", user.Id)
		fmt.Println("day Count: ", dayCount)

		if dayCount%8 == 0 && dayCount >= 7 {

			//delete previous 7 daily reward entries for this user
			query := "delete from user_daily_rewards where user_id =?"
			err := db.QueryExecutor(query, nil, user.Id)
			if err != nil {

				fmt.Println("error in fetching", err.Error())
				return
			}

			//create entry of this daily reward for this user
			var daily_user_reward model.UserDailyRewards
			daily_user_reward.UserId = user.Id
			daily_user_reward.Coins += dailyReward.Coins
			daily_user_reward.Energy += dailyReward.Energy
			daily_user_reward.Gems += dailyReward.Gems

			err = db.CreateRecord(&daily_user_reward)
			if err != nil {
				fmt.Println("error in creating", err.Error())
				return
				//create entry of this daily reward for this user
				var daily_user_reward model.UserDailyRewards
				daily_user_reward.UserId = user.Id
				//find random reward Type
				//append the quantity into reward
				randomInt := rand.Intn(6)
				if randomInt == 3 {
					//gems
					randomInt := int(Multiplier) * (rand.Intn(10))
					daily_user_reward.Gain = int64(randomInt)

				} else if randomInt == 4 {
					//inventory
					//set asset name
					daily_user_reward.AssetName = "egg_hat" //(can be random asset in future)
					daily_user_reward.Gain = 1
				} else if randomInt == 5 {
					//Chest
					//set chest level
					randomInt := rand.Intn(6)
					daily_user_reward.ChestType = int64(randomInt)
					daily_user_reward.Gain = 1
				} else {
					randomInt := int(Multiplier) * (rand.Intn(100) + rand.Intn(50))
					daily_user_reward.Gain = int64(randomInt)
				}
				daily_user_reward.RewardType = int64(randomInt)

				if i == 0 {
					daily_user_reward.Status = utils.UNCLAIMED
				}
				daily_user_reward.Status = utils.UNAVAILABLE
				fmt.Println("entry created for user!!!")
				err = db.CreateRecord(&daily_user_reward)
				if err != nil {
					fmt.Println("error in creating", err.Error())
					return
				}
			}
		}

	}
}

func CreateStarterDailyRewards(userId string) error {

	//check if already created
	var count int64
	query := "select count(*) from user_daily_reward where user_id=?"
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
			Multiplier := utils.UserMultipler(userId)

			//find random reward Type
			//append the quantity into reward
			randomInt := rand.Intn(6)
			if randomInt == 5 {
				//inventory
				//set asset name
				daily_user_reward.AssetName = "egg_hat" //(can be random asset in future)
				daily_user_reward.Gain = 1
			} else if randomInt == 6 {
				//Chest
				//set chest level
				randomInt := rand.Intn(6)
				daily_user_reward.ChestType = int64(randomInt)
				daily_user_reward.Gain = 1
			} else {
				randomInt := int(Multiplier) * (rand.Intn(100) + rand.Intn(50))
				daily_user_reward.Gain = int64(randomInt)
			}
			daily_user_reward.RewardType = int64(randomInt)

			if i == 0 {
				//for the first daily reward
				daily_user_reward.Status = utils.UNCLAIMED
			}
			daily_user_reward.Status = utils.UNAVAILABLE

			err = db.CreateRecord(&daily_user_reward)
			if err != nil {
				fmt.Println("error in creating", err.Error())
				return err
			}
		}
	}
	return nil

}

func CollectDailyReward(ctx *gin.Context, userId string) {

	//get user data
	var userGameStats model.UserGameStats
	query := "select * from user_game_stats where user_id=?"
	err := db.QueryExecutor(query, &userGameStats, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	//get daily reward data
	var userRewardData model.UserDailyRewards
	query = "select * from user_daily_rewards where user_id=?"
	err = db.QueryExecutor(query, &userRewardData)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	//update this userRewardData as claimed true
	userRewardData.Status = utils.CLAIMED
	err = db.UpdateRecord(&userRewardData, userId, "user_id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	switch userRewardData.RewardType {
	case 1:
		fmt.Println("Energy")
		userGameStats.Energy += userRewardData.Gain
	case 2:
		fmt.Println("Coins")
		userGameStats.CurrentCoins += userRewardData.Gain
		userGameStats.TotalCoins += userRewardData.Gain
	case 3:
		fmt.Println("Gems")
		userGameStats.CurrentGems += userRewardData.Gain
		userGameStats.TotalGems += userRewardData.Gain

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
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//fetch the updated daily rewards and give in response
	var dailyRewards []model.UserDailyRewards
	query = "select * from daily_rewards where user_id=?"
	err = db.QueryExecutor(query, &dailyRewards, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("reward claimed successfully ", utils.HTTP_OK, utils.SUCCESS, dailyRewards, ctx)

}

// todo
// get user daily reward data
func GetUserDailyRewardData(ctx *gin.Context, userId string) {

	var UserDailyRewardsData model.UserDailyRewards
	query := "select * from user_daily_rewards where user_id = ?"
	err := db.QueryExecutor(query, &UserDailyRewardsData, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("daily reward fetched successfully ", utils.HTTP_OK, utils.SUCCESS, UserDailyRewardsData, ctx)

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
		query := "select daycount from users where emailverified =true and id=?"
		db.QueryExecutor(query, &dayCount, user.Id)
		if dayCount%7 != 0 {

			//other than last day or first of daily reward weekly pack
			//make the status missed if still unclaimed
			var userDailyReward model.UserDailyRewards
			query := "select * from user_daily_rewards where user_id=? and daycount=?"
			err := db.QueryExecutor(query, userDailyReward, user.Id, dayCount)
			if err != nil {
				fmt.Println("error ", err.Error())
				return
			}
			if userDailyReward.Status == utils.UNCLAIMED {

				//mark as Missed
				userDailyReward.Status = utils.MISSED
				err := db.UpdateRecord(&userDailyReward, user.Id, "user_id").Error
				if err != nil {
					fmt.Println("error ", err.Error())
					return
				}
			}
			//make the next day reward status from unavailbale to unclaimed
			query = "select * from user_daily_rewards where user_id=? and daycount=?"
			err = db.QueryExecutor(query, userDailyReward, user.Id, dayCount+1)
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
