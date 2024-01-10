package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/server/db"
	"main/server/model"
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
		// {"part_customizations", "server/dummyData/partCustomization.json", &[]model.PartCustomization{}},
		// {"default_customisations", "server/dummyData/defaultCustomization.json", &[]model.DefaultCustomisation{}},
		// {"race_types", "server/dummyData/raceTypes.json", &[]model.RaceTypes{}},
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

	default:
		log.Fatal("Invalid modelType provided")
	}

}
