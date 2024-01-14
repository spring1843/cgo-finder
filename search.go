package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// cGOResult stores CGO usage information
type cGOResult struct {
	fileName string
	line     int
}

const (
	// goCLI is the name of the go CLI command
	goCLI = "go"
	// c is the name of the C package that indicates use of CGo
	c = "C"
	// bufferSize is the number of bytes to read from the top of a Go file to search for CGo
	bufferSize = 2048 // 2KB
)

var (
	// goCacheArgs are the arguments to the go CLI to find the Go module cache dir
	goCacheArgs = []string{"env", "GOMODCACHE"}
	// goModDir is the directory where Go caches the modules
	goModDir = findGoModDir()
)

// goCLIPath returns the path to the go CLI.
func goCLIPath() string {
	goPath, err := exec.LookPath(goCLI)
	if err != nil {
		panic(fmt.Errorf("failed to find go CLI. %w", err))
	}
	return goPath
}

// findGoModDir returns the directory where Go caches the modules.
func findGoModDir() string {
	cmd := exec.Command(goCLIPath(), goCacheArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Errorf("failed to execute the go command. %w", err))
	}
	return strings.TrimSuffix(string(output), "\n")
}

// dirExists returns true if directory exists and false otherwise
func dirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(fmt.Errorf("failed to check directory existence. %w", err))
	}
	return true
}

// recursiveFindCGo recursively searches for CGO usage in *.go files in a directory
// and it's subdirectories and returns file names, line numbers and excerpts
func recursiveFindCGo(dir string) []cGOResult {
	output := []cGOResult{}

	searchFileForCGo := func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			imports := parseImports(readFileHead(filePath))
			if lineNumber, ok := imports[c]; ok {
				output = append(output, cGOResult{fileName: filePath, line: lineNumber})
			}
		}
		return nil
	}

	err := filepath.Walk(dir, searchFileForCGo)
	if err != nil {
		panic(fmt.Errorf("failed to walk directory. %w", err))
	}

	return output
}

// readFileHead reads the first bufferSize bytes of a file and returns them as a string
func readFileHead(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to open file. %w", err))
	}
	defer file.Close()

	buffer := make([]byte, bufferSize)

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		panic(fmt.Errorf("failed read file. %w", err))
	}

	return string(buffer[:n])
}
