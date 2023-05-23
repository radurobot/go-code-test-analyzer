package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/urfave/cli/v2"
)

const Version = "v0.1.1"

type CmdArgs struct {
	Markdown       bool
	Html           bool
	LatestVersion  bool
	CurrentVersion bool
	Path           []string
	Exclude        []string
}

func parseArgs() CmdArgs {
	cmdArgs := CmdArgs{}

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "markdown",
				Aliases:     []string{"m"},
				Value:       false,
				Usage:       "Output in markdown format",
				Destination: &cmdArgs.Markdown,
			},
			&cli.BoolFlag{
				Name:        "html",
				Value:       false,
				Usage:       "Output in HTML format",
				Destination: &cmdArgs.Html,
			},
			&cli.BoolFlag{
				Name:        "latest-version",
				Aliases:     []string{"lv"},
				Value:       false,
				Usage:       "Prints the latest version and exit",
				Destination: &cmdArgs.LatestVersion,
			},
			&cli.BoolFlag{
				Name:        "version",
				Aliases:     []string{"v"},
				Value:       false,
				Usage:       "Prints the current version and exit",
				Destination: &cmdArgs.CurrentVersion,
			},
			&cli.StringSliceFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Paths to analyze",
			},
			&cli.StringSliceFlag{
				Name:    "exclude",
				Aliases: []string{"e"},
				Usage:   "Comma separated list of directories to exclude",
			},
		},
		Action: func(c *cli.Context) error {
			cmdArgs.Path = c.StringSlice("path")
			cmdArgs.Exclude = c.StringSlice("exclude")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	return cmdArgs
}

func main() {
	args := parseArgs()
	if args.LatestVersion {
		printLatestVersion()
		return
	}
	if args.CurrentVersion {
		fmt.Println(Version)
		return
	}
	if len(args.Path) == 0 {
		fmt.Println("Please provide a path to analyze")
		return
	}
	if err := checkForConflictsInFlags(args); err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}
	fileCounts := make(map[string]*FileAnalysis)
	files := make(chan string, 100)
	wg := sync.WaitGroup{}
	results := make(chan *FileResult, 100)
	workers := runtime.NumCPU()
	excludeMap := make(map[string]bool)
	for _, exclude := range args.Exclude {
		excludeMap[exclude] = true
	}
	go func() {
		for _, path := range args.Path {
			log.Println("Walking path", path)
			err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if excludeMap[path] {
					return filepath.SkipDir
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
		}
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
	switch {
	case args.Markdown:
		fmt.Println(convertToMarkdown(fileCounts))
	case args.Html:
		fmt.Println(convertToHtmlTable(fileCounts))
	default:
		printCounts(fileCounts)
	}
}

func printLatestVersion() {
	url := "https://api.github.com/repos/radurobot/go-code-test-analyzer/releases/latest"
	resp, err := http.Get(url)
	if resp.StatusCode != http.StatusOK || err != nil {
		fmt.Println("Error getting version from github")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body")
		return
	}
	var release struct {
		TagName string `json:"tag_name"`
	}
	err = json.Unmarshal(body, &release)
	if err != nil {
		fmt.Println("Error parsing json")
		return
	}
	fmt.Println(release.TagName)
}

func checkForConflictsInFlags(args CmdArgs) error {
	if args.Markdown && args.Html {
		return fmt.Errorf("cannot use both markdown and html flags")
	}

	// Create a map of excluded paths for faster lookup
	excludedPaths := make(map[string]bool)
	for _, exclude := range args.Exclude {
		// check if exclude is valid
		if _, err := os.Stat(exclude); err != nil {
			return fmt.Errorf("exclude path %s is not valid", exclude)
		}
		excludedPaths[filepath.Clean(exclude)] = true
	}

	// Check if any paths are also in excludedPaths
	for _, path := range args.Path {
		// check if path is valid
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("path %s is not valid", path)
		}
		if excludedPaths[filepath.Clean(path)] {
			return fmt.Errorf("path %s is also in exclude list", path)
		}
	}

	return nil
}
