package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "domain-finder"
	app.Description = "Search domains in registo.br candidates to monthly release"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "refresh",
			Usage: "Refresh domain list of current month",
		},
		cli.BoolFlag{
			Name:  "diff",
			Usage: "Diff domain list between current month and last month",
		},
	}

	refreshCommand := cli.Command{
		Name: "refresh",
		Action: func(c *cli.Context) error {
			return refreshDomainsList()
		},
	}

	diffCommand := cli.Command{
		Name:        "diff",
		Description: "Diff available domains with last month list",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "removed",
			},
		},
		Action: func(c *cli.Context) error {
			now := time.Now()
			currentMonth := fmt.Sprintf("history/release-%d-%d.txt", now.Month(), now.Year())
			lastMonth := fmt.Sprintf("history/release-%d-%d.txt", now.Month()-1, now.Year())

			newDomains, removedDomains, err := diffFiles(currentMonth, lastMonth)
			if err != nil {
				return err
			}

			if c.Bool("removed") {
				log.Println(removedDomains)
			} else {
				log.Println(newDomains)
			}

			return nil
		},
	}

	app.Commands = []cli.Command{
		refreshCommand,
		diffCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
