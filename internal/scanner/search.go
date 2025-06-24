package scanner

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type SearchResult struct {
	Path string
	Info fs.FileInfo
}

// if verbose is false, skip dir read error messages
func SearchFile(ctx context.Context, root string, re *regexp.Regexp, verbose bool, ch chan<- SearchResult, wg *sync.WaitGroup) {

	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	entries, err := os.ReadDir(root)

	if err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "Failed to read dir %s: %v\n", root, err)
		}
		return
	}

	for _, entry := range entries {

		select {
		case <-ctx.Done():
			return
		default:
		}

		full_path := filepath.Join(root, entry.Name())

		if entry.IsDir() {
			wg.Add(1)
			go SearchFile(ctx, full_path, re, verbose, ch, wg)
		} else {
			if matched := re.MatchString(entry.Name()); matched {
				info, err := entry.Info()
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to get info on %s: %v\n", entry.Name(), err)
					return
				}
				ch <- SearchResult{full_path, info}
			}
		}
	}

}
