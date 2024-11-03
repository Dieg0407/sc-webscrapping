package main

import (
	"dieg0407/seace/internal/scrapper"
	"time"
)

func main() {
	time := time.Now()
	yesterday := time.AddDate(0, 0, -5)
	scrapper.Start(yesterday)
}
