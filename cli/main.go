// Package cli is responsible for the execution of the CLI.
package cli

import (
	"context"
	"log"
	"os"

	"github.com/umatare5/twelvedata-exporter/config"
	"github.com/umatare5/twelvedata-exporter/internal"
	cli "github.com/urfave/cli/v3"
)

// Start is the entrypoint of this CLI
func Start() {
	cmd := &cli.Command{
		Name:      "twelvedata-exporter",
		Usage:     "Fetch quotes from Twelvedata API",
		UsageText: "twelvedata-exporter COMMAND [options...]",
		Version:   getVersion(),
		Flags:     registerFlags(),
		Action: func(ctx context.Context, cli *cli.Command) error {
			config := config.NewConfig(cli)
			server, _ := internal.NewServer(&config)

			server.Start()

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// registerFlags returns global flags
func registerFlags() []cli.Flag {
	flags := []cli.Flag{}
	flags = append(flags, registerWebListenAddressFlag()...)
	flags = append(flags, registerWebListenPortFlag()...)
	flags = append(flags, registerWebScrapePathFlag()...)
	flags = append(flags, registerAPIKeyFlag()...)
	return flags
}

// registerWebListenAddressFlag
func registerWebListenAddressFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    config.WebListenAddressFlagName,
			Usage:   "Set IP address",
			Aliases: []string{"I"},
			Value:   "0.0.0.0",
		},
	}
}

// registerWebListenPortFlag
func registerWebListenPortFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    config.WebListenPortFlagName,
			Usage:   "Set port number",
			Aliases: []string{"P"},
			Value:   10016,
		},
	}
}

// registerWebScrapePathFlag
func registerWebScrapePathFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    config.WebScrapePathFlagName,
			Usage:   "Set the path to expose metrics",
			Aliases: []string{"p"},
			Value:   "/price",
		},
	}
}

// registerAPIKeyFlag
func registerAPIKeyFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     config.TwelvedataAPIKeyFlagName,
			Usage:    "Set key to use twelvedata API",
			Aliases:  []string{"a"},
			Sources:  cli.EnvVars("TWELVEDATA_API_KEY"),
			Required: true,
		},
	}
}
