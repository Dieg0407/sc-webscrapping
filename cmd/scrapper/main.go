package main

import (
	"dieg0407/seace/internal/scrapper"
	"time"
)

func main() {
	time := time.Now()
	yesterday := time.AddDate(0, 0, -1)
	scrapper.Start(yesterday)
}
