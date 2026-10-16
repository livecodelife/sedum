package since

import (
	"testing"
	"time"
)

func TestParseSince(t *testing.T) {
	now := time.Date(2026, 10, 16, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		in   string
		want time.Time
		bad  bool
	}{
		{"", time.Time{}, false},
		{"7d", time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC), false},
		{"36h", time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC), false},
		{"2026-10-01", time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{"2026-10-01T06:00:00Z", time.Date(2026, 10, 1, 6, 0, 0, 0, time.UTC), false},
		{"yesterday", time.Time{}, true},
	}
	for _, tc := range cases {
		got, err := parseSince(tc.in, now)
		if (err != nil) != tc.bad {
			t.Errorf("parseSince(%q) error = %v, want error %v", tc.in, err, tc.bad)
			continue
		}
		if !tc.bad && !got.Equal(tc.want) {
			t.Errorf("parseSince(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
