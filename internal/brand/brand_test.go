package brand

import "testing"

func TestLightPaletteMatchesBrandDeck(t *testing.T) {
	want := map[string]string{
		"Ink":      "#0B0B09",
		"Pine":     "#21663F",
		"Paper":    "#F0EFE9",
		"Moss":     "#49A863",
		"Gold":     "#7A632A",
		"Brass":    "#F0D053",
		"Charcoal": "#3B3F39",
		"Stone":    "#6A6A65",
		"Mist":     "#D9DAD6",
	}
	got := map[string]string{
		"Ink":      LightInk,
		"Pine":     LightPine,
		"Paper":    LightPaper,
		"Moss":     LightMoss,
		"Gold":     LightGold,
		"Brass":    LightBrass,
		"Charcoal": LightCharcoal,
		"Stone":    LightStone,
		"Mist":     LightMist,
	}
	for name, w := range want {
		if got[name] != w {
			t.Errorf("light %s = %s, want %s", name, got[name], w)
		}
	}
}

func TestDarkPaletteMatchesBrandDeck(t *testing.T) {
	want := map[string]string{
		"Black":    "#000000",
		"Moss":     "#49A863",
		"Paper":    "#F7F7F4",
		"Gold":     "#C9A855",
		"Brass":    "#F5D96E",
		"Stone":    "#888888",
		"Charcoal": "#A8A8A4",
		"Mist":     "#8A8A8A",
	}
	got := map[string]string{
		"Black":    DarkBlack,
		"Moss":     DarkMoss,
		"Paper":    DarkPaper,
		"Gold":     DarkGold,
		"Brass":    DarkBrass,
		"Stone":    DarkStone,
		"Charcoal": DarkCharcoal,
		"Mist":     DarkMist,
	}
	for name, w := range want {
		if got[name] != w {
			t.Errorf("dark %s = %s, want %s", name, got[name], w)
		}
	}
}

func TestFontsMatchBrandDeck(t *testing.T) {
	fonts := map[string]string{
		"Display": FontDisplay,
		"Heading": FontHeading,
		"Body":    FontBody,
		"Meta":    FontMeta,
	}
	want := map[string]string{
		"Display": "Big Shoulders Inline Display",
		"Heading": "Cinzel",
		"Body":    "Google Sans Flex",
		"Meta":    "Google Sans Code",
	}
	for name, w := range want {
		if fonts[name] != w {
			t.Errorf("font %s = %q, want %q", name, fonts[name], w)
		}
	}
}

func TestPairingsAreWellFormed(t *testing.T) {
	for _, pairings := range []struct {
		theme string
		ps    []Pairing
	}{{"light", LightPairings()}, {"dark", DarkPairings()}} {
		if len(pairings.ps) == 0 {
			t.Errorf("%s theme declares no pairings", pairings.theme)
		}
		for _, p := range pairings.ps {
			if p.Name == "" {
				t.Errorf("%s: pairing with empty name", pairings.theme)
			}
			if p.Min != 4.5 && p.Min != 3.0 {
				t.Errorf("%s/%s: Min = %v, want 4.5 (body) or 3.0 (large/UI)", pairings.theme, p.Name, p.Min)
			}
		}
	}
}
