package netrasp

import (
	"context"
	"strings"
	"testing"
)

func TestFakeIOSRunCommand(t *testing.T) {
	device, err := New("test1",
		WithUsernamePassword("username", "password"),
		WithInsecureIgnoreHostKey(),
		WithDriver("ios"),
		withFakeFileConnection("testdata/ios/basic"),
	)
	if err != nil {
		t.Fatalf("unable to create device: %v", err)
	}
	defer device.Close(context.Background())

	err = device.Dial(context.Background())
	if err != nil {
		t.Fatalf("could not establish connection. Error: '%v'", err)
	}
	output, err := device.Run(context.Background(), "show version")
	if err != nil {
		t.Fatalf("unable to run command. Error: %v", err)
	}
	want := "Cisco IOS XE Software, Version 16.09.03"
	if !strings.Contains(output, want) {
		t.Fatalf("expected to find '%s' in '%s'", want, output)
	}
}

func TestFakeIOSEnableAndConfigureCommand(t *testing.T) {
	device, err := New("test1",
		WithUsernamePassword("username", "password"),
		WithInsecureIgnoreHostKey(),
		WithDriver("ios"),
		withFakeFileConnection("testdata/ios/enable_with_config"),
	)
	if err != nil {
		t.Fatalf("unable to create device: %v", err)
	}
	defer device.Close(context.Background())

	err = device.Dial(context.Background())
	if err != nil {
		t.Fatalf("could not establish connection. Error: '%v'", err)
	}

	err = device.Enable(context.Background())
	if err != nil {
		t.Fatalf("could not establish connection. Error: '%v'", err)
	}

	config := []string{
		"ip access-list extended netrasp-test",
		"remark running some tests",
		"permit ip any any",
	}

	_, err = device.Configure(context.Background(), config)

	if err != nil {
		t.Fatalf("could not send config. Error: '%v'", err)
	}
}
