#!/bin/bash

curl -s http://localhost:8617/messages | jq -r '.'
