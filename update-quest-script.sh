#!/usr/bin/env bash

set -e

gist_url="https://gist.githubusercontent.com/aamiaa/204cd9d42013ded9faf646fae7f89fbb/raw/CompleteDiscordQuest.md"
output_file="scripts/CompleteDiscordQuest.js"

gist_content=$(curl -sL "$gist_url")

js_code=$(echo "$gist_content" | awk '/```js$/,/```$/' | sed '1d;$d')

if [ -z "$js_code" ]; then
  echo "err: could not extract JS code from gist"
  exit 1
fi

if [ -f "$output_file" ]; then
  existing_js_code=$(awk '/^ \*\/$/{flag=1; next} flag' "$output_file" | sed '1{/^$/d}')

  if [ "$js_code" = "$existing_js_code" ]; then
    echo "no changes detected in script content"
    exit 0
  fi
fi

sync_date=$(date -u +"%Y-%m-%d %H:%M:%S UTC")

header="/*
 * Complete Recent Discord Quest
 * Source: https://gist.github.com/aamiaa/204cd9d42013ded9faf646fae7f89fbb
 * Author: aamiaa
 * License: GPL-3.0
 *
 * This script automates the completion of Discord quests by spoofing game activity,
 * video watching, or streaming depending on the quest requirements.
 *
 * Last synced: $sync_date
 */
"

echo "$header" > "$output_file"
echo "$js_code" >> "$output_file"

echo "script updated successfully"
