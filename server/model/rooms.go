package model

import (
	"time"

	"gorm.io/gorm"
)

type Rooms struct {
	RoomId          string         `json:"roomId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	Capacity        int64          `json:"capaccity"`
	CurrentCapacity int64          `json:"current_capacity"`
	UserId          string         `json:"userId"`
	User            User           `json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Is_Open         bool           `json:"is_open"`
	CreatedAt       time.Time      `json:"-"`
	UpdatedAt       time.Time      `json:"-"`
	DeletedAt       gorm.DeletedAt `json:"-"`
}

type UsersInRooms struct {
	RoomId    string         `json:"roomId"`
	Room      Rooms          `json:"-" gorm:"foreignKey:RoomId;constraint:OnDelete:CASCADE"`
	UserId    string         `json:"userId"`
	User      User           `json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Actor_Nr  int            `json:"actor_nr"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
