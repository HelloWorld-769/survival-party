package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id              string    `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	Email           string    `json:"email"  gorm:"unique"`
	Emailverified   bool      `json:"emailverified"`
	Password        string    `json:"password"`
	Username        string    `json:"username"  gorm:"unique"`
	Avatar          int64     `json:"avatar"`
	Coins           int64     `json:"coins"`
	Gems            int64     `json:"gems"`
	Energy          int64     `json:"energy"`
	TotalCoins      int64     `json:"totalCoins"`
	CurrentCoins    int64     `json:"currentCoins"`
	TotalGems       int64     `json:"totalGems"`
	CurrentGems     int64     `json:"currentGems"`
	CurrentTrophies int64     `json:"currentTrophies"`
	HighestTrophies int64     `json:"highestTrophies"`
	XP              int64     `json:"xp"`
	Level           int64     `json:"level"`
	SocialId        string    `json:"socialId"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt
}

type UseGameStats struct {
	UserId         string    `json:"userId"`
	User           User      `json:"-" gorm:"refrences:UserId;constraint:OnDelete:CASCADE"`
	MatchesPlayed  int64     `json:"matchesPlayed"`
	MatchesWon     int64     `json:"matchesWon"`
	TotalTimeSpent time.Time `json:"totalTimeSpent"`
	TotalKills     int64     `json:"totalKills"`
}
