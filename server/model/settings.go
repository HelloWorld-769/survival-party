package model

import (
	"time"

	"gorm.io/gorm"
)

type UserSettings struct {
	Id             string         `json:"Id" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	UserId         string         `json:"userId"`
	User           User           `json:"-" gorm:"references:Id;constraint:OnDelete:CASCADE"`
	Sound          float64        `json:"sound"`
	Music          float64        `json:"music"`
	Vibration      bool           `json:"vibration"`
	VoicePack      bool           `json:"voicePack"`
	Notifications  bool           `json:"notifications"`
	FriendRequests bool           `json:"friendRequests"`
	JoystickSize   float64        `json:"joystickSize"`
	Language       string         `json:"language"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-"`
}
