package main

import (
	"slices"
	"testing"
)

// Found while screencapping: "Münzen" rendered as "M?nzen" because only ASCII
// glyphs were loaded.
func TestFontRunesCoverUmlauts(t *testing.T) {
	rs := FontRunes()
	for _, r := range "äöüÄÖÜ–" {
		if !slices.Contains(rs, r) {
			t.Errorf("glyph %q missing from font set", r)
		}
	}
}
