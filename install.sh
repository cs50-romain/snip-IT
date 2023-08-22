#!/bin/bash

cp data.json temp.json

go build -o snip-it.go

mkdir -p /usr/local/bin

mv snip-it /usr/local/bin

echo "snip-it is now installed. You can run 'snip-it' command."
