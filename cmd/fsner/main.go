package main

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"sync"

	"github.com/Ege-Guler/fsner/internal/scanner"
	"github.com/urfave/cli/v2"
)

func main() {

	var wg sync.WaitGroup
	ch := make(chan fs.FileInfo, 100)

	var pattern string
	var verbose bool
	var root string

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "Print the version of fsner",
	}

	app := &cli.App{
		Name:    "fsner",
		Usage:   "A file system search tool",
		Version: "0.1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "pattern",
				Aliases:     []string{"p"},
				Usage:       "Pattern to search for in file names",
				Required:    true,
				Destination: &pattern,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "Enable verbose output, defaults to false",
				Required:    false,
				Value:       false,
				Destination: &verbose,
			},
			&cli.StringFlag{
				Name:        "root",
				Aliases:     []string{"r"},
				Usage:       "Root directory to start the search from, defaults to /",
				Value:       "/",
				Required:    true,
				Destination: &root,
			}},
		Action: func(c *cli.Context) error {
			if pattern == "" {
				return fmt.Errorf("pattern is required")
			}
			// validate regex once
			re := regexp.MustCompile(pattern)
			wg.Add(1)
			go scanner.SearchFile(root, re, verbose, ch, &wg)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "\nError running app: %v\n", err)
		return
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := range ch {
		fmt.Println(i.Name())
	}

}
