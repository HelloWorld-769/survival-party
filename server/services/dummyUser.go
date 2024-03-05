package services

import (
	"errors"
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	dailygoal "main/server/services/daily_goal"
	"main/server/services/rewards"
	"main/server/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

func AddDummyUsers(input request.SigupRequest) {

	err := utils.IsPassValid(input.User.Password)
	if err != nil {
		return
	}

	encryptedPassword, err := utils.HashPassword(input.User.Password)
	if err != nil {
		return
	}

	userRecord := model.User{
		Email:             input.User.Email,
		Password:          *encryptedPassword,
		Username:          strings.ToLower(input.User.Username),
		Avatar:            input.User.Avatar,
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
		Language:       "english",
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

	err = rewards.CreateStarterDailyRewards(userRecord.Id)
	if err != nil {
		return
	}

	dailygoal.DailyGoalGeneration(true, &userRecord.Id)

	rewards.GenerateLevelReward(userRecord.Id)

	fmt.Println("Transaction edy to commit")

}
