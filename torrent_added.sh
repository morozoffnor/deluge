#!/usr/bin/env bash

torrentid=$1
torrentname=$2
torrentpath=$3

URL="localhost:80"

curl -X POST "$URL" \
  -H "Content-Type: application/json" \
  -d @- <<EOF
  {
    "event_type": "added",
    "id": "$torrentid",
    "name": "$torrentname",
    "path": "$torrentpath"
  }
EOF
