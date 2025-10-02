package quests

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"time"

	"discord-quest-watcher/internal/types"
	"discord-quest-watcher/internal/webhook"

	"github.com/go-rod/rod"
)

//go:embed extract-quests.js
var extractQuestsScript string

func filterQuests(quests []types.Quest, fn func(types.Quest) bool) []types.Quest {
	var result []types.Quest
	for _, q := range quests {
		if fn(q) {
			result = append(result, q)
		}
	}
	return result
}

func contains(quests []types.Quest, id string) bool {
	for _, q := range quests {
		if q.ID == id {
			return true
		}
	}
	return false
}

func CheckQuests(browser *rod.Browser, webhookURL, rewardFilter string) error {

	allQuests, _ := extractQuests(browser)
	log.Printf("extracted quests: count=%d", len(allQuests))

	// keep only quests we care about
	wantedQuests := allQuests
	if rewardFilter == "orbs" {
		wantedQuests = filterQuests(allQuests, func(q types.Quest) bool { return q.RewardType == "orbs" })
	}
	log.Printf("filtered quests: count=%d filter=%s", len(wantedQuests), rewardFilter)

	// find which ones are actually new
	previousQuests := questStorage(nil)
	newQuests := filterQuests(wantedQuests, func(current types.Quest) bool {
		return !contains(previousQuests, current.ID)
	})
	log.Printf("new quests: count=%d", len(newQuests))

	// notify about new ones
	if len(newQuests) > 0 {
		log.Printf("sending notifications: count=%d", len(newQuests))
		webhook.Send(webhookURL, newQuests)
	}

	// remember what we found
	log.Printf("saving quests: count=%d", len(wantedQuests))
	questStorage(wantedQuests)
	return nil
}

func extractQuests(browser *rod.Browser) ([]types.Quest, error) {
	page := browser.MustPage("https://discord.com/discovery/quests").MustWaitLoad()
	time.Sleep(10 * time.Second) // wait for react to load

	result, err := page.Eval("() => {" + extractQuestsScript + "}")

	if err != nil {
		return nil, err
	}

	var quests []types.Quest
	return quests, json.Unmarshal([]byte(result.Value.String()), &quests)
}

func questStorage(quests []types.Quest) []types.Quest {
	const file = "/data/known-quests.json"

	// load existing if reading
	if quests == nil {
		var loaded []types.Quest
		if data, _ := os.ReadFile(file); data != nil {
			json.Unmarshal(data, &loaded)
		}
		return loaded
	}

	// save if writing
	os.MkdirAll("/data", 0755)
	if data, _ := json.MarshalIndent(quests, "", "  "); data != nil {
		os.WriteFile(file, data, 0644)
	}
	return quests
}
