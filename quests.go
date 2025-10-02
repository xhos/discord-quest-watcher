package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
)

//go:embed scripts/extract-quests.js
var extractQuestsScript string

func filterQuests(quests []Quest, fn func(Quest) bool) []Quest {
	var result []Quest
	for _, q := range quests {
		if fn(q) {
			result = append(result, q)
		}
	}
	return result
}

func contains(quests []Quest, id string) bool {
	for _, q := range quests {
		if q.ID == id {
			return true
		}
	}
	return false
}

func checkQuests(token, webhook, rewardFilter string) error {
	browser, _ := createBrowser()
	defer browser.MustClose()

	authenticateWithToken(browser, token)

	allQuests, _ := extractQuests(browser)
	log.Printf("extracted quests: count=%d", len(allQuests))

	// keep only quests we care about
	wantedQuests := allQuests
	if rewardFilter == "orbs" {
		wantedQuests = filterQuests(allQuests, func(q Quest) bool { return q.RewardType == "orbs" })
	}
	log.Printf("filtered quests: count=%d filter=%s", len(wantedQuests), rewardFilter)

	// find which ones are actually new
	previousQuests := questStorage(nil)
	newQuests := filterQuests(wantedQuests, func(current Quest) bool {
		return !contains(previousQuests, current.ID)
	})
	log.Printf("new quests: count=%d", len(newQuests))

	// notify about new ones
	if len(newQuests) > 0 {
		log.Printf("sending notifications: count=%d", len(newQuests))
		sendNotifications(webhook, newQuests)
	}

	// remember what we found
	log.Printf("saving quests: count=%d", len(wantedQuests))
	questStorage(wantedQuests)
	return nil
}

func extractQuests(browser *rod.Browser) ([]Quest, error) {
	page := browser.MustPage("https://discord.com/discovery/quests").MustWaitLoad()
	time.Sleep(10 * time.Second) // wait for react to load

	result, err := page.Eval("() => {" + extractQuestsScript + "}")

	if err != nil {
		return nil, err
	}

	var quests []Quest
	return quests, json.Unmarshal([]byte(result.Value.String()), &quests)
}

func questStorage(quests []Quest) []Quest {
	const file = "/data/known-quests.json"

	// load existing if reading
	if quests == nil {
		var loaded []Quest
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
