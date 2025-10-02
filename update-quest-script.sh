#!/usr/bin/env bash

set -e

GIST_URL="https://gist.githubusercontent.com/aamiaa/204cd9d42013ded9faf646fae7f89fbb/raw/CompleteDiscordQuest.md"
OUTPUT_FILE="scripts/CompleteDiscordQuest.js"

GIST_CONTENT=$(curl -sL "$GIST_URL")

JS_CODE=$(echo "$GIST_CONTENT" | awk '/```js$/,/```$/' | sed '1d;$d')

if [ -z "$JS_CODE" ]; then
  echo "err: could not extract JS code from gist"
  exit 1
fi

SYNC_DATE=$(date -u +"%Y-%m-%d %H:%M:%S UTC")

HEADER="/*
 * Complete Recent Discord Quest
 * Source: https://gist.github.com/aamiaa/204cd9d42013ded9faf646fae7f89fbb
 * Author: aamiaa
 * License: GPL-3.0
 *
 * This script automates the completion of Discord quests by spoofing game activity,
 * video watching, or streaming depending on the quest requirements.
 *
 * Last synced: $SYNC_DATE
 */
"

echo "$HEADER" > "$OUTPUT_FILE"
echo "$JS_CODE" >> "$OUTPUT_FILE"

echo "script updated successfully"
