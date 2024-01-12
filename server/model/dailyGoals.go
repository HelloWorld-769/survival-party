package model

type UserDailyGoals struct {
	UserId          string `json:"userId"`
	GoalType        int64  `json:"goalType"`
	Coins           int64  `json:"coins"`
	Gems            int64  `json:"gems"`
	TotalProgress   int64  `json:"totalProgress"`
	CurrentProgress int64  `json:"currentProgress"`
	CurrencyType    int64  `json:"currencyType"`
	Price           int64  `json:"price"`
}
