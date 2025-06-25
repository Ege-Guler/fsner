package main

import (
	"fmt"
	"regexp"

	"github.com/urfave/cli/v2"
)

func setupApp(cfg *Config) *cli.App {

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "Print the version of fsner",
	}

	return &cli.App{
		Name:    "fsner",
		Usage:   "A file system search tool",
		Version: "0.1.0",
		Flags:   cliFlags(cfg),
		Action: func(c *cli.Context) error {
			if cfg.Pattern == "" {
				return fmt.Errorf("pattern is required")
			}
			// validate regex once
			re, err := regexp.Compile(cfg.Pattern)
			if err != nil {
				return fmt.Errorf("invalid regex pattern, %w", err)
			}
			cfg.Regex = re
			return nil
		},
	}
}

func cliFlags(cfg *Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "pattern",
			Aliases:     []string{"p"},
			Usage:       "Pattern to search for in file names",
			Required:    true,
			Destination: &cfg.Pattern,
		},
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"v"},
			Usage:       "Enable verbose output, defaults to false",
			Required:    false,
			Value:       false,
			Destination: &cfg.Verbose,
		},
		&cli.StringFlag{
			Name:        "root",
			Aliases:     []string{"r"},
			Usage:       "Root directory to start the search from, defaults to /",
			Value:       "/",
			Required:    true,
			Destination: &cfg.Root,
		},
		&cli.Int64Flag{
			Name:        "max",
			Aliases:     []string{"m"},
			Usage:       "Maximum number of results to return, defaults to unlimited",
			Value:       -1,
			Required:    false,
			Destination: &cfg.MaxResults,
		},
		&cli.BoolFlag{
			Name:        "file-size",
			Aliases:     []string{"s"},
			Usage:       "Display file size in MiB",
			Value:       false,
			Destination: &cfg.FileSize,
		},
	}
}
