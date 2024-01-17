package model

import (
	"time"

	"gorm.io/gorm"
)

type LevelRewards struct {
	RewardId      string    `json:"rewardId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	RewardType    int64     `json:"rewardType"`
	Quantity      int64     `json:"quantity"`
	LevelRequired int64     `json:"levelRequired"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt
}

type UserLevelRewards struct {
	UserId    string       `json:"userId"`
	User      User         `json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	RewardId  string       `json:"rewardId"`
	Rewards   LevelRewards `json:"-" gorm:"foreignKey:RewardId;constraint:OnDelete:CASCADE"`
	Claimed   bool         `json:"claimed"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

type DailyRewards struct {
	RewardId  string    `json:"rewardId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	DayCount  int64     `json:"dayCount"`
	Coins     int64     `json:"coins"`
	Gems      int64     `json:"gems"`
	Energy    int64     `json:"energy"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

type UserDailyRewards struct {
	Id         string `json:"Id" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	UserId     string `json:"userId" gorm:"constraint:unique"`
	User       User   `json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Status     int64  `json:"status"`
	DayCount   int64  `json:"dayCount"`
	RewardType int64  `json:"rewardType"`
	Gain       int64  `json:"gain" `
	AssetName  string `json:"assetName,omitempty" `
	Name       string `json:"name ,omitempty"`
	ChestType  int64  `json:"chestType,omitempty"`
}
