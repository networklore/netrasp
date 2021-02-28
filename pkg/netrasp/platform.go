package netrasp

import (
	"context"
	"errors"
	"regexp"
)

var errInvalidPlatformError = errors.New("invalid platform selected")

type Platform interface {
	// Disconnect from a device
	Close(context.Context) error
	// Configure device with the provided commands
	Configure(context.Context, []string) (string, error)
	// Open a connection to a device
	Dial(context.Context) error
	// Elevate privileges on device
	Enable(context.Context) error
	// Run a command against a device
	Run(context.Context, string) (string, error)
	// Run a command against a device and search for a specific prompt
	RunUntil(context.Context, string, *regexp.Regexp) (string, error)
}

func InitDevice(platform string, connection Connection) (Platform, error) {
	switch platform {
	case "asa":
		driver := &Asa{}
		driver.Connection = connection

		return driver, nil
	case "nxos":
		driver := &Nxos{}
		driver.Connection = connection

		return driver, nil
	case "ios":
		driver := &Ios{}
		driver.Connection = connection

		return driver, nil
	}

	return nil, errInvalidPlatformError
}
