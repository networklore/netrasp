package netrasp

import (
	"context"
	"strings"
	"testing"
)

func TestFakeSmartaxRunCommand(t *testing.T) {
	device, err := New("test1",
		WithUsernamePassword("username", "password"),
		WithInsecureIgnoreHostKey(),
		WithDriver("smartax"),
		withFakeFileConnection("testdata/smartax/basic"),
	)
	if err != nil {
		t.Fatalf("unable to create device: %v", err)
	}
	defer device.Close(context.Background())

	err = device.Dial(context.Background())
	if err != nil {
		t.Fatalf("could not establish connection. Error: '%v'", err)
	}
	output, err := device.Run(context.Background(), "display version")
	if err != nil {
		t.Fatalf("unable to run command. Error: %v", err)
	}
	want := "VERSION : MA5800V100R015C10"
	if !strings.Contains(output, want) {
		t.Fatalf("expected to find '%s' in '%s'", want, output)
	}
}

func TestFakeSmartaxEnableAndConfigureCommand(t *testing.T) {
	device, err := New("test1",
		WithUsernamePassword("username", "password"),
		WithInsecureIgnoreHostKey(),
		WithDriver("smartax"),
		withFakeFileConnection("testdata/smartax/enable_with_config"),
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

	config := []string{"display current-configuration"}

	_, err = device.Configure(context.Background(), config)

	if err != nil {
		t.Fatalf("could not send config. Error: '%v'", err)
	}
}
