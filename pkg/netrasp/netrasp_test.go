package netrasp

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestFailWithNoUser(t *testing.T) {
	_, err := New("testdevice", WithDriver("ios"))
	if errors.Is(err, errUserNotSpecified) != true {
		t.Fatalf("expencted to encounter error '%v', actual error '%v'", errUserNotSpecified, err)
	}
}

func TestConfigOptions(t *testing.T) {
	_, err := New("testdevice",
		WithDriver("asa"),
		WithUsernamePassword("admin", "password"),
		WithInsecureIgnoreHostKey(),
		WithSSHPort(2222),
		WithSSHCipher("aes128-cbc"),
	)
	if err != nil {
		t.Fatalf("expected device to be initialized: %v", err)
	}
}

func TestNxosTimingRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping nxos integration test in short mode")
	}
	connection := connect(t, "sbx-nxos-mgmt.cisco.com", "nxos", "admin", "Admin_1234!")
	defer connection.Close(context.Background())

	cases := []struct {
		name    string
		command string
		want    string
	}{
		{
			name:    "show_version",
			command: "show version",
			want:    "Cisco Nexus Operating System (NX-OS) Software",
		},
		{
			name:    "show_inventory",
			command: "show inventory",
			want:    "Nexus 9000v",
		},
	}
	err := connection.Enable(context.Background())
	if err != nil {
		t.Fatalf("unable to run enable")
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := connection.Run(context.Background(), tc.command)
			if err != nil {
				t.Fatalf("unable to run command '%s'. Error: %v", tc.command, err)
			}

			if !strings.Contains(result, tc.want) {
				t.Fatalf("expected to find '%s' in '%s'", tc.want, result)
			}
		})
	}
}

func TestIosEnable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ios integration test in short mode")
	}
	connection := connect(t, "ios-xe-mgmt.cisco.com", "ios", "developer", "C1sco12345")
	defer connection.Close(context.Background())

	_, err := connection.Run(context.Background(), "disable")
	if err != nil {
		t.Fatalf("could not send disable command. Error: '%v'", err)
	}
	err = connection.Enable(context.Background())
	if err != nil {
		t.Fatalf("cound not enter enable mode. Error: '%v'", err)
	}
}

func TestIosConfigure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ios integration test in short mode")
	}
	connection := connect(t, "ios-xe-mgmt.cisco.com", "ios", "developer", "C1sco12345")
	defer connection.Close(context.Background())

	config := []string{
		"ip access-list extended netrasp-test",
		"remark running some tests",
		"permit ip any any",
	}

	_, err := connection.Configure(context.Background(), config)

	if err != nil {
		t.Fatalf("could not send config. Error: '%v'", err)
	}

	result, err := connection.Run(context.Background(), "show ip access-list netrasp-test")
	if err != nil {
		t.Fatalf("could run command. Error: '%v'", err)
	}

	if want := "permit ip any any"; !strings.Contains(result, want) {
		t.Fatalf("expected to find '%s' in '%s'", want, result)
	}

	_, err = connection.Configure(context.Background(), []string{"no ip access-list extended netrasp-test"})
	if err != nil {
		t.Fatalf("unable to remove config. Error: '%v'", err)
	}
}

func TestIosTimingRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ios integration test in short mode")
	}
	connection := connect(t, "ios-xe-mgmt.cisco.com", "ios", "developer", "C1sco12345")
	defer connection.Close(context.Background())

	cases := []struct {
		name    string
		command string
		want    string
	}{
		{
			name:    "show_version",
			command: "show version",
			want:    "Cisco IOS Software",
		},
		{
			name:    "show_inventory",
			command: "show inventory",
			want:    "CSR1000V",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := connection.Run(context.Background(), tc.command)
			if err != nil {
				t.Fatalf("unable to run command '%s'. Error: %v", tc.command, err)
			}

			if !strings.Contains(result, tc.want) {
				t.Fatalf("expected to find '%s' in '%s'", tc.want, result)
			}
		})
	}
}

func connect(t *testing.T, host string, platform string, username string, password string) Platform {
	t.Helper()

	device, err := New(host,
		WithUsernamePassword(username, password),
		WithDriver(platform),
		WithInsecureIgnoreHostKey(),
		WithSSHPort(8181),
	)
	if err != nil {
		t.Fatalf("unable to create device. Error: '%v':", err)
	}

	err = device.Dial(context.Background())
	if err != nil {
		t.Fatalf("could not establish connection. Error: '%v'", err)
	}

	return device
}
