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
)

type Config struct {
	Pattern    string
	Verbose    bool
	Root       string
	MaxResults int64
	FileSize   bool
	Regex      *regexp.Regexp
}

func run() int {
	cfg := &Config{}
	app := setupApp(cfg)
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	return runSearch(cfg)
}

func runSearch(cfg *Config) int {

	const MiB = 1 << 20 // 1 MiB in bytes
	const KiB = 1 << 10 // 1 KiB in bytes

	var wg sync.WaitGroup
	// buffered channel to hold search results
	ch := make(chan scanner.SearchResult, 100)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var counter int64 = 0

	wg.Add(1)
	go scanner.SearchFile(ctx, cfg.Root, cfg.Regex, cfg.Verbose, ch, &wg)

	go func() {
		wg.Wait()
		close(ch)
		signal.Stop(signalChan)
	}()

	for {

		select {
		case sig, ok := <-signalChan:
			if !ok {
				return 0
			}
			fmt.Println("\nReceived shutdown signal, exiting gracefully:", sig)
			return -1

		case result, ok := <-ch:
			if !ok {
				return 0
			}
			fmt.Printf("%s", result.Path)

			if cfg.FileSize {
				mib := float64(result.Info.Size()) / float64(MiB)
				if mib < 1 {
					fmt.Printf(" : %.2f KiB", float64(result.Info.Size())/float64(KiB))
				} else {
					fmt.Printf(" : %.2f MiB", mib)
				}
			}
			fmt.Println()

			if cfg.MaxResults > 0 {
				counter++
				if counter >= cfg.MaxResults {
					fmt.Printf("\nReached maximum results limit of %d, exiting...\n", cfg.MaxResults)
					cancel() // cancel the context to stop the search
					return 0
				}

			}

		}

	}

}
