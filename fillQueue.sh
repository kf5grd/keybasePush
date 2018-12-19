#!/bin/bash

for i in {1..10}; do
    curl -H "Content-Type: application/json" \
	 -d "{\"title\":\"message$i\",\"content\":\"asdf\",\"targets\":[\"server\"]}" \
	 http://localhost:8617/messages
done
