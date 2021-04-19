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

func runBenchmarkReaderOriginalNLines(b *testing.B, i int) {
	b.Helper()
	var s string
	var err error

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		r := strings.NewReader(generatePrompt(i))
		s, err = readUntilPromptOriginal(context.Background(), r, generalPrompt)
		if err != nil {
			b.Fail()
		}
	}
	ns = s
}

func runBenchmarkReaderStringsBuilderNLines(b *testing.B, i int) {
	b.Helper()
	var s string
	var err error
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		r := strings.NewReader(generatePrompt(i))
		s, err = readUntilPromptWithStringsBuilder(context.Background(), r, generalPrompt)
		if err != nil {
			b.Fail()
		}
	}
	ns = s
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

// 1K lines.
func BenchmarkReaderOriginal100Lines(b *testing.B) {
	runBenchmarkReaderOriginalNLines(b, 100)
}
func BenchmarkReaderStringsBuilder100Lines(b *testing.B) {
	runBenchmarkReaderStringsBuilderNLines(b, 100)
}
func BenchmarkReaderBytesBuffer100Lines(b *testing.B) {
	runBenchmarkReaderBytesBufferNLines(b, 100)
}

// 1K lines.
func BenchmarkReaderOriginal1KLines(b *testing.B) {
	runBenchmarkReaderOriginalNLines(b, 1000)
}
func BenchmarkReaderStringsBuilder1KLines(b *testing.B) {
	runBenchmarkReaderStringsBuilderNLines(b, 1000)
}
func BenchmarkReaderBytesBuffer1KLines(b *testing.B) {
	runBenchmarkReaderBytesBufferNLines(b, 1000)
}

// 10K lines.
func BenchmarkReaderOriginal10KLines(b *testing.B) {
	runBenchmarkReaderOriginalNLines(b, 10000)
}
func BenchmarkReaderStringsBuilder10KLines(b *testing.B) {
	runBenchmarkReaderStringsBuilderNLines(b, 10000)
}
func BenchmarkReaderBytesBuffer10KLines(b *testing.B) {
	runBenchmarkReaderBytesBufferNLines(b, 10000)
}

// 50k Lines.
func BenchmarkReaderOriginal50KLines(b *testing.B) {
	runBenchmarkReaderOriginalNLines(b, 50000)
}
func BenchmarkReaderStringsBuilder50KLines(b *testing.B) {
	runBenchmarkReaderStringsBuilderNLines(b, 50000)
}
func BenchmarkReaderBytesBuffer50KLines(b *testing.B) {
	runBenchmarkReaderBytesBufferNLines(b, 50000)
}

// 100k Lines.
func BenchmarkReaderOriginal100KLines(b *testing.B) {
	runBenchmarkReaderOriginalNLines(b, 100000)
}

func BenchmarkReaderStringsBuilder100KLines(b *testing.B) {
	runBenchmarkReaderStringsBuilderNLines(b, 100000)
}

func BenchmarkReaderBytesBuffer100KLines(b *testing.B) {
	runBenchmarkReaderBytesBufferNLines(b, 100000)
}
