// Package netrasp provides an easy way to communicate with network devices
//
// Using an SSH connection it lets you send commands and configure devices
// that only supports screen scraping.
package netrasp

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

var errUserNotSpecified = errors.New("username is not specified")

// Config contains the information needed to connect to a device.
type config struct {
	Platform   Platform
	Connection connection
	SSHConfig  *ssh.ClientConfig
	Host       *host
	driver     string
}

// New creates a new SSH connction to the device.
func New(hostname string, opts ...ConfigOpt) (Platform, error) {
	config := &config{
		SSHConfig: &ssh.ClientConfig{},
		Host: &host{
			Address:  hostname,
			Port:     22,
			password: "",
		},
		driver: "",
	}
	config.SSHConfig.SetDefaults()
	for _, opt := range opts {
		opt.apply(config)
	}
	if len(config.SSHConfig.User) == 0 {
		return nil, errUserNotSpecified
	}

	if config.SSHConfig.HostKeyCallback == nil {
		hostKeyCallback, err := knownHosts(defaultKnownHosts())
		if err != nil {
			return nil, err
		}
		config.SSHConfig.HostKeyCallback = hostKeyCallback
	}

	if config.Connection == nil {
		config.Connection = &sshConnection{Config: config.SSHConfig, Host: config.Host}
	}

	if config.Platform == nil {
		device, err := initDevice(config.driver, config.Connection)
		if err != nil {
			return nil, fmt.Errorf("unable to use selected platform: %w", err)
		}

		return device, nil
	}

	return config.Platform, nil
}

// ConfigOpt configures a Netrasp connection and platform.
type ConfigOpt interface {
	apply(*config)
}

type funcConfigOpt struct {
	f func(*config)
}

func (fco *funcConfigOpt) apply(c *config) {
	fco.f(c)
}

func newFuncConfigOpt(f func(*config)) *funcConfigOpt {
	return &funcConfigOpt{
		f: f,
	}
}

func WithUsernamePassword(username string, password string) ConfigOpt {
	return newFuncConfigOpt(func(c *config) {
		c.SSHConfig.User = username
		c.SSHConfig.Auth = []ssh.AuthMethod{
			ssh.Password(password),
		}
		c.Host.password = password
	})
}

// WithSSHPort allows you to specify an alternate SSH port, defaults to 22.
func WithSSHPort(port int) ConfigOpt {
	return newFuncConfigOpt(func(c *config) {
		c.Host.Port = port
	})
}

// WithInsecureIgnoreHostKey allows you to ignore the validation of the public
// SSH key of a device against a reference in a known_hosts file. Using this
// option should be considered a security risk.
func WithInsecureIgnoreHostKey() ConfigOpt {
	return newFuncConfigOpt(func(c *config) {
		c.SSHConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey() // nolint: gosec
	})
}

// WithDriver tells Netrasp which network driver to use for the connection
// for example "asa", "ios", "nxos".
func WithDriver(name string) ConfigOpt {
	return newFuncConfigOpt(func(c *config) {
		c.driver = name
	})
}

// WithSSHCipher allows you to configure additional SSH Ciphers that the connection
// will use. The parameter can be useful if your device doesn't support the default
// ciphers.
func WithSSHCipher(name string) ConfigOpt {
	return newFuncConfigOpt(func(c *config) {
		c.SSHConfig.Ciphers = append(c.SSHConfig.Ciphers, name)
	})
}
