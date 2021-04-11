package netrasp

import (
	"testing"
)

func TestInvalidKnownHostsFiles(t *testing.T) {
	cases := []struct {
		name string
		file []string
	}{
		{
			name: "invalid_file",
			file: []string{"testdata/invalid_known_hosts"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := knownHosts(tc.file)
			if err == nil {
				t.Fatalf("parsing invalid file '%s' was expected to fail but didn't", tc.file)
			}
		})
	}
}
