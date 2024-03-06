package request

type DeductAmount struct {
	Coins  bool `json:"coins"`
	Gems   bool `json:"gems"`
	Energy bool `json:"energy"`
}
