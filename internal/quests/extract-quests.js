const quests = [];
document.querySelectorAll('[id^="quest-tile-"]').forEach(tile => {
	const name = tile.querySelector('[class*="questName"]')?.textContent?.trim();
	const reward = tile.querySelector('[class*="header"]')?.textContent?.replace('Claim ', '')?.trim();
	const allText = tile.textContent;
	const expiresMatch = allText.match(/Ends (\d{2}\/\d{2})/)?.[1];

	if (name && reward && expiresMatch && !allText.includes('Quest ended')) {
		let rewardType = 'other';
		if (reward.toLowerCase().includes('orb')) rewardType = 'orbs';
		else if (reward.toLowerCase().includes('avatar decoration')) rewardType = 'decor';

		const [day, month] = expiresMatch.split('/').map(Number);
		const now = new Date();
		let year = now.getFullYear();
		const expiryDate = new Date(Date.UTC(year, month - 1, day, 23, 59, 59));

		if (expiryDate < now) {
			expiryDate.setFullYear(year + 1);
		}

		const expiresTimestamp = Math.floor(expiryDate.getTime() / 1000);

		quests.push({
			id: tile.id,
			name: name,
			reward: reward,
			reward_type: rewardType,
			expires_at: expiresTimestamp.toString()
		});
	}
});
return JSON.stringify(quests);
