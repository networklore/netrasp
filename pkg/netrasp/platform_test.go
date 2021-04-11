package netrasp

import (
	"testing"
)

func TestNewDevice(t *testing.T) {
	cases := []struct {
		name     string
		platform string
		succeed  bool
	}{
		{
			name:     "cisco_asa",
			platform: "asa",
			succeed:  true,
		},
		{
			name:     "cisco_ios",
			platform: "ios",
			succeed:  true,
		},
		{
			name:     "cisco_nxos",
			platform: "nxos",
			succeed:  true,
		},
		{
			name:     "driver_missing",
			platform: "",
			succeed:  false,
		},
		{
			name:     "bogus_driver",
			platform: "doesn_not_exist",
			succeed:  false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := New("device1", WithUsernamePassword("admin", "password"), WithDriver(tc.platform))
			status := err == nil
			if status != tc.succeed {
				t.Fatalf("platform '%s' was expected to return '%v', actual result '%v'", tc.platform, tc.succeed, status)
			}
		})
	}
}
