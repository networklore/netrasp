package netrasp

import (
	"context"
	"testing"
)

func TestFakeNxos(t *testing.T) {
	device, err := New("test1",
		WithUsernamePassword("username", "password"),
		WithInsecureIgnoreHostKey(),
		WithDriver("nxos"),
		withFakeFileConnection("testdata/nxos/basic"),
	)
	if err != nil {
		t.Fatalf("unable to create device: %v", err)
	}

	err = device.Dial(context.Background())
	if err != nil {
		t.Fatalf("could not establish connection. Error: '%v'", err)
	}

	err = device.Enable(context.Background())
	if err != nil {
		t.Fatalf("enable should always work on nxos. Error: '%v'", err)
	}

	defer device.Close(context.Background())
}
