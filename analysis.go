package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type FileAnalysis struct {
	Files         int
	Functions     int
	TestFunctions int
}
type FileResult struct {
	Path     string
	Analysis *FileAnalysis
}

func analyzeFile(filename string) (*FileAnalysis, error) {
	analysis := &FileAnalysis{}

	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return nil, errors.New("Error parsing file: " + err.Error())
	}

	analysis.Functions = countFunctions(parsedFile)

	testFilename := strings.TrimSuffix(filename, ".go") + "_test.go"
	if _, err := os.Stat(testFilename); err == nil {
		testAnalysis, err := analyzeFile(testFilename)
		if err != nil {
			return nil, err
		}
		analysis.TestFunctions = testAnalysis.Functions
		analysis.Files = testAnalysis.Files + 1
	} else {
		analysis.Files = 1
	}

	return analysis, nil
}

func countFunctions(file *ast.File) int {
	functionCount := 0
	for _, decl := range file.Decls {
		if _, ok := decl.(*ast.FuncDecl); ok {
			functionCount++
		}
	}
	return functionCount
}

func aggregateCounts(fileCounts map[string]*FileAnalysis) *FileAnalysis {
	total := &FileAnalysis{
		Files:         len(fileCounts),
		Functions:     0,
		TestFunctions: 0,
	}
	for _, analysis := range fileCounts {
		total.Functions += analysis.Functions
		total.TestFunctions += analysis.TestFunctions
	}
	return total
}
