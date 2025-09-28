package main

import (
	"log"
	"os"
	"time"
)

type Quest struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Reward     string `json:"reward"`
	RewardType string `json:"reward_type"` // "orbs", "decor", "other"
	ExpiresAt  string `json:"expires_at"`
}

func main() {
	token, webhook := os.Getenv("TOKEN"), os.Getenv("DISCORD_WEBHOOK_URL")
	rewardFilter := func() string {
		if value := os.Getenv("REWARD_FILTER"); value != "" {
			return value
		}
		return "all"
	}()

	if token == "" || webhook == "" {
		log.Fatal("TOKEN and DISCORD_WEBHOOK_URL required")
	}

	log.Printf("starting Discord quest monitor with reward_filter=%s", rewardFilter)

	for {
		log.Println("checking for new quests")
		if err := checkQuests(token, webhook, rewardFilter); err != nil {
			log.Printf("quest check failed: %v", err)
		}
		time.Sleep(30 * time.Minute)
	}
}
