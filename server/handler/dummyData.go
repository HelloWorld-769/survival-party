package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/services"
	"main/server/utils"
)

func ReadJSONFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func AddDummyDataHandler() {
	dataFiles := []struct {
		tableName string
		filePath  string
		dataPtr   interface{}
	}{
		{"level_rewards", "server/dummyData/rewards.json", &[]model.LevelRewards{}},
		{"daily_rewards", "server/dummyData/dailyRewards.json", &[]model.DailyRewards{}},
		{"shops", "server/dummyData/shop.json", &[]model.Shop{}},
		{"special_offers", "server/dummyData/specialOffer.json", &[]model.SpecialOffer{}},
		// {"race_rewards", "server/dummyData/rewards.json", &[]model.RaceRewards{}},
		// {"reward_data", "server/dummyData/collectables_rewards.json", &[]model.RewardData{}},
		// {"arena_level_perks", "server/dummyData/arenaPerks.json", &[]model.ArenaLevelPerks{}},
	}

	fmt.Println("dummy handler")

	for _, dataFile := range dataFiles {
		if !utils.TableIsEmpty(dataFile.tableName) {
			addtoDb(dataFile.filePath, dataFile.dataPtr)
		}
	}

}
func AddDummyUsers() {
	filePath := "server/dummyData/dummyUsers.json"
	if !utils.TableIsEmpty("users") {
		Data, err := ReadJSONFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		var modelType []request.SigupRequest
		err = json.Unmarshal(Data, &modelType)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("   ", modelType)

		for _, item := range modelType {
			services.AddDummyUsers(item)
		}
	}

}
func addtoDb(filePath string, modelType interface{}) {

	Data, err := ReadJSONFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(Data, &modelType)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("adding data to DB function")
	switch slice := modelType.(type) {
	case *[]model.LevelRewards:
		for _, item := range *slice {
			db.CreateRecord(&item)
		}
	case *[]model.DailyRewards:
		for _, item := range *slice {
			db.CreateRecord(&item)
		}
	case *[]model.Shop:
		for _, item := range *slice {
			db.CreateRecord(&item)
		}
	case *[]model.SpecialOffer:
		for _, item := range *slice {
			db.CreateRecord(&item)
		}
	default:
		log.Fatal("Invalid modelType provided")
	}

}
