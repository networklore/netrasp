package netrasp

import (
	"errors"
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var errKnownHostsMissingError = errors.New("unable to locate a known_hosts file")
var errHomeDirLookupFailure = errors.New("unable to locate user home directory to check for ssh known hosts file")

func missingHostFile(_ string, _ net.Addr, _ ssh.PublicKey) error {
	return errKnownHostsMissingError
}

// KnownHosts loads ssh known_hosts from default locations.
func KnownHosts() (ssh.HostKeyCallback, error) {
	userHome, err := os.UserHomeDir()

	if err != nil {
		return nil, errHomeDirLookupFailure
	}

	var hostFiles []string
	hostFileCandidates := []string{
		fmt.Sprintf("%s/.ssh/known_hosts", userHome),
		"/etc/ssh/ssh_known_hosts",
	}
	for _, hostFile := range hostFileCandidates {
		_, err := os.Stat(hostFile)
		if err == nil {
			hostFiles = append(hostFiles, hostFile)
		}
	}

	if len(hostFiles) == 0 {
		// If we can't locate a known_hosts file return a separate callback that won't
		// error out immediately, the reason for this is that some users might still
		// want to change HostKeyCallback to ssh.InsecureIgnoreHostKey after using
		// netrasp.New()
		return missingHostFile, nil
	}

	callback, err := knownhosts.New(hostFiles...)
	if err != nil {
		return nil, fmt.Errorf("unable to parse known_hosts files: %w", err)
	}

	return callback, nil
}
