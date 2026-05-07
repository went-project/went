package utils

import "testing"

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		name     string
		left     string
		right    string
		expected int
	}{
		{"stable greater than beta", "v1.0.0", "v1.0.0-beta.1", 1},
		{"beta less than stable", "v1.0.0-beta.1", "v1.0.0", -1},
		{"beta numeric ordering", "v1.0.0-beta.1", "v1.0.0-beta.2", -1},
		{"beta suffix ordering", "v1.0.0-beta", "v1.0.0-beta.1", -1},
		{"stable patch ordering", "v1.0.0", "v1.0.1", -1},
		{"pre-release alphabetical", "v1.0.0-alpha", "v1.0.0-beta", -1},
		{"mixed numeric and lexical", "v1.10.0", "v1.2.0", 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := CompareVersions(tc.left, tc.right)
			if result != tc.expected {
				t.Fatalf("CompareVersions(%q, %q) = %d, want %d", tc.left, tc.right, result, tc.expected)
			}
		})
	}
}

func TestParseChannelFlag(t *testing.T) {
	if ParseChannelFlag(false) != ChannelStable {
		t.Fatal("expected stable channel when beta flag is false")
	}
	if ParseChannelFlag(true) != ChannelBeta {
		t.Fatal("expected beta channel when beta flag is true")
	}
}
