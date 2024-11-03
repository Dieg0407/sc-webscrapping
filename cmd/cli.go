package main

import (
	"dieg0407/seace/internal/scrapper"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	logger := log.New(os.Stderr, "[cli] ", log.LstdFlags)

	var dateString string
	var date time.Time
	var layout = "2006-01-02"
	var err error

	app := &cli.App{
		Name:  "scrapper",
		Usage: "Use this to extract information from the page",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "process-date",
				Aliases:     []string{"d"},
				Usage:       "The date to process",
				Destination: &dateString,
			},
		},
		Action: func(*cli.Context) error {
			if dateString == "" {
				return fmt.Errorf("You must provide a date")
			}
			if date, err = time.Parse(layout, dateString); err != nil {
				return fmt.Errorf("Invalid date format, you should use YYYY-MM-DD")
			}
			scrapper.Start(date)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
