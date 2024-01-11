package model

import (
	"time"

	"gorm.io/gorm"
)

type SpecialOffer struct {
	Id           string    `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	ProductId    string    `json:"productId"`
	Coins        int64     `json:"coins"`
	Gems         int64     `json:"gems"`
	Inventory    int64     `json:"inventory"`
	CurrencyType int       `json:"currencyType"`
	Price        int64     `json:"price"`
	IsAvailable  bool      `json:"isAvailable"`
	ExpireAt     time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt
}
