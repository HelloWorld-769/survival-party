package gameplay

import (
	"encoding/json"
	"fmt"
	"io"
	"main/server/db"
	"main/server/response"
	"sort"

	"main/server/model"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

type WebRpc struct {
	UserId    string
	RpcParams struct {
		GameId     string  `json:"gameId"`
		ActionType float64 `json:"actionType"`
		Data       struct {
			ActorNr int `json:"actorNr"`
		} `json:"data"`
	}
}

// This is the func that is used to chnge the state of pplayer like if player kills a player or zombie or complete a game
func InGameState(ctx *gin.Context) {
	var data WebRpc
	body, _ := io.ReadAll(ctx.Request.Body)

	fmt.Println("Bdy is", string(body))

	err := json.Unmarshal(body, &data)
	if err != nil {
		ctx.JSON(int(utils.HTTP_INTERNAL_SERVER_ERROR), map[string]interface{}{
			"ResultCode": "0",
		})
		fmt.Println("Error in unmarshalling the resposne from the hook")
		return
	}

	if data.RpcParams.ActionType == float64(utils.MINI_GAME_COMPLETION) {

		var userGameData model.GameState
		query := "SELECT * FROM game_states WHERE user_id=? AND game_id=?"
		err := db.QueryExecutor(query, &userGameData, data.UserId, data.RpcParams.GameId)
		if err != nil {
			fmt.Println("Error is", err.Error())
			return
		}

		if userGameData.GamesCompleted < userGameData.TotalGames {

			userGameData.GamesCompleted++
			userGameData.Xp += 5
		}

		query = "UPDATE game_states SET xp=?,games_completed=? WHERE user_id=? AND game_id=?"
		err = db.RawExecutor(query, userGameData.Xp, userGameData.GamesCompleted, data.UserId, data.RpcParams.GameId)
		if err != nil {
			//resposne the error
			ctx.JSON(int(utils.HTTP_INTERNAL_SERVER_ERROR), map[string]interface{}{
				"ResultCode": "0",
				// "Data": map[string]interface{}{
				// 	"coins": 100,
				// 	"gems":  10,
				// },
			})
			return
		}

		var userStats model.UserGameStats
		query = "SELECT * FROM user_game_stats WHERE user_id=?"
		err = db.QueryExecutor(query, &userStats, data.UserId)
		if err != nil {
			//resposne the error
			return
		}

		reward, err := ProcessDailyGoal(int64(utils.MINI_GAMES_PLAYED), data.UserId, &userStats)
		if err != nil {

			fmt.Println("Error in processing the daily goal")
			return

		}

		err = db.UpdateRecord(&userStats, data.UserId, "user_id").Error
		if err != nil {
			//resposne the error
			return
		}
		if reward != nil {

			ctx.JSON(200, map[string]interface{}{
				"ResultCode": "1",
				"Data": map[string]interface{}{
					"achievments": reward,
				},
			})
			return
		}

		ctx.JSON(200, map[string]interface{}{
			"ResultCode": "0",
			"Data":       "Data updated sucessfully",
		})

	} else if data.RpcParams.ActionType == float64(utils.KILL_PLAYER) {
		// fmt.Println("asdbkjasbdjkbakdb")

		//If userId is empty the it means its a bot
		if data.UserId == "" {

			//Marking another player dead
			query := "UPDATE game_states set is_dead=true where actor_nr=? and game_id=?"
			err = db.RawExecutor(query, data.RpcParams.Data.ActorNr, data.RpcParams.GameId)
			if err != nil {
				//resposne the error
				return
			}

			ctx.JSON(200, map[string]interface{}{
				"ResultCode": "0",
				"Data":       "Data updated sucessfully",
			})

		} else if data.RpcParams.Data.ActorNr == 0 {
			//Giving the points to the player
			var userGameData model.GameState
			query := "SELECT * FROM game_states WHERE user_id=? AND game_id=?"
			err := db.QueryExecutor(query, &userGameData, data.UserId, data.RpcParams.GameId)
			if err != nil {
				fmt.Println("Error is", err.Error())
				return
			}

			userGameData.Xp += 3
			userGameData.Kills++

			//update userGameData

			query = "UPDATE game_states SET xp=?,kills=? WHERE user_id=? AND game_id=?"
			err = db.RawExecutor(query, userGameData.Xp, userGameData.Kills, data.UserId, data.RpcParams.GameId)
			if err != nil {
				//resposne the error
				return
			}

			//Increasing the total kills
			var userStats model.UserGameStats
			query = "SELECT * FROM user_game_stats WHERE user_id=?"
			err = db.QueryExecutor(query, &userStats, data.UserId)
			if err != nil {
				//resposne the error
				return
			}

			userStats.TotalKills++

			reward, err := ProcessDailyGoal(int64(utils.PLAYERS_KILLED), data.UserId, &userStats)
			if err != nil {

				fmt.Println("Error in processing the daily goal")
				return

			}

			err = db.UpdateRecord(userStats, data.UserId, "user_id").Error
			if err != nil {
				//resposne the error
				return
			}

			if reward != nil {

				ctx.JSON(200, map[string]interface{}{
					"ResultCode": "0",
					"Data": map[string]interface{}{
						"achievments": reward,
					},
				})
				return
			}

			ctx.JSON(200, map[string]interface{}{
				"ResultCode": "0",
				"Data":       "Data updated sucessfully",
			})

		} else {
			//Giving the points to the player
			var userGameData model.GameState
			query := "SELECT * FROM game_states WHERE user_id=? AND game_id=?"
			err := db.QueryExecutor(query, &userGameData, data.UserId, data.RpcParams.GameId)
			if err != nil {
				fmt.Println("Error is", err.Error())
				return
			}

			userGameData.Xp += 2
			userGameData.Kills++

			//update userGameData

			query = "UPDATE game_states SET xp=?,kills=? WHERE user_id=? AND game_id=?"
			err = db.RawExecutor(query, userGameData.Xp, userGameData.Kills, data.UserId, data.RpcParams.GameId)
			if err != nil {
				//resposne the error
				return
			}

			//Increasing the total kills
			var userStats model.UserGameStats
			query = "SELECT * FROM user_game_stats WHERE user_id=?"
			err = db.QueryExecutor(query, &userStats, data.UserId)
			if err != nil {
				//resposne the error
				return
			}

			userStats.TotalKills++

			reward, err := ProcessDailyGoal(int64(utils.PLAYERS_KILLED), data.UserId, &userStats)
			if err != nil {

				fmt.Println("Error in processing the daily goal")
				return

			}

			err = db.UpdateRecord(userStats, data.UserId, "user_id").Error
			if err != nil {
				//resposne the error
				return
			}

			//Marking another player dead
			query = "UPDATE game_states set is_dead=true where actor_nr=? and game_id=?"
			err = db.RawExecutor(query, data.RpcParams.Data.ActorNr, data.RpcParams.GameId)
			if err != nil {
				//resposne the error
				return
			}

			if reward != nil {

				ctx.JSON(200, map[string]interface{}{
					"ResultCode": "0",
					"Data": map[string]interface{}{
						"achievments": reward,
					},
				})
				return
			}
			ctx.JSON(200, map[string]interface{}{
				"ResultCode": "0",
				"Data":       "Data updated sucessfully",
			})
		}

	} else if data.RpcParams.ActionType == float64(utils.KILL_ZOMBIE) {

		var userGameData model.GameState
		query := "SELECT * FROM game_states WHERE user_id=? AND game_id=?"
		err := db.QueryExecutor(query, &userGameData, data.UserId, data.RpcParams.GameId)
		if err != nil {
			fmt.Println("Error is", err.Error())
			return
		}

		userGameData.Xp += 1
		userGameData.Kills++
		query = "UPDATE game_states SET xp=?,kills=? WHERE user_id=? AND game_id=?"
		err = db.RawExecutor(query, userGameData.Xp, userGameData.Kills, data.UserId, data.RpcParams.GameId)
		if err != nil {
			//resposne the error
			return
		}

		//Increasing the total kills
		var userStats model.UserGameStats
		query = "SELECT * FROM user_game_stats WHERE user_id=?"
		err = db.QueryExecutor(query, &userStats, data.UserId)
		if err != nil {
			//resposne the error
			return
		}

		userStats.TotalKills++
		err = db.UpdateRecord(userStats, data.UserId, "user_id").Error
		if err != nil {
			//resposne the error
			return
		}

		reward, err := ProcessDailyGoal(int64(utils.ZOMBIES_KILLED), data.UserId, &userStats)
		if err != nil {

			fmt.Println("Error in processing the daily goal")
			return

		}

		err = db.UpdateRecord(userStats, data.UserId, "user_id").Error
		if err != nil {
			//resposne the error
			return
		}

		if reward != nil {

			ctx.JSON(200, map[string]interface{}{
				"ResultCode": "1",
				"Data": map[string]interface{}{
					"achievments": reward,
				},
			})
		}
		ctx.JSON(200, map[string]interface{}{
			"ResultCode": "0",
			"Data":       "Data updated sucessfully",
		})

	} else if data.RpcParams.ActionType == float64(utils.MAKE_ZOMBIE) {

		//Player made bot zmombie
		if data.RpcParams.Data.ActorNr == 0 {
			var userGameData model.GameState
			query := "SELECT * FROM game_states WHERE user_id=? AND game_id=?"
			err := db.QueryExecutor(query, &userGameData, data.UserId, data.RpcParams.GameId)
			if err != nil {
				fmt.Println("Error is", err.Error())
				return
			}
			var userStats model.UserGameStats
			query = "SELECT * FROM user_game_stats WHERE user_id=?"
			err = db.QueryExecutor(query, &userStats, data.UserId)
			if err != nil {
				//resposne the error
				return
			}

			//increase the xp of zombie who attacked the player

			query = "UPDATE game_states SET xp=xp+? WHERE game_id=? AND user_id=?"
			err = db.RawExecutor(query, 50, data.RpcParams.GameId, data.UserId)
			if err != nil {

				fmt.Println("Error in processing the game_stats")
				return
			}

			reward, err := ProcessDailyGoal(int64(utils.BECAME_ZOMBIE), data.UserId, &userStats)
			if err != nil {

				fmt.Println("Error in processing the daily goal")
				return

			}

			err = db.UpdateRecord(userStats, data.UserId, "user_id").Error
			if err != nil {
				//resposne the error
				return
			}

			if reward != nil {

				ctx.JSON(200, map[string]interface{}{
					"ResultCode": "1",
					"Data": map[string]interface{}{
						"achievments": reward,
					},
				})
			}
			ctx.JSON(200, map[string]interface{}{
				"ResultCode": "0",
				"Data":       "Data updated sucessfully",
			})

		} else if data.UserId == "" {

			//Bot made player a zombie
			query := "UPDATE game_states SET xp=xp+?,is_zombie=true WHERE game_id=? AND actor_nr=?"

			zombieAttack := 5
			err = db.RawExecutor(query, -zombieAttack, data.RpcParams.GameId, data.RpcParams.Data.ActorNr)
			if err != nil {

				fmt.Println("Error in processing the game_stats")
				return
			}
			ctx.JSON(200, map[string]interface{}{
				"ResultCode": "0",
				"Data":       "Data updated sucessfully",
			})
		} else {
			var userGameData model.GameState
			query := "SELECT * FROM game_states WHERE user_id=? AND game_id=?"
			err := db.QueryExecutor(query, &userGameData, data.UserId, data.RpcParams.GameId)
			if err != nil {
				fmt.Println("Error is", err.Error())
				return
			}
			var userStats model.UserGameStats
			query = "SELECT * FROM user_game_stats WHERE user_id=?"
			err = db.QueryExecutor(query, &userStats, data.UserId)
			if err != nil {
				//resposne the error
				return
			}

			//decrease the xp of the player who get affected by zombie attack

			//set the xp of  the player affected by zombie attack

			query = "UPDATE game_states SET xp=xp+? WHERE game_id=? AND actor_nr=?"

			zombieAttack := 5
			err = db.RawExecutor(query, -zombieAttack, data.RpcParams.GameId, data.RpcParams.Data.ActorNr)
			if err != nil {

				fmt.Println("Error in processing the game_stats")
				return
			}
			fmt.Println("xp decreased for user having actor nr:", data.RpcParams.Data.ActorNr)

			//increase the xp of zombie who attacked the player

			query = "UPDATE game_states SET xp=xp+? WHERE game_id=? AND user_id=?"
			err = db.RawExecutor(query, zombieAttack, data.RpcParams.GameId, data.UserId)
			if err != nil {

				fmt.Println("Error in processing the game_stats")
				return
			}

			reward, err := ProcessDailyGoal(int64(utils.BECAME_ZOMBIE), data.UserId, &userStats)
			if err != nil {

				fmt.Println("Error in processing the daily goal")
				return

			}

			err = db.UpdateRecord(userStats, data.UserId, "user_id").Error
			if err != nil {
				//resposne the error
				return
			}

			if reward != nil {

				ctx.JSON(200, map[string]interface{}{
					"ResultCode": "1",
					"Data": map[string]interface{}{
						"achievments": reward,
					},
				})
			}
			ctx.JSON(200, map[string]interface{}{
				"ResultCode": "0",
				"Data":       "Data updated sucessfully",
			})

			fmt.Println("xp increased for user having user_id:", data.UserId)
		}

	}

}

type GoalReward struct {
	Id              string                    `json:"id"`
	GoalType        int64                     `json:"goalType"`
	CurrentProgress int64                     `json:"currentProgress"`
	TotalProgress   int64                     `json:"totalProgress"`
	CurrencyType    int64                     `json:"currencyType"`
	Price           int64                     `json:"price"`
	RewardData      []response.RewardResponse `json:"rewardData"`
}

func ProcessDailyGoal(goalType int64, userId string, userGameStats *model.UserGameStats) ([]*GoalReward, error) {

	fmt.Println("Process daily goals called")
	var userDailyGoals *model.UserDailyGoals
	var res []*GoalReward
	query := "SELECT * FROM user_daily_goals WHERE user_id=? AND goal_type=?"
	err := db.QueryExecutor(query, &userDailyGoals, userId, goalType)
	if err != nil {
		// response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return nil, err
	}

	// fmt.Println("User Daily Goal", userDailyGoals)

	if userDailyGoals != nil {
		if userDailyGoals.CurrentProgress < userDailyGoals.TotalProgress {
			userDailyGoals.CurrentProgress++
			if userDailyGoals.CurrentProgress == userDailyGoals.TotalProgress {
				rewardData := []response.RewardResponse{
					{
						RewardType: utils.Coins,
						Quantity:   userDailyGoals.Coins,
					},
					{
						RewardType: utils.Gems,
						Quantity:   userDailyGoals.Gems,
					},
				}

				res = append(res, &GoalReward{
					Id:              userDailyGoals.Id,
					GoalType:        userDailyGoals.GoalType,
					CurrentProgress: userDailyGoals.CurrentProgress,
					TotalProgress:   userDailyGoals.TotalProgress,
					CurrencyType:    userDailyGoals.CurrencyType,
					Price:           userDailyGoals.Price,
					RewardData:      rewardData,
				})

				userGameStats.CurrentCoins += userDailyGoals.Coins
				userGameStats.TotalCoins += userDailyGoals.Coins
				userGameStats.TotalGems += userDailyGoals.Gems
				userGameStats.CurrentGems += userDailyGoals.Gems
			}
		}
		err = db.UpdateRecord(&userDailyGoals, userDailyGoals.Id, "id").Error
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

type GameEndRpc struct {

	RpcParams struct {
		GameId string `json:"GameId"`
		Data   struct {
			Bots []struct {
				BotName           string  `json:"botName"`
				KillsAsSurvivor   float64 `json:"killsAsSurvivor"`
				KillsAsZombie     float64 `json:"killsAsZombie"`
				MiniGameCompleted float64 `json:"miniGamesCompleted"`
				Avatar            float64 `json:"avatar"`
				IsZombie          bool    `json:"isZombie"`
			} `json:"bots"`
			//1 for survivors and 2 for zombies
			Win float64 `json:"wins"`
		} `json:"data"`
	}
}
type Rewards struct {
	UserId      string                    `json:"userId"`
	RewardsData []response.RewardResponse `json:"rewards"`
}

	UserId    string
	RpcParams struct {
		GameId     string  `json:"gameId"`
		ActionType float64 `json:"actionType"`
		Data       struct {
			Bots []struct {
				BotName           string `json:"botName"`
				Kills             int    `json:"kills"`
				MiniGameCompleted int    `json:"miniGameCompleted"`
				IsZombie          bool   `json:"isZombie"`
			} `json:"bots"`
			//1 for survivors and 2 for zombies
			Win int `json:"win"`
		} `json:"data"`
	}
}


func GameEnd(ctx *gin.Context) {
	var data GameEndRpc
	body, _ := io.ReadAll(ctx.Request.Body)

	fmt.Println("Body is", string(body))

	err := json.Unmarshal(body, &data)
	if err != nil {

		fmt.Println("Error in unmarshalling the resposne from the hook", err)
		return
	}

	fmt.Println("data", data)
	//Cehck whether the game id is present in the db or not

	var queryRes []struct {
		ActorNr        int    `json:"ActorNr"`
		GameId         string `json:"gameId"`
		Time           int    `json:"time"`
		UserId         string `json:"userId"`
		TotalGames     int    `json:"totalGames"`
		GamesCompleted int    `json:"gamesCompleted"`
		IsZombie       bool   `json:"isZombie"`
		Xp             int64  `json:"xp" `
		Kills          int    `json:"kills"`
		IsRunning      bool   `json:"isRunning"`
		XP             int64  `json:"xp"`
		Level          int64  `json:"level"`
		Energy         int64  `json:"energy"`
		TotalCoins     int64  `json:"totalCoins"`
		CurrentCoins   int64  `json:"currentCoins"`
		TotalGems      int64  `json:"totalGems"`
		CurrentGems    int64  `json:"currentGems"`
		MatchesWon     int64  `json:"matchesWon"`
		MatchesLost    int64  `json:"matchesLost"`
		TotalKills     int64  `json:"totalKills"`
		Username       string `json:"username"`
		Avatar         int    `json:"avatar"`
	}

	query := ` SELECT gs.*, ugs.* ,u.username,u.avatar
	FROM game_states gs 
	JOIN user_game_stats ugs ON gs.user_id = ugs.user_id 
	JOIN users u ON ugs.user_id=u.id
	WHERE gs.game_id = ? AND gs.is_running = 'true'`

	err = db.QueryExecutor(query, &queryRes, data.RpcParams.GameId)
	if err != nil {

		fmt.Println("error in query form the db", err)
		return
	}

	fmt.Println("Query res is", queryRes)
	var rewards []Rewards
	var scoreBoard []response.Scoreboard

	//Get the user game stas for each  user and add the xp and kills gained in a game
	// Update user game stats for each player
	for _, player := range queryRes {

		winCoins := utils.RoundToNearestMultipleOf5(float64(player.Xp * utils.WIN_MULTIPLIER))
		winGems := 1
		winEnergy := 0

		lostCoins := utils.RoundToNearestMultipleOf5(float64(player.Xp * utils.LOST_MULTIPLIER))
		lostGems := 0
		lostEnergy := 0

		player.XP += player.Xp
		player.TotalKills += int64(player.Kills)
		if data.RpcParams.Data.Win == 1 && !player.IsZombie {
			fmt.Println("asvdjhasjdhsffggghhh")
			player.MatchesWon++
			player.CurrentCoins += int64(winCoins)
			player.TotalCoins += int64(winCoins)
			player.CurrentGems += int64(winGems)
			player.TotalGems += int64(winGems)
			player.Energy += int64(winEnergy)

			rewards = append(rewards, Rewards{
				UserId: player.UserId,
				RewardsData: []response.RewardResponse{
					{
						RewardType: utils.Coins,
						Quantity:   int64(winCoins),
					},
					{
						RewardType: utils.Gems,
						Quantity:   int64(winGems),
					},
					{
						RewardType: utils.Energy,
						Quantity:   int64(winEnergy),
					},
				},
			})
		} else if data.RpcParams.Data.Win == 2 && player.IsZombie {
			fmt.Println("askvdjhasjdad")
			player.MatchesWon++
			player.CurrentCoins += int64(winCoins)
			player.TotalCoins += int64(winCoins)
			player.CurrentGems += int64(winGems)
			player.TotalGems += int64(winGems)
			player.Energy += int64(winEnergy)

			rewards = append(rewards, Rewards{
				UserId: player.UserId,
				RewardsData: []response.RewardResponse{
					{
						RewardType: utils.Coins,
						Quantity:   int64(winCoins),
					},
					{
						RewardType: utils.Gems,
						Quantity:   int64(winGems),
					},
					{
						RewardType: utils.Energy,
						Quantity:   int64(winEnergy),
					},
				},
			})
		} else {
			fmt.Println("ajdvjav,,,,,,,jdas")
			player.MatchesLost++
			player.CurrentCoins += int64(lostCoins)
			player.TotalCoins += int64(lostCoins)
			player.CurrentGems += int64(lostGems)
			player.TotalGems += int64(lostGems)
			player.Energy += int64(lostEnergy)

			rewards = append(rewards, Rewards{
				UserId: player.UserId,
				RewardsData: []response.RewardResponse{
					{
						RewardType: utils.Coins,
						Quantity:   int64(lostCoins),
					},
					{
						RewardType: utils.Gems,
						Quantity:   int64(lostGems),
					},
					{
						RewardType: utils.Energy,
						Quantity:   int64(lostEnergy),
					},
				},
			})
		}

		// Check after adding the XP if player level is increased or not
		// If yes, unlock the rewards and update the user level
		var playerLevel model.LevelRewards
		query = "SELECT * FROM level_rewards WHERE xp_required <= ? ORDER BY level_required DESC LIMIT 1"
		err := db.QueryExecutor(query, &playerLevel, player.XP)

		fmt.Println("Error in unmarshalling the resposne from the hook")
		return
	}

	//Get all the players in the game
	var players []model.GameState
	query := `SELECT * FROM game_states WHERE game_id=? AND is_running='true'`
	err = db.QueryExecutor(query, &players, data.RpcParams.GameId)
	if err != nil {
		fmt.Println("Error in getting the players")
		return
	}

	//Get the user game stas for each  user and add the xp and kills gained in a game
	for _, player := range players {
		var userStats model.UserGameStats
		query := "SELECT * FROM user_game_stats WHERE user_id=?"
		err = db.QueryExecutor(query, &userStats, player.UserId)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		userStats.XP += player.Xp
		userStats.TotalKills += int64(player.Kills)
		if data.RpcParams.Data.Win == 1 && !player.IsZombie {
			userStats.MatchesWon++
		} else if data.RpcParams.Data.Win == 2 && player.IsZombie {
			userStats.MatchesWon++
		} else {
			userStats.MatchesLost++
		}

		//Check after adding the xp player level is increased or not if yes then unlock the rewards
		//and update the user level
		var playerLevel model.LevelRewards
		query = "select * from level_rewards WHERE xp_required<=?  ORDER By level  DESC LIMIT 1	"
		err := db.QueryExecutor(query, &playerLevel, userStats.XP)

		if err != nil {
			fmt.Println("Error in getting the level rewards")
			return
		}

		if playerLevel.LevelRequired != player.Level {
			// Mark status of the level_reward for that player as UNCLAIMED till that level
			for i := player.Level; i <= playerLevel.LevelRequired; i++ {
				query = "UPDATE user_level_rewards SET status = ? WHERE user_id = ? AND level = ?"

				err = db.RawExecutor(query, utils.UNCLAIMED, player.UserId, i)
				if err != nil {
					fmt.Println("Error in updating the level rewards")
					return
				}

			}

			player.Level = playerLevel.LevelRequired

		}

		//Query to update user game states
		query = "UPDATE user_game_stats set xp=?,level=?,matches_lost=?,matches_won=?,energy=?,total_coins=?,total_gems=?,current_gems=?,current_coins=?,total_kills=? WHERE user_id=?"
		err = db.RawExecutor(query, player.XP, player.Level, player.MatchesLost, player.MatchesWon, player.Energy, player.TotalCoins, player.TotalGems, player.CurrentGems, player.CurrentCoins, player.TotalKills, player.UserId)
		if err != nil {
			fmt.Println("Error in updating the records of the player", err)
			return
		}
		// err = db.UpdateRecord(&player, player.UserId, "user_id").Error
		// if err != nil {
		// 	fmt.Println("Error in updating the game stats")
		// 	return
		// }
		scoreBoard = append(scoreBoard, response.Scoreboard{
			Name:     player.Username,
			XPGained: player.Xp,
			Avatar:   player.Avatar,
			Level:    player.Level,
		})

	}

	// The rest of the code remains the same
	if data.RpcParams.Data.Bots != nil {
		for i, bot := range data.RpcParams.Data.Bots {
			scoreBoard = append(scoreBoard, response.Scoreboard{
				Name:     bot.BotName,
				XPGained: int64((bot.KillsAsSurvivor * 2) + (bot.KillsAsZombie * 1) + (bot.MiniGameCompleted * 5)),
				Avatar:   int(data.RpcParams.Data.Bots[i].Avatar),
				Level:    1,
			})
		}
	}
	// fmt.Println("Scoreboard is before sorting", scoreBoard)

	//Sorting on the basis of the XP gained
	sort.SliceStable(scoreBoard, func(i, j int) bool {
		return (scoreBoard[i].XPGained) > scoreBoard[j].XPGained
	})

	//Giving rewads to the player on wins

	fmt.Println("ScoreBorad data:", scoreBoard)
	fmt.Println("Rewards data:", rewards)

	ctx.JSON(200, map[string]interface{}{
		"Result": 1,
		"Data": map[string]interface{}{
			"scoreboard": scoreBoard,
			"rewards":    rewards,
		},
	})

}


			}

			userStats.Level = playerLevel.LevelRequired

		}

		err = db.UpdateRecord(&userStats, player.UserId, "user_id").Error
		if err != nil {
			fmt.Println("Error in updating the game stats")
			return
		}

	}

}


