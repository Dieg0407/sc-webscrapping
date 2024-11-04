#!/bin/sh

go build -o sc-scrubber cmd/cli.go 
sudo mv sc-scrubber /usr/local/bin/sc-scrubber
