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
