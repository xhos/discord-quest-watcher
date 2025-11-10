package webhook

import (
	"bytes"
	"discord-quest-watcher/internal/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Send(webhook string, quests []types.Quest) {
	customMsg := os.Getenv("WEBHOOK_MESSAGE")
	colors := map[string]int{"orbs": 0x5865F2, "decor": 0x57F287}

	for _, quest := range quests {
		color := 0x99AAB5 // default gray
		if c, ok := colors[quest.RewardType]; ok {
			color = c
		}

		payload := map[string]any{
			"embeds": []any{map[string]any{
				"title":       fmt.Sprintf("🔮 New %s Quest!", quest.RewardType),
				"description": fmt.Sprintf("**%s**\n%s\nExpires: %s", quest.Name, quest.Reward, quest.ExpiresAt),
				"color":       color,
			}},
		}

		if customMsg != "" {
			payload["content"] = customMsg
		}

		if data, _ := json.Marshal(payload); data != nil {
			if resp, _ := http.Post(webhook, "application/json", bytes.NewBuffer(data)); resp != nil {
				resp.Body.Close()
			}
		}

		log.Printf("sent notification for quest: %s", quest.Name)
	}
}
