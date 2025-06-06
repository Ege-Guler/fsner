package main

import (
	"fmt"
	"io/fs"
	"regexp"
	"sync"

	"github.com/Ege-Guler/fsner/internal/scanner"
)

func main() {

	var wg sync.WaitGroup

	ch := make(chan fs.FileInfo, 100)

	pattern := `ders`

	// validate regex once
	re := regexp.MustCompile(pattern)

	wg.Add(1)
	go scanner.SearchFile("/", re, false, ch, &wg)

	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := range ch {
		fmt.Println(i.Name())
	}

}
