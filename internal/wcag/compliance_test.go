package wcag

import (
	"os"
	"strings"
	"testing"

	"adamndegwa/internal/brand"
)

// TestPaletteCompliance is the WCAG2 AA proof: every declared text pairing in
// both themes must meet its minimum contrast ratio.
func TestPaletteCompliance(t *testing.T) {
	for _, group := range []struct {
		theme    string
		pairings []brand.Pairing
	}{{"light", brand.LightPairings()}, {"dark", brand.DarkPairings()}} {
		for _, p := range group.pairings {
			ratio, err := ContrastRatio(p.FG, p.BG)
			if err != nil {
				t.Fatalf("%s/%s: %v", group.theme, p.Name, err)
			}
			if ratio < p.Min {
				t.Errorf("%s/%s: contrast %.2f:1 below required %.1f:1 (%s on %s)",
					group.theme, p.Name, ratio, p.Min, p.FG, p.BG)
			}
		}
	}
}

// TestStylesheetTokensMatchBrand ensures static/style.css stays in sync with
// the brand constants — the CSS is where the palette is actually applied.
func TestStylesheetTokensMatchBrand(t *testing.T) {
	css, err := os.ReadFile("../../static/style.css")
	if err != nil {
		t.Fatal(err)
	}
	tokens := map[string]string{
		"--ink":      brand.LightInk,
		"--paper":    brand.LightPaper,
		"--pine":     brand.LightPine,
		"--moss":     brand.LightMoss,
		"--gold":     brand.LightGold,
		"--brass":    brand.LightBrass,
		"--charcoal": brand.LightCharcoal,
		"--stone":    brand.LightStone,
		"--mist":     brand.LightMist,
	}
	for name, hex := range tokens {
		if !strings.Contains(string(css), name+":"+hex) {
			t.Errorf("style.css light token %s should be %s", name, hex)
		}
	}
	// Dark overrides live inside the [data-theme="dark"] block.
	darkIdx := strings.Index(string(css), `[data-theme="dark"]{`)
	if darkIdx < 0 {
		t.Fatal("style.css has no dark theme block")
	}
	darkBlock := string(css)[darkIdx:]
	darkTokens := map[string]string{
		"--ink":      brand.DarkPaper,
		"--paper":    brand.DarkBlack,
		"--pine":     brand.DarkMoss,
		"--moss":     brand.DarkMoss,
		"--gold":     brand.DarkGold,
		"--brass":    brand.DarkBrass,
		"--charcoal": brand.DarkCharcoal,
		"--stone":    brand.DarkStone,
		"--mist":     brand.DarkMist,
	}
	for name, hex := range darkTokens {
		if !strings.Contains(darkBlock, name+":"+hex) {
			t.Errorf("style.css dark token %s should be %s", name, hex)
		}
	}
}

// TestStylesheetA11yFeatures asserts the CSS carries the features WCAG
// requires beyond colour: visible focus, reduced motion, a skip link.
func TestStylesheetA11yFeatures(t *testing.T) {
	css, err := os.ReadFile("../../static/style.css")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		":focus-visible",
		"prefers-reduced-motion",
		".skip-link",
	} {
		if !strings.Contains(string(css), want) {
			t.Errorf("style.css missing %s", want)
		}
	}
}
