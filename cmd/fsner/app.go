package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
)

func setupApp(cfg *Config) *cli.App {

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "Print the version of fsner",
	}

	return &cli.App{
		Name:                 "fsner",
		Usage:                "A file system search tool",
		EnableBashCompletion: true,
		Version:              "0.1.0",
		Flags:                cliFlags(cfg),
		BashComplete: func(c *cli.Context) {

			args := os.Args

			if len(args) > 1 {
				last := args[len(args)-1]
				prev := args[len(args)-2]

				if last == "--root" || strings.HasPrefix(prev, "--root") {
					suggestDirectories()
					return
				}
				if last == "--max" || strings.HasPrefix(prev, "--max") {
					suggestMax()
					return
				}
			}

			if !c.IsSet("pattern") {
				fmt.Println("--pattern")
			}
			if !c.IsSet("verbose") {
				fmt.Println("--verbose")
			}
			if !c.IsSet("root") {
				fmt.Println("--root")
			}
			if !c.IsSet("max") {
				fmt.Println("--max")
			}

			if !c.IsSet("file-size") {
				fmt.Println("--file-size")
			}

		},
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

func suggestDirectories() {

	fmt.Println("/")
	home, err := os.UserHomeDir()
	if err == nil {
		fmt.Println(home)
		fmt.Println(filepath.Join(home, "Downloads"))
		fmt.Println(filepath.Join(home, "Documents"))
		fmt.Println(filepath.Join(home, "Pictures"))
	}
}

func suggestMax() {
	fmt.Println("5")
	fmt.Println("10")
	fmt.Println("20")
	fmt.Println("50")
	fmt.Println("100")

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
		&cli.BoolFlag{
			Name:   "generate-bash-completion",
			Hidden: true,
		},
	}
}
