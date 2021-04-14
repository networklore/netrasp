package netrasp

import (
	"context"
	"strings"
	"testing"
)

func TestFakeSROSRunCommand(t *testing.T) {
	device, err := New("test1",
		WithUsernamePassword("username", "password"),
		WithInsecureIgnoreHostKey(),
		WithDriver("sros"),
		withFakeFileConnection("testdata/sros/basic"),
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
	want := "TiMOS-B-21.2.R1 both/x86_64 Nokia 7750 SR Copyright (c) 2000-2021 Nokia."
	if !strings.Contains(output, want) {
		t.Fatalf("expected to find '%s' in '%s'", want, output)
	}
}
