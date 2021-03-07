//+build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"os"
)

// Run Integrationtests
func Integration() error {
	params := []string{
		"test",
		"-v",
		"./...",
		"-p",
		"1",
		"-coverprofile=cover.out",
		"-coverprofile=coverage.txt",
		"-covermode=atomic",
	}
	return sh.RunV("go", params...)
}

// Run Unittests
func Test() error {
	params := []string{
		"test",
		"-v",
		"./...",
		"-short",
		"-coverprofile=cover.out",
		"-coverprofile=coverage.txt",
		"-covermode=atomic",
	}
	return sh.RunV("go", params...)
}

// Run linting checks
func Lint() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	params := []string{
		"run",
		"--rm",
		"-v",
		fmt.Sprintf("%s:/go/src/netrasp", cwd),
		"-w",
		"/go/src/netrasp",
		"-e",
		"GO111MODULE=on",
		"-e",
		"GOPROXY=https://proxy.golang.org",
		"golangci/golangci-lint:v1.36",
		"golangci-lint",
		"run",
		"--timeout",
		"300s",
	}
	return sh.RunV("docker", params...)
}
