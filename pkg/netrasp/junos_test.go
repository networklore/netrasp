package netrasp

import (
	"context"
	"strings"
	"testing"
)

func TestFakeJunosRunCommand(t *testing.T) {
	device, err := New("test1",
		WithUsernamePassword("username", "Password"),
		WithInsecureIgnoreHostKey(),
		WithDriver("junos"),
		withFakeFileConnection("testdata/junos/basic"),
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
	want := "Junos: 14.1R1.10"
	if !strings.Contains(output, want) {
		t.Fatalf("expected to find '%s' in '%s'", want, output)
	}
}