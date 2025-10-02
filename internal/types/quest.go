package types

type Quest struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Reward     string `json:"reward"`
	RewardType string `json:"reward_type"` // "orbs", "decor", "other"
	ExpiresAt  string `json:"expires_at"`
}
