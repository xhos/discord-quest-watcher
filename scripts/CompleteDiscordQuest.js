/*
 * Complete Recent Discord Quest
 * Source: https://gist.github.com/aamiaa/204cd9d42013ded9faf646fae7f89fbb
 * Author: aamiaa
 * License: GPL-3.0
 *
 * This script automates the completion of Discord quests by spoofing game activity,
 * video watching, or streaming depending on the quest requirements.
 *
 * Last synced: 2025-10-02 20:14:13 UTC
 */

delete window.$;
let wpRequire = webpackChunkdiscord_app.push([[Symbol()], {}, r => r]);
webpackChunkdiscord_app.pop();

let ApplicationStreamingStore = Object.values(wpRequire.c).find(x => x?.exports?.Z?.__proto__?.getStreamerActiveStreamMetadata).exports.Z;
let RunningGameStore = Object.values(wpRequire.c).find(x => x?.exports?.ZP?.getRunningGames).exports.ZP;
let QuestsStore = Object.values(wpRequire.c).find(x => x?.exports?.Z?.__proto__?.getQuest).exports.Z;
let ChannelStore = Object.values(wpRequire.c).find(x => x?.exports?.Z?.__proto__?.getAllThreadsForParent).exports.Z;
let GuildChannelStore = Object.values(wpRequire.c).find(x => x?.exports?.ZP?.getSFWDefaultChannel).exports.ZP;
let FluxDispatcher = Object.values(wpRequire.c).find(x => x?.exports?.Z?.__proto__?.flushWaitQueue).exports.Z;
let api = Object.values(wpRequire.c).find(x => x?.exports?.tn?.get).exports.tn;

let quest = [...QuestsStore.quests.values()].find(x => x.id !== "1412491570820812933" && x.userStatus?.enrolledAt && !x.userStatus?.completedAt && new Date(x.config.expiresAt).getTime() > Date.now())
let isApp = typeof DiscordNative !== "undefined"

		}
fn()
console.log(`Spoofing video for ${questName}.`)
	} else if (taskName === "PLAY_ON_DESKTOP") {
  if (!isApp) {
    console.log("This no longer works in browser for non-video quests. Use the discord desktop app to complete the", questName, "quest!")
  } else {
    api.get({ url: `/applications/public?application_ids=${applicationId}` }).then(res => {
      const appData = res.body[0]
      const exeName = appData.executables.find(x => x.os === "win32").name.replace(">", "")

      const fakeGame = {
        cmdLine: `C:\\Program Files\\${appData.name}\\${exeName}`,
        exeName,
        exePath: `c:/program files/${appData.name.toLowerCase()}/${exeName}`,
        hidden: false,
        isLauncher: false,
        id: applicationId,
        name: appData.name,
        pid: pid,
        pidPath: [pid],
        processName: appData.name,
        start: Date.now(),
      }
      const realGames = RunningGameStore.getRunningGames()
      const fakeGames = [fakeGame]
      const realGetRunningGames = RunningGameStore.getRunningGames
      const realGetGameForPID = RunningGameStore.getGameForPID
      RunningGameStore.getRunningGames = () => fakeGames
      RunningGameStore.getGameForPID = (pid) => fakeGames.find(x => x.pid === pid)
      FluxDispatcher.dispatch({ type: "RUNNING_GAMES_CHANGE", removed: realGames, added: [fakeGame], games: fakeGames })

      let fn = data => {
        let progress = quest.config.configVersion === 1 ? data.userStatus.streamProgressSeconds : Math.floor(data.userStatus.progress.PLAY_ON_DESKTOP.value)
        console.log(`Quest progress: ${progress}/${secondsNeeded}`)

        if (progress >= secondsNeeded) {
          console.log("Quest completed!")

          RunningGameStore.getRunningGames = realGetRunningGames
          RunningGameStore.getGameForPID = realGetGameForPID
          FluxDispatcher.dispatch({ type: "RUNNING_GAMES_CHANGE", removed: [fakeGame], added: [], games: [] })
          FluxDispatcher.unsubscribe("QUESTS_SEND_HEARTBEAT_SUCCESS", fn)
        }
      }
      FluxDispatcher.subscribe("QUESTS_SEND_HEARTBEAT_SUCCESS", fn)

      console.log(`Spoofed your game to ${applicationName}. Wait for ${Math.ceil((secondsNeeded - secondsDone) / 60)} more minutes.`)
    })
  }
} else if (taskName === "STREAM_ON_DESKTOP") {
  if (!isApp) {
    console.log("This no longer works in browser for non-video quests. Use the discord desktop app to complete the", questName, "quest!")
  } else {
    let realFunc = ApplicationStreamingStore.getStreamerActiveStreamMetadata
    ApplicationStreamingStore.getStreamerActiveStreamMetadata = () => ({
      id: applicationId,
      pid,
      sourceName: null
    })

    let fn = data => {
      let progress = quest.config.configVersion === 1 ? data.userStatus.streamProgressSeconds : Math.floor(data.userStatus.progress.STREAM_ON_DESKTOP.value)
      console.log(`Quest progress: ${progress}/${secondsNeeded}`)

      if (progress >= secondsNeeded) {
        console.log("Quest completed!")

        ApplicationStreamingStore.getStreamerActiveStreamMetadata = realFunc
        FluxDispatcher.unsubscribe("QUESTS_SEND_HEARTBEAT_SUCCESS", fn)
      }
    }
    FluxDispatcher.subscribe("QUESTS_SEND_HEARTBEAT_SUCCESS", fn)

    console.log(`Spoofed your stream to ${applicationName}. Stream any window in vc for ${Math.ceil((secondsNeeded - secondsDone) / 60)} more minutes.`)
    console.log("Remember that you need at least 1 other person to be in the vc!")
  }
} else if (taskName === "PLAY_ACTIVITY") {
  const channelId = ChannelStore.getSortedPrivateChannels()[0]?.id ?? Object.values(GuildChannelStore.getAllGuilds()).find(x => x != null && x.VOCAL.length > 0).VOCAL[0].channel.id
  const streamKey = `call:${channelId}:1`

  let fn = async () => {
    console.log("Completing quest", questName, "-", quest.config.messages.questName)

    while (true) {
      const res = await api.post({ url: `/quests/${quest.id}/heartbeat`, body: { stream_key: streamKey, terminal: false } })
      const progress = res.body.progress.PLAY_ACTIVITY.value
      console.log(`Quest progress: ${progress}/${secondsNeeded}`)

      await new Promise(resolve => setTimeout(resolve, 20 * 1000))

      if (progress >= secondsNeeded) {
        await api.post({ url: `/quests/${quest.id}/heartbeat`, body: { stream_key: streamKey, terminal: true } })
        break
      }
    }

    console.log("Quest completed!")
  }
  fn()
}
}
