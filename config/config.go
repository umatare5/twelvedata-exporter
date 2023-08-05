// Package config is responsible for the execution of the CLI.
package config

import (
	"errors"
	"log"

	"github.com/jinzhu/configor"
	"github.com/urfave/cli/v2"
)

// Config flag names
const (
	WebListenAddressFlagName = "web.listen-address"
	WebListenPortFlagName    = "web.listen-port"
	WebScrapePathFlagName    = "web.scrape-path"
	TwelvedataAPIKeyFlagName = "twelvedata.api-key"
)

// Config struct
type Config struct {
	WebListenAddress string
	WebListenPort    int
	WebScrapePath    string
	TwelvedataAPIKey string
}

// NewConfig returns Config struct
func NewConfig(ctx *cli.Context) Config {
	config := Config{
		WebListenAddress: ctx.String(WebListenAddressFlagName),
		WebListenPort:    ctx.Int(WebListenPortFlagName),
		WebScrapePath:    ctx.String(WebScrapePathFlagName),
		TwelvedataAPIKey: ctx.String(TwelvedataAPIKeyFlagName),
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
		return errors.New("Environment variable 'TWELVEDATA_API_KEY' is not set")
	}

	return nil
}
