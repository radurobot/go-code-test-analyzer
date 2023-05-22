package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type CmdArgs struct {
	Markdown bool
}

func parseArgs() CmdArgs {
	markdown := flag.Bool("markdown", false, "Output in markdown format")
	flag.Parse()

	return CmdArgs{
		Markdown: *markdown,
	}
}

func main() {
	args := parseArgs()

	fileCounts := make(map[string]*FileAnalysis)
	files := make(chan string, 100)
	wg := sync.WaitGroup{}
	results := make(chan *FileResult, 100)
	workers := runtime.NumCPU()

	go func() {
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				files <- path
			}
			return nil
		})

		if err != nil {
			fmt.Println("Error walking the path", err)
			return
		}

		close(files)
	}()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range files {
				analysis, err := analyzeFile(path)
				if err != nil {
					fmt.Println("Error analyzing file", err)
					continue
				}
				results <- &FileResult{Path: path, Analysis: analysis}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fileCounts[result.Path] = result.Analysis
	}
	if args.Markdown {
		fmt.Println(convertToMarkdown(fileCounts))
	} else {
		printCounts(fileCounts)
	}
}
