# discord-quest-watcher

receive notifications when new discord quests are released. filter for orb quests or monitor all quest types.

minimalistic, single-dependency go app that reliably logs in via a user token, bypassing captchas and rate limits. fully self-hostable and private.

[how to get your discord token](https://gist.github.com/MarvNC/e601f3603df22f36ebd3102c501116c6#file-get-discord-token-from-browser-md)

> [!WARNING]  
> this app uses your user token which technically breaks Discord ToS, so use at your own risk.

## features

- filters for orb rewards or all quest types
- sends discord webhook notifications when new quests appear
- docker support
- checks for new quests every 30 minutes

## environment variables

| variable              | required | default | description                                                                 |
|-----------------------|----------|---------|-----------------------------------------------------------------------------|
| `TOKEN`               | yes      | —       | discord authentication token                                                |
| `DISCORD_WEBHOOK_URL` | yes      | —       | webhook URL used for sending notifications                                  |
| `REWARD_FILTER`       | no       | `all`   | filter for rewards: `orbs` (only orbs) or `all` (include all rewards)       |

## usage

### docker (recommended)

```shell
docker run -d \
           --name discord-quest-watcher \
           --restart=unless-stopped \
           -e TOKEN=your-token \
           -e DISCORD_WEBHOOK_URL=your-webhook-url \
           -e REWARD_FILTER=orbs \
           ghcr.io/xhos/discord-quest-watcher:latest
```

### local

```shell
export TOKEN="your_discord_token"
export DISCORD_WEBHOOK_URL="your_webhook_url"
export REWARD_FILTER="orbs"

go run .
```

## how it works

it authenticates with your discord token and navigates to the quests page in a headless chromium browser. it then scrapes the page, extracting the quest data, compares against previously seen quests, and sends webhook notifications for any new entries. state is tracked in /data/known-quests.json with 30-minute check intervals.
