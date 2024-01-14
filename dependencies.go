package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type dependency struct {
	name           string
	version        string
	isIncompatible bool
	isIndirect     bool
	cGoResults     []cGOResult
}

const (
	incompatible = "+incompatible"
	indirect     = "// indirect"
	vendorDir    = "vendor"
	require      = "require"
	module       = "module"
	goVersion    = "go"
	closeBracket = ")"
)

// checkForCGo inspects the dependency for CGO usage and prints information if the
// dependency does use CGO. It returns true if the dependency uses CGO and false otherwise.
func (d *dependency) checkForCGo() bool {
	d.cGoResults = recursiveFindCGo(d.sourceDirectory())
	if len(d.cGoResults) > 0 {
		fmt.Printf("Dependency %s uses CGO\n", d.name)
		for _, result := range d.cGoResults {
			fmt.Printf("\t%s:%d\n", result.fileName, result.line)
		}
		return true
	}
	return false
}

// sourceDirectory returns the source directory of the dependency
func (d *dependency) sourceDirectory() string {
	vendor := fmt.Sprintf("%s/%s", vendorDir, d.name)
	if dirExists(vendor) {
		return vendor
	}
	goModCache := fmt.Sprintf("%s/%s@%s", goModDir, exclamationFormat(d.name), d.version)
	if !dirExists(goModCache) {
		panic(fmt.Errorf("go module cache for %s does not exit in %s\n. Please run `go mod tidy` or `go mod vendor`", d.name, goModCache))
	}
	return goModCache
}

// exclamationFormat changes capital letters to lower case and adds an exclamation mark
// example: FooBar -> !foo!bar
func exclamationFormat(dependencyName string) string {
	formattedName := ""
	for _, char := range dependencyName {
		if char >= 'A' && char <= 'Z' {
			formattedName += "!" + string(char+32)
		} else {
			formattedName += string(char)
		}
	}
	return formattedName
}

// readGoMod reads the contents of the go.mod file
func readGoMod() string {
	filePath := "go.mod"
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return string(content)
}

// lineIsNotDependency returns true if a line in go.mod should not be interpreted as a dependency
func lineIsNotDependency(line string) bool {
	return strings.HasPrefix(line, module) ||
		strings.HasPrefix(line, require) ||
		strings.HasPrefix(line, goVersion) ||
		strings.HasPrefix(line, closeBracket)
}

// parseGoModDependencies parses the contents of the go.mod file and returns a slice of dependencies
func parseGoModDependencies(goModContent string) []dependency {
	var dependencies []dependency
	lines := strings.Split(goModContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if lineIsNotDependency(line) {
			continue
		}
		dependencyLine := strings.TrimSpace(line)
		if dependencyLine != "" {
			dependencies = append(dependencies, toDependency(dependencyLine))
		}
	}
	return dependencies
}

// parseImports parses the imports of a go file and returns a map of import names to line numbers
func parseImports(input string) map[string]int {
	imports := make(map[string]int)
	lines := strings.Split(input, "\n")

	re := regexp.MustCompile(`^import\s+"(.*)"$|^"(.*)"$`)

	for i, line := range lines {
		line = strings.TrimSpace(line)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 0 {
			importName := matches[1]
			if importName == "" {
				importName = matches[2]
			}
			imports[importName] = i + 1 // line numbers start at 1
		}
	}

	return imports
}

// toDependency converts a line in go.mod to a dependency
func toDependency(line string) dependency {
	parts := strings.Fields(line)
	first := parts[0]
	second := parts[1]
	version := strings.ReplaceAll(second, indirect, "")
	return dependency{
		name:           first,
		version:        version,
		isIncompatible: strings.Contains(second, incompatible),
		isIndirect:     strings.Contains(second, indirect),
		cGoResults:     []cGOResult{},
	}
}
