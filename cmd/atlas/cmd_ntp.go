package main

import "github.com/urfave/cli"

// init injects our "ntp" related commands/options.
func init() {
	// Fill-in the various commands
	cliCommands = append(cliCommands, cli.Command{
		Name:        "ntp",
		Usage:       "get time from ntp server",
		Description: "send NTP queries to an host/IP",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "6, ipv6",
				Usage:       "displays only IPv6",
				Destination: &fWant6,
			},
			cli.BoolFlag{
				Name:        "4, ipv4",
				Usage:       "displays only IPv4",
				Destination: &fWant4,
			},
		},
		Action: cmdNTP,
	})
}

func cmdNTP(c *cli.Context) error {
	return nil
}
