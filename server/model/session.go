package model

import (
	"time"

	"gorm.io/gorm"
)

// DB model to store session information
type Session struct {
	SessionId string         `json:"sessionId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	UserId    string         `json:"userId"`
	Token     string         `json:"token"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type ResetSession struct {
	SessionId string         `json:"sessionId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	UserEmail string         `json:"userEmail"`
	Otp       int64          `json:"otp"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
