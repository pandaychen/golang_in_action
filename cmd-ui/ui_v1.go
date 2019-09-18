package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"sort"
)

func main() {
	var language string

	app := cli.NewApp()
	app.Name = "GoTest"
	app.Usage = "hello world"
	app.Version = "1.2.3"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 8000,
			Usage: "listening port",
		},
		cli.StringFlag{
			Name:        "lang, l",
			Value:       "english",
			Usage:       "read from `FILE`",
			Destination: &language,
		},
	}
	app.Commands = []cli.Command{
		{
			//use ./ui_v1 add 
			Name:     "add",
			Aliases:  []string{"a"},
			Usage:    "calc 1+1",
			Category: "arithmetic",
			Action: func(c *cli.Context) error {
				fmt.Println("1 + 1 = ", 1+1)
				return nil
			},
		},
		{
			//use ./ui_v1 sub
			Name:     "sub",
			Aliases:  []string{"s"},
			Usage:    "calc 5-3",
			Category: "arithmetic",
			Action: func(c *cli.Context) error {
				fmt.Println("5 - 3 = ", 5-3)
				return nil
			},
		},
		{
			Name:     "db",
			Usage:    "database operations",
			Category: "database",
			Subcommands: []cli.Command{
				{
					Name:  "insert",
					Usage: "insert data",
					Action: func(c *cli.Context) error {
						fmt.Println("insert subcommand")
						return nil
					},
				},
				{
					Name:  "delete",
					Usage: "delete data",
					Action: func(c *cli.Context) error {
						fmt.Println("delete subcommand")
						return nil
					},
				},
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println("BOOM!")
		fmt.Println(c.String("lang"), c.Int("port"))
		fmt.Println(language)

		// if c.Int("port") == 8000 {
		// 	return cli.NewExitError("invalid port", 88)
		// }

		return nil
	}
	app.Before = func(c *cli.Context) error {
		fmt.Println("app Before")
		return nil
	}
	app.After = func(c *cli.Context) error {
		fmt.Println("app After")
		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	cli.HelpFlag = cli.BoolFlag{
		Name:  "help, h",
		Usage: "Help!Help!",
	}

	cli.VersionFlag = cli.BoolFlag{
		Name:  "print-version, v",
		Usage: "print version",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
