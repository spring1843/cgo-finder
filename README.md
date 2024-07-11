# CGO-Finder

CGO-Finder is a tool that helps identify and analyze CGo usage in Go project dependencies, including indirect dependencies.

## Background

Using dependencies that utilize CGo or have dependencies that rely on CGo can introduce challenges in terms of deployability and portability of the code. While it may be necessary to use such dependencies in certain cases, it is generally recommended to prioritize pure Go alternatives whenever possible.

Some of the negative consequences of relying on CGo dependencies include:

1. Inability to cross-compile the code for different platforms.
2. Difficulty in creating minimal container images for deployment.

There is no straight forward way to identify which dependencies 

## Why cgo-finder?

It's not easy to find dependencies that use CGo. This tool was made to address this need.

## Usage

1. Change current directory to the one that contains your code base and `go.mod`
Make sure all dependencies are present by running `go mod download`
2. Install cgo-finder: `go install github.com/spring1843/cgo-finder@latest`
3. Run the tool to generate report `cgo-finder`
