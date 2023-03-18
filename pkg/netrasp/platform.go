package netrasp

import (
	"context"
	"errors"
	"regexp"
)

var errInvalidPlatformError = errors.New("invalid platform selected")

// Platform defines an interface for network drivers.
type Platform interface {
	// Disconnect from a device
	Close(context.Context) error
	// Configure device with the provided commands
	Configure(context.Context, []string) (ConfigResult, error)
	// Open a connection to a device
	Dial(context.Context) error
	// Elevate privileges on device
	Enable(context.Context) error
	// Run a command against a device
	Run(context.Context, string) (string, error)
	// Run a command against a device and search for a specific prompt
	RunUntil(context.Context, string, *regexp.Regexp) (string, error)
}

// initDevice returns a platform / network driver.
func initDevice(platform string, conn connection) (Platform, error) {
	switch platform {
	case "asa":
		driver := &asa{}
		driver.Connection = conn

		return driver, nil
	case "nxos":
		driver := &nxos{}
		driver.Connection = conn

		return driver, nil
	case "ios":
		driver := &ios{}
		driver.Connection = conn

		return driver, nil
	case "sros":
		driver := &sros{}
		driver.Connection = conn

		return driver, nil
	case "smartax":
		driver := &smartax{}
		driver.Connection = conn

		return driver, nil
	}

	return nil, errInvalidPlatformError
}
