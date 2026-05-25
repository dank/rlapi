package rlapi

import "testing"

func TestDecodeBuildID(t *testing.T) {
	tests := []struct {
		input string
		want  int32
	}{
		{"260316.80791.512269", 1210528741},
		{"260420.86069.515605", 1273328361},
		{"260506.26700.517210", -1652286008},
	}
	for _, tt := range tests {
		got := decodeBuildID(tt.input)
		if got != tt.want {
			t.Errorf("decodeBuildID(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
