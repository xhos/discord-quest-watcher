package main

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

type Quest struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Reward     string `json:"reward"`
	RewardType string `json:"reward_type"` // "orbs", "decor", "other"
	ExpiresAt  string `json:"expires_at"`
}

func main() {
	token := os.Getenv("TOKEN")
	webhook := os.Getenv("DISCORD_WEBHOOK_URL")
	rewardFilter := getEnv("REWARD_FILTER", "all") // "all" or "orbs"

	if token == "" || webhook == "" {
		log.Fatal("TOKEN and DISCORD_WEBHOOK_URL required")
	}

	log.Info("starting Discord quest monitor", "reward_filter", rewardFilter)

	for {
		log.Info("checking for new quests")
		if err := checkQuests(token, webhook, rewardFilter); err != nil {
			log.Error("quest check failed", "error", err)
		}
		time.Sleep(30 * time.Minute)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
