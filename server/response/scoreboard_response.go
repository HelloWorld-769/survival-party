package response

type Scoreboard struct {
	Name     string `json:"name"`
	XPGained int64  `json:"experienceGained"`
	Avatar   int    `json:"avatar"`
	Level    int64  `json:"level"`
}
