// Package config is responsible for the execution of the CLI.
package config

import (
	"errors"

	"github.com/jinzhu/configor"
	cli "github.com/urfave/cli/v3"

	"github.com/umatare5/twelvedata-exporter/log"
)

// Config flag names.
const (
	WebListenAddressFlagName = "web.listen-address"
	WebListenPortFlagName    = "web.listen-port"
	WebScrapePathFlagName    = "web.scrape-path"
	TwelvedataAPIKeyFlagName = "twelvedata.api-key" //nolint:gosec // Flag name, not a credential.
)

// Config struct.
type Config struct {
	WebListenAddress string
	WebListenPort    int
	WebScrapePath    string
	TwelvedataAPIKey string
}

// NewConfig returns Config struct.
func NewConfig(cmd *cli.Command) Config {
	config := Config{
		WebListenAddress: cmd.String(WebListenAddressFlagName),
		WebListenPort:    cmd.Int(WebListenPortFlagName),
		WebScrapePath:    cmd.String(WebScrapePathFlagName),
		TwelvedataAPIKey: cmd.String(TwelvedataAPIKeyFlagName),
	}

	err := configor.New(&configor.Config{}).Load(&config)
	if err != nil {
		log.Fatal(err)
	}

	if err := isValidWebListenAddressFlag(config.WebListenAddress); err != nil {
		log.Fatal(err)
	}

	if err := isValidWebListenPortFlag(config.WebListenPort); err != nil {
		log.Fatal(err)
	}

	if err := isValidWebScrapePathFlag(config.WebScrapePath); err != nil {
		log.Fatal(err)
	}

	if err := isValidTwelvedataAPIKeyFlag(config.TwelvedataAPIKey); err != nil {
		log.Fatal(err)
	}

	return config
}

func isValidWebListenAddressFlag(_ string) error {
	return nil
}

func isValidWebListenPortFlag(_ int) error {
	return nil
}

func isValidWebScrapePathFlag(_ string) error {
	return nil
}

func isValidTwelvedataAPIKeyFlag(apikey string) error {
	if apikey == "" {
		return errors.New("environment variable 'TWELVEDATA_API_KEY' is not set")
	}

	return nil
}
