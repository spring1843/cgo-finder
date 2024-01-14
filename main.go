package main

import "fmt"

func main() {
	dependencies := parseGoModDependencies(readGoMod())
	foundCGo := false
	for _, dependency := range dependencies {
		foundCGo = foundCGo && dependency.checkForCGo()
	}
	if !foundCGo {
		fmt.Println("No CGo dependencies found.")
	}
}
