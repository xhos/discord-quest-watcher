package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

func sendNotifications(webhook string, quests []Quest) {
	for _, quest := range quests {
		color := 0x99AAB5 // gray
		switch quest.RewardType {
		case "orbs":
			color = 0x5865F2 // blue
		case "decor":
			color = 0x57F287 // green
		}

		embed := map[string]any{
			"title":       fmt.Sprintf("🔮 New %s Quest!", quest.RewardType),
			"description": fmt.Sprintf("**%s**\n%s\nExpires: %s", quest.Name, quest.Reward, quest.ExpiresAt),
			"color":       color,
		}

		payload := map[string]any{
			"embeds": []any{embed},
		}

		data, _ := json.Marshal(payload)
		resp, _ := http.Post(webhook, "application/json", bytes.NewBuffer(data))
		resp.Body.Close()

		log.Info("sent notification", "quest", quest.Name)
	}
}
