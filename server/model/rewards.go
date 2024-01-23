package model

import (
	"time"

	"gorm.io/gorm"
)

type LevelRewards struct {
	RewardId      string         `json:"rewardId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	RewardType    int64          `json:"rewardType"`
	Quantity      int64          `json:"quantity"`
	LevelRequired int64          `json:"levelRequired"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `json:"-"`
}

type UserLevelRewards struct {
	UserId     string         `json:"userId"`
	User       User           `json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	RewardType int64          `json:"rewardType"`
	Quantity   int64          `json:"quantity"`
	Status     int64          `json:"status"`
	Level      int64          `json:"level"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-"`
}

type DailyRewards struct {
	RewardId  string         `json:"rewardId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	DayCount  int64          `json:"dayCount"`
	Coins     int64          `json:"coins"`
	Gems      int64          `json:"gems"`
	Energy    int64          `json:"energy"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type UserDailyRewards struct {
	Id         string `json:"Id" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	UserId     string `json:"userId" gorm:"constraint:unique"`
	User       User   `json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Status     int64  `json:"status"`
	DayCount   int64  `json:"dayCount"`
	RewardType int64  `json:"rewardType"`
	Gain       int64  `json:"gain" `
	AssetName  string `json:"assetName" `
	Name       string `json:"name"`
	ChestType  int64  `json:"chestType"`
}
