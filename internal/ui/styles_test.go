package ui

import (
	"strings"
	"testing"
)

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
	}

	for _, c := range cases {
		out := FormatBytes(c.bytes)
		if out != c.expected {
			t.Errorf("FormatBytes(%d) = %s; want %s", c.bytes, out, c.expected)
		}
	}
}

func TestFormatSpeed(t *testing.T) {
	if s := FormatSpeed(0); s != "0 B/s" {
		t.Errorf("expected '0 B/s', got '%s'", s)
	}
	if s := FormatSpeed(10485760); s != "10.00 MB/s" {
		t.Errorf("expected '10.00 MB/s', got '%s'", s)
	}
}

func TestRenderMultiSegmentProgress(t *testing.T) {
	bar := RenderMultiSegmentProgress(40, 0.5, 4)
	if !strings.Contains(bar, "█") {
		t.Errorf("expected filled bar block, got: %s", bar)
	}
	if !strings.Contains(bar, "░") {
		t.Errorf("expected empty bar block, got: %s", bar)
	}
}
