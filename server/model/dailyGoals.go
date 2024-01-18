package model

import (
	"time"

	"gorm.io/gorm"
)

type UserDailyGoals struct {
	Id               string           `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	UserId           string           `json:"userId"`
	GoalType         int64            `json:"goalType"`
	Coins            int64            `json:"coins"`
	Gems             int64            `json:"gems"`
	TotalProgress    int64            `json:"totalProgress"`
	CurrentProgress  int64            `json:"currentProgress"`
	CurrencyType     int64            `json:"currencyType"`
	Price            int64            `json:"price"`
	DailyRewardId    string           `json:"dailyRewardId"`
	DailyGoalRewards DailyGoalRewards `json:"-" gorm:"foreignKey:DailyRewardId;constraint:OnDelete:CASCADE"`
	CreatedAt        time.Time        `json:"-"`
	UpdatedAt        time.Time        `json:"-"`
	DeletedAt        gorm.DeletedAt   `json:"-"`
}

type DailyGoalRewards struct {
	Id        string         `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	Coins     int64          `json:"coins"`
	Gems      int64          `json:"gems"`
	Energy    int64          `json:"energy"`
	Chest     int64          `json:"chest"`
	Claimed   bool           `json:"claimed"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
