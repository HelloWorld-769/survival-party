package model

import "time"

type GameState struct {
	ActorNr int    `json:"ActorNr"`
	GameId  string `json:"gameId"`
	Time    int    `json:"time"`
	UserId  string `json:"userId"`
	// User           User      `json:"-" gorm:"references:Id;constraint:OnDelete:CASCADE"`
	TotalGames     int       `json:"totalGames"`
	GamesCompleted int       `json:"gamesCompleted"`
	IsDead         bool      `json:"isDead"`
	IsZombie       bool      `json:"isZombie"`
	Xp             int64     `json:"xp" `
	Kills          int       `json:"kills"`
	IsRunning      bool      `json:"isRunning"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
