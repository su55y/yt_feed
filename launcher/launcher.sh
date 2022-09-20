#!/bin/sh

SCRIPTPATH="$(cd -- "$(dirname "$0")" >/dev/null 2>&1 || exit 1 ; pwd -P)"
[ ! -f "$SCRIPTPATH/yt_feed" ] && \
    notify-send "yt_feed" "blocks wrapper not found"

rofi  -modi blocks \
    -show blocks \
    -theme yt_search.rasi \
    -normal-window \
    -blocks-wrap "$SCRIPTPATH/yt_feed"
