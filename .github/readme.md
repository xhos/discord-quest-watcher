# discord-quest-watcher

receive notifications when new discord quests are released. filter for orb quests or monitor all quest types.

minimalistic, single-dependency go app that reliably logs in via a user token, bypassing captchas and rate limits. fully self-hostable and private.

> why?

I couldn't find any existing tools that would reliably ping me when a new orb quest drops, so i threw this little tool together in an evening.

> but why would one care about those quests in the first place?

quest give orb -> orb allow monkey brain get free shiny thing. 

no but seriously, I like the look of some discord user decor, but there's no way im paying for it. oh and by the way, you dont have to actually *complete* the quest, this script does it for you: https://gist.github.com/aamiaa/204cd9d42013ded9faf646fae7f89fbb

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
| `FETCH_INTERVAL`      | no       | `30`    | interval in minutes between quest checks (must be a positive integer)       |
| `RUN_ONCE`            | no       | `false` | if `true`, the application will run once and then exit                    |

## usage

> [!WARNING]  
> this app uses your user token which technically breaks Discord ToS, so use at your own risk.

[how to get your user token](https://gist.github.com/MarvNC/e601f3603df22f36ebd3102c501116c6#file-get-discord-token-from-browser-md)

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
