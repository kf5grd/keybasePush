#!/bin/bash

msgid="$1"
curl -s http://localhost:8617/messages/$msgid | jq -r '.'
