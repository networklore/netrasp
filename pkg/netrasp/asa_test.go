package netrasp

import (
	"context"
	"testing"
)

func TestFakeASA(t *testing.T) {
	device, err := New("test1",
		WithUsernamePassword("username", "password"),
		WithInsecureIgnoreHostKey(),
		WithDriver("asa"),
		withFakeFileConnection("testdata/asa/basic"),
	)
	if err != nil {
		t.Fatalf("unable to create device: %v", err)
	}

	err = device.Dial(context.Background())
	if err != nil {
		t.Fatalf("could not establish connection. Error: '%v'", err)
	}
	defer device.Close(context.Background())
}
