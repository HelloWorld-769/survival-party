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

func GetPlayerLevelRewards(ctx *gin.Context, userId string) {

	var userReward []model.UserLevelRewards
	query := "select * from user_level_rewards where user_id=? order by level"
	err := db.QueryExecutor(query, &userReward, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	response.ShowResponse("playerRewards Get success", utils.HTTP_OK, utils.SUCCESS, userReward, ctx)
}

func GenerateLevelReward(userId string) error {

	//create user level reward for this user
	var levelRewards []model.LevelRewards

	//fetch all the level rewards from db
	query := "select * from level_rewards"
	err := db.QueryExecutor(query, &levelRewards)
	if err != nil {
		fmt.Println("error fetching level rewards", err)
		return err
	}

	for _, r := range levelRewards {

		var levelReward model.UserLevelRewards
		levelReward.UserId = userId
		levelReward.Status = utils.UNAVAILABLE
		if r.LevelRequired == 0 {

			levelReward.Status = utils.UNCLAIMED
		}
		levelReward.RewardType = r.RewardType
		levelReward.Quantity = r.Quantity
		levelReward.Level = r.LevelRequired

		err := db.CreateRecord(&levelReward)
		if err != nil {
			return err
		}

	}

	return nil

}

func PlayerLevelRewardCollect(ctx *gin.Context, userId string, req request.PlayerLevelRewardCollectRequest) {

	//check whether user has enough level to collect the levelReward
	//user level reward dummy for testing
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
	var res response.RewardResponse
	res.Quantity = LevelReward.Quantity
	switch LevelReward.RewardType {
	case utils.Coins:
		fmt.Println("Coins")
		userData.CurrentCoins += LevelReward.Quantity
		userData.TotalCoins += LevelReward.Quantity
		res.RewardType = utils.Coins
	case utils.Energy:
		fmt.Println("Energy")
		userData.Energy += LevelReward.Quantity
		res.RewardType = utils.Energy
	case utils.Gems:
		fmt.Println("Gems")
		userData.CurrentGems += LevelReward.Quantity
		userData.TotalGems += LevelReward.Quantity
		res.RewardType = utils.Gems

	default:
		fmt.Println("nothing for reward")
	}

	err = db.UpdateRecord(&userData, userId, "user_id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	//update the claimed field in user level reward
	query = "UPDATE user_level_rewards SET status=? WHERE user_id=? AND level=?"
	db.RawExecutor(query, utils.CLAIMED, userId, req.Level)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("reward collected successfully", utils.HTTP_OK, utils.SUCCESS, res, ctx)

}

func AddPlayerLevel() {
	if !utils.TableIsEmpty("level_rewards") {
		multiplier := 1.45

		rewardsRange := map[int][]int{
			1: {2, 5},
			2: {100, 1000},
			3: {1, 10},
		}

		// Initialize the result slice with the first two elements
		res := []model.LevelRewards{
			{LevelRequired: 0, XpRequired: 0, Quantity: 10, RewardType: 2},
			{LevelRequired: 1, XpRequired: 50, Quantity: 25, RewardType: 1},
		}

		// Loop from 2 to 49
		for i := 2; i <= 10; i++ {
			// Determine the multiplier based on the range

			// Calculate the new value and round it to the nearest multiple of 5
			val := float64(res[i-1].XpRequired) * multiplier
			xp := utils.RoundToNearestMultipleOf5(val)

			rewardType := rand.Intn(4-2) + 2
			rewardQuant := utils.RoundToNearestMultiple(int64(rand.Intn(rewardsRange[rewardType][1]-rewardsRange[rewardType][0])+rewardsRange[rewardType][0]), 10)

			// Append the values to the result slice
			res = append(res, model.LevelRewards{LevelRequired: int64(i), XpRequired: int64(xp), RewardType: int64(rewardType), Quantity: rewardQuant})

		}
		err := db.CreateRecord(&res)
		if err != nil {
			fmt.Println("Error is:", err)
			return
		}
	}

}
