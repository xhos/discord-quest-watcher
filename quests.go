package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-rod/rod"
)

func checkQuests(token, webhook, rewardFilter string) error {
	browser, _ := createBrowser()
	defer browser.MustClose()

	authenticateWithToken(browser, token)

	allQuests, _ := extractQuests(browser)
	log.Info("extracted quests", "count", len(allQuests))

	// keep only quests we care about
	wantedQuests := allQuests
	if rewardFilter == "orbs" {
		wantedQuests = []Quest{}
		for _, quest := range allQuests {
			if quest.RewardType == "orbs" {
				wantedQuests = append(wantedQuests, quest)
			}
		}
	}
	log.Info("filtered quests", "count", len(wantedQuests), "filter", rewardFilter)

	// find which ones are actually new
	previousQuests := questStorage(nil)
	newQuests := []Quest{}
	for _, current := range wantedQuests {
		isAlreadyKnown := false
		for _, previous := range previousQuests {
			if current.ID == previous.ID {
				isAlreadyKnown = true
				break
			}
		}
		if !isAlreadyKnown {
			newQuests = append(newQuests, current)
		}
	}
	log.Info("new quests", "count", len(newQuests))

	// notify about new ones
	if len(newQuests) > 0 {
		log.Info("sending notifications", "count", len(newQuests))
		sendNotifications(webhook, newQuests)
	}

	// remember what we found
	log.Info("saving quests", "count", len(wantedQuests))
	questStorage(wantedQuests)
	return nil
}

func extractQuests(browser *rod.Browser) ([]Quest, error) {
	page := browser.MustPage("https://discord.com/discovery/quests").MustWaitLoad()
	time.Sleep(10 * time.Second) // wait for react to load

	result, err := page.Eval(`() => {
		const quests = [];
		document.querySelectorAll('[id^="quest-tile-"]').forEach(tile => {
			const name = tile.querySelector('[class*="questName"]')?.textContent?.trim();
			const reward = tile.querySelector('[class*="header"]')?.textContent?.replace('Claim ', '')?.trim();
			const allText = tile.textContent;
			const expires = allText.match(/Ends (\d{2}\/\d{2})/)?.[1];

			if (name && reward && expires && !allText.includes('Quest ended')) {
				let rewardType = 'other';
				if (reward.toLowerCase().includes('orb')) rewardType = 'orbs';
				else if (reward.toLowerCase().includes('avatar decoration')) rewardType = 'decor';

				quests.push({
					id: tile.id,
					name: name,
					reward: reward,
					reward_type: rewardType,
					expires_at: expires
				});
			}
		});
		return JSON.stringify(quests);
	}`)

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
		data, _ := os.ReadFile(file)
		var loaded []Quest
		json.Unmarshal(data, &loaded)
		return loaded
	}

	// save if writing
	os.MkdirAll("/data", 0755)
	data, _ := json.MarshalIndent(quests, "", "  ")
	os.WriteFile(file, data, 0644)
	return quests
}
