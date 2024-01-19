package request

type GoalRequest struct {
	Id string `json:"id"`
}

type UpdateGoalReq struct {
	KillAsZomb int64 `json:"killsAsZombie"`
	KillAsSur  int64 `json:"killsAsSurvivor"`
	GamesComp  int64 `json:"miniGamesCompleted"`
	IsZombie   bool  `json:"isZombie"`
	IsDead     bool  `json:"isDead"`
}
