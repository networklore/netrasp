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

	got, err := ReadUntilPrompt(context.Background(), r, generalPrompt)
	if err != nil {
		t.Fatalf("error reading until prompt: %v", err)
	}
	if got != testReaderBasic {
		t.Fatalf("expected '%s' got '%s'", testReaderBasic, got)
	}
}

func TestReaderTimeout(t *testing.T) {
	r := strings.NewReader(testReaderBasic)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	_, err := ReadUntilPrompt(ctx, r, generalPrompt)
	if err == nil {
		t.Fatalf("expected read to fail")
	}
}
