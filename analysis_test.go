package main

import (
	"go/ast"
	"os"
	"testing"
)

func TestAnalyzeFile(t *testing.T) {
	// Test a valid Go file with no associated test file
	filetestfile, err := os.CreateTemp("", "testfile*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filetestfile.Name())

	code := `
    package main

    func main() {}
    `
	if _, err := filetestfile.Write([]byte(code)); err != nil {
		t.Fatal(err)
	}

	analysis, err := analyzeFile(filetestfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if analysis.Functions != 1 {
		t.Errorf("Expected 1 function, got %d", analysis.Functions)
	}
	if analysis.TestFunctions != 0 {
		t.Errorf("Expected 0 test functions, got %d", analysis.TestFunctions)
	}
	if analysis.Files != 1 {
		t.Errorf("Expected 1 file, got %d", analysis.Files)
	}

	// Test a valid Go file with an associated test file
	testFile, err := os.Create(filetestfile.Name()[:len(filetestfile.Name())-3] + "_test.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile.Name())

	testCode := `
    package main

    func TestMain(t *testing.T) {}
    `
	if _, err := testFile.Write([]byte(testCode)); err != nil {
		t.Fatal(err)
	}

	analysis, err = analyzeFile(filetestfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if analysis.Functions != 1 {
		t.Errorf("Expected 1 function, got %d", analysis.Functions)
	}
	if analysis.TestFunctions != 1 {
		t.Errorf("Expected 1 test function, got %d", analysis.TestFunctions)
	}
	if analysis.Files != 2 {
		t.Errorf("Expected 2 files, got %d", analysis.Files)
	}

	// Test an invalid Go file
	invalidFile, err := os.CreateTemp("", "invalidfile*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(invalidFile.Name())

	invalidCode := `
    package main
	package main
    func main() {
        fmt.Println("Hello, world!")
    }
    `
	if _, err := invalidFile.Write([]byte(invalidCode)); err != nil {
		t.Fatal(err)
	}

	_, err = analyzeFile(invalidFile.Name())
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Test a non-existent file
	_, err = analyzeFile("nonexistent.go")
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestCountFunctions(t *testing.T) {
	// Test a file with no functions
	file := &ast.File{
		Name:  &ast.Ident{Name: "test"},
		Decls: []ast.Decl{},
	}
	functionCount := countFunctions(file)
	if functionCount != 0 {
		t.Errorf("Expected 0 functions, got %d", functionCount)
	}

	// Test a file with one function
	file = &ast.File{
		Name: &ast.Ident{Name: "test"},
		Decls: []ast.Decl{
			&ast.FuncDecl{
				Name: &ast.Ident{Name: "testFunc"},
			},
		},
	}
	functionCount = countFunctions(file)
	if functionCount != 1 {
		t.Errorf("Expected 1 function, got %d", functionCount)
	}

	// Test a file with multiple functions
	file = &ast.File{
		Name: &ast.Ident{Name: "test"},
		Decls: []ast.Decl{
			&ast.FuncDecl{
				Name: &ast.Ident{Name: "testFunc1"},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{Name: "testFunc2"},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{Name: "testFunc3"},
			},
		},
	}
	functionCount = countFunctions(file)
	if functionCount != 3 {
		t.Errorf("Expected 3 functions, got %d", functionCount)
	}
}

func TestAggregateCounts(t *testing.T) {
	// Test with no file counts
	fileCounts := map[string]*FileAnalysis{}
	total := aggregateCounts(fileCounts)
	if total.Files != 0 {
		t.Errorf("Expected 0 files, got %d", total.Files)
	}
	if total.Functions != 0 {
		t.Errorf("Expected 0 functions, got %d", total.Functions)
	}
	if total.TestFunctions != 0 {
		t.Errorf("Expected 0 test functions, got %d", total.TestFunctions)
	}

	// Test with one file count
	fileCounts = map[string]*FileAnalysis{
		"file1.go": {
			Files:         1,
			Functions:     5,
			TestFunctions: 2,
		},
	}
	total = aggregateCounts(fileCounts)
	if total.Files != 1 {
		t.Errorf("Expected 1 file, got %d", total.Files)
	}
	if total.Functions != 5 {
		t.Errorf("Expected 5 functions, got %d", total.Functions)
	}
	if total.TestFunctions != 2 {
		t.Errorf("Expected 2 test functions, got %d", total.TestFunctions)
	}

	// Test with multiple file counts
	fileCounts = map[string]*FileAnalysis{
		"file1.go": {
			Files:         1,
			Functions:     5,
			TestFunctions: 2,
		},
		"file2.go": {
			Files:         1,
			Functions:     3,
			TestFunctions: 1,
		},
		"file3.go": {
			Files:         1,
			Functions:     2,
			TestFunctions: 0,
		},
	}
	total = aggregateCounts(fileCounts)
	if total.Files != 3 {
		t.Errorf("Expected 3 files, got %d", total.Files)
	}
	if total.Functions != 10 {
		t.Errorf("Expected 10 functions, got %d", total.Functions)
	}
	if total.TestFunctions != 3 {
		t.Errorf("Expected 3 test functions, got %d", total.TestFunctions)
	}
}
