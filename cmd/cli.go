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
		Usage: "Utiliza esto para extraer información de la página",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "fecha-proceso",
				Aliases:     []string{"d"},
				Usage:       "La fecha a procesar",
				Destination: &dateString,
			},
		},
		Action: func(*cli.Context) error {
			if dateString == "" {
				return fmt.Errorf("Debes proporcionar una fecha")
			}
			if date, err = time.Parse(layout, dateString); err != nil {
				return fmt.Errorf("Formato de fecha inválido, debes usar YYYY-MM-DD")
			}
			scrapper.Start(date)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
