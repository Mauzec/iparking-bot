#!/usr/bin/env bash
DIR="$(cd "$(dirname "$0")"/../data && pwd)"
TMP="$DIR/data.json.tmp"
DEST="$DIR/data.json"

while true; do
    D=$(( RANDOM % 100 ))
    echo "{\"distance\": $D}" > "$TMP" && mv "$TMP" "$DEST"
    sleep 2
done