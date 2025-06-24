package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"

	"github.com/Ege-Guler/fsner/internal/scanner"
	"github.com/urfave/cli/v2"
)

func main() {

	const MiB = 1 << 20 // 1 MiB in bytes

	var wg sync.WaitGroup
	// buffered channel to hold search results
	ch := make(chan scanner.SearchResult, 100)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var pattern string
	var verbose bool
	var root string
	var maxResults int64 = -1

	var counter int64 = 0

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
			},
			&cli.Int64Flag{
				Name:        "max",
				Aliases:     []string{"m"},
				Usage:       "Maximum number of results to return, defaults to unlimited",
				Value:       -1,
				Required:    false,
				Destination: &maxResults,
			},
		},
		Action: func(c *cli.Context) error {
			if pattern == "" {
				return fmt.Errorf("pattern is required")
			}
			// validate regex once
			re := regexp.MustCompile(pattern)
			wg.Add(1)
			go scanner.SearchFile(ctx, root, re, verbose, ch, &wg)
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
		signal.Stop(signalChan)
	}()

	for {
		select {
		case sig, ok := <-signalChan:
			if !ok {
				return // channel closed, exit
			}
			fmt.Println("\nReceived shutdown signal, exiting gracefully:", sig)
			return
		case i, ok := <-ch:
			if !ok {
				return // channel closed, exit
			}
			fmt.Printf("%s: %2.f MiB\n", i.Path, float64(i.Info.Size())/float64(MiB))
			if maxResults > 0 {
				counter++
				if counter >= maxResults {
					fmt.Printf("\nReached maximum results limit of %d, exiting...\n", maxResults)
					cancel() // cancel the context to stop the search
					return
				}
			}
		}
	}
}
