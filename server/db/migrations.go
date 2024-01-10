package db

import (
	"fmt"
	"main/server/model"

	"gorm.io/gorm"
)

func AutoMigrateDatabase(db *gorm.DB) {

	var dbVersion model.DbVersion
	err := db.First(&dbVersion).Error
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println("db version is:", dbVersion.Version)
	if dbVersion.Version < 1 {
		err := db.AutoMigrate(&model.User{}, &model.UserGameStats{}, &model.ResetSession{}, &model.UserSettings{}, &model.Session{}, &model.LevelRewards{}, &model.UserLevelRewards{}, model.UserBadges{}, model.DailyRewards{}, model.UserDailyRewards{})
		if err != nil {
			panic(err)
		}
		db.Create(&model.DbVersion{
			Version: 1,
		})
		dbVersion.Version = 1
	}

}
