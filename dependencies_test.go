package main

import (
	"reflect"
	"testing"
)

func TestParseDependencies(t *testing.T) {
	goModContent := `module github.com/foo/bar
go 1.23

require (
	bitbucket.org/a/b v0.5.0
	github.com/c/d v68.0.0+incompatible
	github.com/e/f v0.11.29
	github.com/g/h v0.9.23
)

require (
	cloud.google.com/i/j v1.23.0 // indirect
	cloud.google.com/l/m v0.2.3 // indirect
)`

	dependencies := parseGoModDependencies(goModContent)
	want := []dependency{
		{name: "bitbucket.org/a/b", version: "v0.5.0", isIncompatible: false, isIndirect: false, cGoResults: []cGOResult{}},
		{name: "github.com/c/d", version: "v68.0.0+incompatible", isIncompatible: true, isIndirect: false, cGoResults: []cGOResult{}},
		{name: "github.com/e/f", version: "v0.11.29", isIncompatible: false, isIndirect: false, cGoResults: []cGOResult{}},
		{name: "github.com/g/h", version: "v0.9.23", isIncompatible: false, isIndirect: false, cGoResults: []cGOResult{}},
		{name: "cloud.google.com/i/j", version: "v1.23.0", isIncompatible: false, isIndirect: false, cGoResults: []cGOResult{}},
		{name: "cloud.google.com/l/m", version: "v0.2.3", isIncompatible: false, isIndirect: false, cGoResults: []cGOResult{}},
	}

	if !reflect.DeepEqual(dependencies, want) {
		t.Errorf("got %#v want %#v", dependencies, want)
	}
}

func TestParseImports(t *testing.T) {
	tests := []struct {
		input   string
		imports map[string]int
	}{
		{
			`// package main is a test package		
package main

import "a"
import "b"

import "c"`,
			map[string]int{"a": 4, "b": 5, "c": 7},
		},
		{
			`package main

import (
	"a"

	"b"

	"c"
)

type foo struct {
	a int
}`,
			map[string]int{"a": 4, "b": 6, "c": 8},
		},
		{
			`package main

import "a"

import (
	"b"

	"c"
)`,
			map[string]int{"a": 3, "b": 6, "c": 8},
		},
		{
			`package main

import (
	"a"

	"b"
)

import "c"
`,
			map[string]int{"a": 4, "b": 6, "c": 9},
		},
	}

	for i, test := range tests {
		got := parseImports(test.input)
		if !reflect.DeepEqual(got, test.imports) {
			t.Errorf("failed test %d, got %#v want %#v", i, got, test.imports)
		}
	}
}

func TestExclamationFormat(t *testing.T) {
	got := exclamationFormat("FooBarBAZ")
	want := "!foo!bar!b!a!z"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
