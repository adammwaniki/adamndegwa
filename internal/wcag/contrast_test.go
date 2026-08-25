package wcag

import (
	"math"
	"testing"
)

func TestParseHex(t *testing.T) {
	r, g, b, err := ParseHex("#21663F")
	if err != nil {
		t.Fatal(err)
	}
	if r != 0x21 || g != 0x66 || b != 0x3F {
		t.Errorf("got %02x %02x %02x", r, g, b)
	}
	for _, bad := range []string{"", "#12345", "21663F", "#GGGGGG", "#1234567"} {
		if _, _, _, err := ParseHex(bad); err == nil {
			t.Errorf("ParseHex(%q) should fail", bad)
		}
	}
}

func TestRelativeLuminanceKnownValues(t *testing.T) {
	cases := map[string]float64{
		"#000000": 0.0,
		"#FFFFFF": 1.0,
	}
	for hex, want := range cases {
		got, err := RelativeLuminance(hex)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("L(%s) = %v, want %v", hex, got, want)
		}
	}
	if _, err := RelativeLuminance("nope"); err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestContrastRatioKnownPairs(t *testing.T) {
	// WCAG's canonical example: black on white is 21:1, a colour on itself 1:1.
	if r, _ := ContrastRatio("#000000", "#FFFFFF"); math.Abs(r-21) > 0.01 {
		t.Errorf("black/white = %v, want 21", r)
	}
	if r, _ := ContrastRatio("#21663F", "#21663F"); math.Abs(r-1) > 0.001 {
		t.Errorf("same colour = %v, want 1", r)
	}
	// Ratio is symmetric.
	a, _ := ContrastRatio("#0B0B09", "#F7F7F4")
	b, _ := ContrastRatio("#F7F7F4", "#0B0B09")
	if math.Abs(a-b) > 1e-9 {
		t.Errorf("ratio not symmetric: %v vs %v", a, b)
	}
	if _, err := ContrastRatio("bad", "#FFFFFF"); err == nil {
		t.Error("expected error for invalid fg")
	}
	if _, err := ContrastRatio("#FFFFFF", "bad"); err == nil {
		t.Error("expected error for invalid bg")
	}
}
