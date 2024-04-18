package services

import (
	"errors"
	"fmt"
	"main/server/db"
	"main/server/model"
	dailygoal "main/server/services/daily_goal"
	"main/server/services/rewards"
	"main/server/utils"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

var wg = &sync.WaitGroup{}

func AddDummyUsers() {

	// err := utils.IsPassValid(input.User.Password)
	// if err != nil {
	// 	return
	// }

	dummyUsersCount := 10

	if !utils.TableIsEmpty("users") {

		for i := 1; i <= dummyUsersCount; i++ {
			// wg.
			wg.Add(1)
			go func(i int) {
				encryptedPassword, err := utils.HashPassword("User!234")
				if err != nil {
					return
				}
				userRecord := model.User{
					Email:             "dummy" + strconv.Itoa(i) + "@yopmail.com",
					Password:          *encryptedPassword,
					Username:          "dummy" + strconv.Itoa(i),
					Avatar:            1,
					EmailVerifiedAt:   time.Now(),
					UsernameUpdatedAt: time.Now(),
					EmailVerified:     true,
					DayCount:          1,
				}

				err = db.CreateRecord(&userRecord)
				if err != nil {
					if errors.Is(err, gorm.ErrDuplicatedKey) {

						return
					}
					return
				}

				fmt.Println("User id is", userRecord.Id)

				userGameStats := model.UserGameStats{
					UserId:         userRecord.Id,
					CurrentCoins:   10000,
					CurrentGems:    10000,
					TotalTimeSpent: 0,
					TotalKills:     0,
				}

				userSettings := model.UserSettings{
					UserId:         userRecord.Id,
					Sound:          1,
					Music:          1,
					Vibration:      false,
					VoicePack:      false,
					Notifications:  false,
					FriendRequests: false,
					Language:       "English",
				}

				err = db.CreateRecord(&userSettings)
				if err != nil {
					return
				}

				err = db.CreateRecord(&userGameStats)
				if err != nil {
					return
				}

				var specailOfferId string
				query := "SELECT id FROM special_offers order by created_at ASC limit 1"
				err = db.QueryExecutor(query, &specailOfferId)
				if err != nil {
					return
				}

				//Giving the starter pack to the user after signup
				//For 7 days starter pack will be valid
				userStartPack := model.UserSpecialOffer{
					SpecialOfferId: specailOfferId,
					UserId:         userRecord.Id,
					Purchased:      false,
				}

				err = db.CreateRecord(&userStartPack)
				if err != nil {
					return
				}

				func(userid string) {
					err = rewards.CreateStarterDailyRewards(userid)
					if err != nil {
						return
					}

					// fmt.Println("Daily goal generation")
					err = dailygoal.DailyGoalGeneration(true, &userid)
					if err != nil {
						fmt.Println("Error in daily goal generation", err)
						return
					}
					// fmt.Println("Daily goal generation done")

					// fmt.Println("Level reward generation")
					rewards.GenerateLevelReward(userid)
					// fmt.Println("Level reward generation done")
				}(userRecord.Id)

				wg.Done()
			}(i)

			wg.Wait()

		}
	}
	// fmt.Println("Transaction edy to commit")

}
