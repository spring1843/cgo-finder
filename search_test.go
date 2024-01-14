package main

import "testing"

func TestGoModDir(t *testing.T) {
	got := findGoModDir()
	if len(got) == 0 {
		t.Error("got an empty directory")
	}
}
