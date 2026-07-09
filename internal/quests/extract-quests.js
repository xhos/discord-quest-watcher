const quests = [];
const tiles = document.querySelectorAll('[id^="quest-tile-"]');

tiles.forEach((tile, index) => {
    const allText = tile.textContent;

    const name = tile.querySelector('[class*="questName"]')?.textContent?.trim();
    const reward = tile.querySelector('[class*="header"]')?.textContent?.replace('Claim ', '')?.trim();
    const expiresMatch = allText.match(/Ends (\d{1,2}\/\d{1,2})/)?.[1];
    const containsEnded = allText.includes('Quest ended');
    
    if (!name || !reward || !expiresMatch || containsEnded) return;

    let rewardType = 'other';
    if (reward.toLowerCase().includes('orb')) rewardType = 'orbs';
    else if (reward.toLowerCase().includes('avatar decoration')) rewardType = 'decor';

    const [first, second] = expiresMatch.split('/').map(Number);
    const now = new Date();
    let year = now.getFullYear();

    // date parsing sucks, but this might help:

    // if first > 12, it must be DD/MM, else assume MM/DD
    let month, day;
    if (first > 12) {
        // must be DD/MM format
        day = first;
        month = second;
    } else if (second > 12) {
        // must be MM/DD format
        month = first;
        day = second;
    } else {
        // try both and pick the one that makes more sense
        // assume mm/dd first
        month = first;
        day = second;
    }

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
});

return JSON.stringify(quests);
