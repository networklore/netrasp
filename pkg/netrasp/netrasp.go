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
type Config struct {
	Platform   Platform
	Connection Connection
	SSHConfig  *ssh.ClientConfig
	Host       *Host
	driver     string
}

// New creates a new SSH connction to the device.
func New(host string, opts ...ConfigOpt) (Platform, error) {
	config := &Config{
		SSHConfig: &ssh.ClientConfig{},
		Host: &Host{
			Address:  host,
			Port:     22,
			password: "",
		},
		driver: "",
	}
	config.SSHConfig.SetDefaults()
	for _, opt := range opts {
		opt.Apply(config)
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

type ConfigOpt interface {
	Apply(*Config)
}

type funcConfigOpt struct {
	f func(*Config)
}

func (fco *funcConfigOpt) Apply(c *Config) {
	fco.f(c)
}

func newFuncConfigOpt(f func(*Config)) *funcConfigOpt {
	return &funcConfigOpt{
		f: f,
	}
}

func WithUsernamePassword(username string, password string) ConfigOpt {
	return newFuncConfigOpt(func(c *Config) {
		c.SSHConfig.User = username
		c.SSHConfig.Auth = []ssh.AuthMethod{
			ssh.Password(password),
		}
		c.Host.password = password
	})
}

func WithSSHPort(port int) ConfigOpt {
	return newFuncConfigOpt(func(c *Config) {
		c.Host.Port = port
	})
}

func WithInsecureIgnoreHostKey() ConfigOpt {
	return newFuncConfigOpt(func(c *Config) {
		c.SSHConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey() // nolint: gosec
	})
}

func WithDriver(name string) ConfigOpt {
	return newFuncConfigOpt(func(c *Config) {
		c.driver = name
	})
}

func WithSSHCipher(name string) ConfigOpt {
	return newFuncConfigOpt(func(c *Config) {
		c.SSHConfig.Ciphers = append(c.SSHConfig.Ciphers, name)
	})
}
