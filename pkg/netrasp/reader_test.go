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

var ns string

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

func generatePrompt(lines int) string {
	return strings.Repeat(strings.Repeat("dummy", 10)+"\r\r\n", lines-1) + "prompt#"
}

func runBenchmarkReaderBytesBufferNLines(b *testing.B, i int) {
	b.Helper()
	var s string
	var err error
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		r := strings.NewReader(generatePrompt(i))
		s, err = readUntilPrompt(context.Background(), r, generalPrompt)
		if err != nil {
			b.Fail()
		}
	}
	ns = s
}

// 100 lines.
func BenchmarkReaderBytesBuffer100Lines(b *testing.B) {
	runBenchmarkReaderBytesBufferNLines(b, 100)
}

// 1K lines.
func BenchmarkReaderBytesBuffer1KLines(b *testing.B) {
	runBenchmarkReaderBytesBufferNLines(b, 1000)
}

// 10K lines.
func BenchmarkReaderBytesBuffer10KLines(b *testing.B) {
	runBenchmarkReaderBytesBufferNLines(b, 10000)
}

// 50k Lines.
func BenchmarkReaderBytesBuffer50KLines(b *testing.B) {
	runBenchmarkReaderBytesBufferNLines(b, 50000)
}

// 100k Lines.
func BenchmarkReaderBytesBuffer100KLines(b *testing.B) {
	runBenchmarkReaderBytesBufferNLines(b, 100000)
}
