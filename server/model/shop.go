package model

import (
	"time"

	"gorm.io/gorm"
)

type Shop struct {
	Id           string         `json:"-" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	ProductId    string         `json:"productId"`
	RewardType   int64          `json:"rewardType"`
	CurrencyType int64          `json:"currencyType"`
	Quantity     int64          `json:"quantity"`
	Price        int64          `json:"price"`
	Popup        bool           `json:"popup"`
	IsAvailable  bool           `json:"isAvailable"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}
