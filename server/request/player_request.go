package request

// "main/server/validation"

type UpdatePlayer struct {
	Username string `json:"username"`
	Avatar   int64  `json:"avatar"`
}

type UpdatePlayerSettingsRequest struct {
	Settings struct {
		Sound        int64   `json:"sound"`
		Music        int64   `json:"music"`
		JoystickSize float64 `json:"joystick_size"`
		Vibration    bool    `json:"vibration"`
	} `json:"setting"`
}

type PlayerLevelRewardCollectRequest struct {
	Level int64 `json:"level"`
}

type CollectDailyRewardsRequest struct {
	RewardId string `json:"reward_id"`
}
