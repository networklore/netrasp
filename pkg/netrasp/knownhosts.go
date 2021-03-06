package netrasp

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func defaultKnownHosts() []string {
	hostFiles := []string{"/etc/ssh/ssh_known_hosts"}
	userHome, err := os.UserHomeDir()
	if err == nil {
		hostFiles = append(hostFiles, fmt.Sprintf("%s/.ssh/known_hosts", userHome))
	}

	return hostFiles
}

// KnownHosts loads ssh known_hosts from default locations.
func knownHosts(hostFiles []string) (ssh.HostKeyCallback, error) {
	var existingFiles []string
	for _, hostFile := range hostFiles {
		_, err := os.Stat(hostFile)
		if err == nil {
			existingFiles = append(existingFiles, hostFile)
		}
	}

	callback, err := knownhosts.New(existingFiles...)
	if err != nil {
		return nil, fmt.Errorf("unable to parse known_hosts files: %w", err)
	}

	return callback, nil
}
