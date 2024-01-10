package model

import (
	"time"

	"gorm.io/gorm"
)

type LevelRewards struct {
	RewardId      string `json:"rewardId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	RewardType    int64  `json:"rewardType"`
	Quantity      int64  `json:"quantity"`
	LevelRequired int64  `json:"levelRequired"`
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
}
