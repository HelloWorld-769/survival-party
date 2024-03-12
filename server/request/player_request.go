package request

// "main/server/validation"

type UpdatePlayer struct {
	Username string `json:"username"`
	Avatar   int64  `json:"avatar"`
}

type UpdatePlayerSettingsRequest struct {
	Sound          int64   `json:"sound"`
	Music          int64   `json:"music"`
	JoystickSize   float64 `json:"joystick_size"`
	Vibration      bool    `json:"vibration"`
	Language       string  `json:"language"`
	FriendRequests bool    `json:"friendRequests"`
	Notifications  bool    `json:"notifications"`
}

// type UserSettings struct {
// 	Id             string    `json:"Id" gorm:"default:uuid_generate_v4();unique;primaryKey"`
// 	UserId         string    `json:"userId"`
// 	User           User      `json:"-" gorm:"references:Id;constraint:OnDelete:CASCADE"`
// 	Sound          int64     `json:"sound"`
// 	Music          int64     `json:"music"`
// 	Vibration      bool      `json:"vibration"`
// 	VoicePack      bool      `json:"voicePack"`
// 	Notifications  bool      `json:"notifications"`
// 	FriendRequests bool      `json:"friendRequests"`
// 	JoystickSize   float64   `json:"joystickSize"`
// 	Language       string    `json:"language"`
// 	CreatedAt      time.Time `json:"created_at"`
// 	UpdatedAt      time.Time `json:"updated_at"`
// 	DeletedAt      gorm.DeletedAt
// }

// vibration, voicepack, frienRequests.language sound ,music, joystickSize

type PlayerLevelRewardCollectRequest struct {
	Level int64 `json:"level"`
}

type CollectDailyRewardsRequest struct {
	RewardId string `json:"reward_id"`
}

type DailyRewardMuti struct {
	Type int `json:"type"`
}
