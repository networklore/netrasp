package netrasp

import (
	"context"
	"strings"
	"testing"
	"time"
)

var testReaderBasic = `
here is your string
router1#`

func TestReader(t *testing.T) {
	r := strings.NewReader(testReaderBasic)

	got, err := readUntilPrompt(context.Background(), r, generalPrompt)
	if err != nil {
		t.Fatalf("error reading until prompt: %v", err)
	}
	if got != testReaderBasic {
		t.Fatalf("expected '%s' got '%s'", testReaderBasic, got)
	}
}

func TestReaderTimeout(t *testing.T) {
	r := strings.NewReader(testReaderBasic)

	ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
	defer cancel()
	out, err := readUntilPrompt(ctx, r, generalPrompt)
	if err == nil {
		t.Fatalf("expected reader to timeout, but got %v", out)
	}
}
