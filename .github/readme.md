# discord-quest-watcher

<img width="400" height="248" alt="image" src="https://github.com/user-attachments/assets/3771e972-9898-450a-b1c5-32aff759fbc8" />

receive notifications when new discord quests are released. filter for orb quests or monitor all quest types.

minimalistic, single-dependency go app that reliably logs in via a user token, bypassing captchas and rate limits. fully self-hostable and private.

> why?

I couldn't find any existing tools that would reliably ping me when a new orb quest drops, so i threw this little tool together in an evening.

> but why would one care about those quests in the first place?

quest give orb -> orb give free shiny things. 

no but seriously, I like the look of some discord user decor, but there's no way im paying for it. oh and by the way, you dont have to actually *complete* the quest, this script does it for you: https://gist.github.com/aamiaa/204cd9d42013ded9faf646fae7f89fbb

or you can use a vencord/equicord/whateverispopularnowcord plugin. 

## environment variables

| variable              | required | default | description                                                                 |
|-----------------------|----------|---------|-----------------------------------------------------------------------------|
| `TOKEN`               | yes      | —       | discord authentication token                                                |
| `DISCORD_WEBHOOK_URL` | yes      | —       | webhook URL used for sending notifications                                  |
| `REWARD_FILTER`       | no       | `all`   | filter for rewards: `orbs` (only orbs) or `all` (include all rewards)       |
| `FETCH_INTERVAL`      | no       | `30`    | interval in minutes between quest checks (must be a positive integer)       |
| `RUN_ONCE`            | no       | `false` | if `true`, the application will run once and then exit                    |
| `WEBHOOK_MESSAGE`     | no       | —       | additional text appended to notifications (e.g., role pings)                |

### custom webhook messages

notifications are sent as discord embeds. you can add a text message above the embed using the `WEBHOOK_MESSAGE` environment variable.

example: ping a role

```shell
WEBHOOK_MESSAGE=<@&`1234567890123456789`>
```

## usage

> [!WARNING]  
> this app uses your user token which technically breaks Discord ToS, so use at your own risk.

[how to get your user token](https://gist.github.com/MarvNC/e601f3603df22f36ebd3102c501116c6#file-get-discord-token-from-browser-md)

```shell
docker run -d \
           --name discord-quest-watcher \
           --restart=unless-stopped \
           -e TOKEN=your-token \
           -e DISCORD_WEBHOOK_URL=your-webhook-url \
           -e REWARD_FILTER=orbs \
           -e WEBHOOK_MESSAGE=<@&role-id> \
           ghcr.io/xhos/discord-quest-watcher:latest
```

## how it works

it authenticates with your discord token and navigates to the quests page in a headless chromium browser. it then scrapes the page, extracting the quest data, compares against previously seen quests, and sends webhook notifications for any new entries. state is tracked in /data/known-quests.json with 30-minute check intervals.

## releasing

to publish a new docker image, push a semver tag:

```shell
git tag v1.0.0
git push origin v1.0.0
```

this creates images tagged as `1.0.0`, `1.0`, `1`, and `latest`.
