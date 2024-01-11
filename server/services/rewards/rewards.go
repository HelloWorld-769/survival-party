package rewards

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"

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
		fmt.Println("userdata level", userData.User.Level)
		fmt.Println("req.Level", req.Level)
		fmt.Println("user level", userData.Level)
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

	//fetch all the users
	var allUsers []model.User
	query := "select * from users"
	err := db.QueryExecutor(query, &allUsers)
	if err != nil {
		fmt.Println("error in fetching users query:", err.Error())
	}

	for _, user := range allUsers {

		//create user daily reward entry based on user daycount
		//if user daycount is in between 1-7 then allot starting rewards else different formula based
		//calculate user daycount
		dayCount := utils.CalculateDays(user.CreatedAt) + 1
		var dailyReward []model.DailyRewards
		if user.Emailverified {

			if dayCount <= 6 {

				query := "select * from daily_rewards"
				err := db.QueryExecutor(query, &dailyReward)
				if err != nil {

					fmt.Println("error in fetching", err.Error())
					return
				}

				//create entry of this daily reward for this user
				for i := 0; i < 7; i++ {
					var daily_user_reward model.UserDailyRewards
					daily_user_reward.UserId = user.Id
					daily_user_reward.Coins += dailyReward[i].Coins
					daily_user_reward.Energy += dailyReward[i].Energy
					daily_user_reward.Gems += dailyReward[i].Gems

					err = db.CreateRecord(&daily_user_reward)
					if err != nil {
						fmt.Println("error in creating", err.Error())
						return
					}
				}
			} else {
				//generate daily reward based on formula
				//formula based on users gameplay_time and users created_at
				Multiplier := utils.UserMultipler(user.Id)

				//create entry of this daily reward for this user
				var daily_user_reward model.UserDailyRewards
				daily_user_reward.UserId = user.Id
				daily_user_reward.Coins += 10 * Multiplier
				daily_user_reward.Energy += 10 * Multiplier
				daily_user_reward.Gems += 1 * Multiplier

				err = db.CreateRecord(&daily_user_reward)
				if err != nil {
					fmt.Println("error in creating", err.Error())
					return
				}

			}
		}

	}
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
	userRewardData.Claimed = true
	err = db.UpdateRecord(&userRewardData, userId, "user_id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//update user game stats with reward data
	userGameStats.CurrentCoins += userRewardData.Coins
	userGameStats.TotalCoins += userRewardData.Coins
	userGameStats.CurrentGems += userRewardData.Gems
	userGameStats.TotalGems += userRewardData.Gems
	userGameStats.Energy += userRewardData.Energy

	err = db.UpdateRecord(&userGameStats, userId, "user_id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("reward claimed successfully ", utils.HTTP_OK, utils.SUCCESS, nil, ctx)

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

}
