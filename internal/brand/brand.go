package brand

// Brand palette from references/brand-wcag2.pptx. These constants are the
// single source of truth; static/style.css must match them (enforced by
// tests in internal/wcag).

// Light theme (manga paper canvas).
const (
	LightInk      = "#0B0B09" // primary text, dark backgrounds
	LightPine     = "#21663F" // brand color, secondary
	LightPaper    = "#F7F7F4" // base canvas
	LightMoss     = "#49A863" // accent: highlights, CTA
	LightGold     = "#876E2F" // accent: details, lines (WCAG AA)
	LightBrass    = "#F0D053" // spark: rare moments
	LightCharcoal = "#3B3F39" // neutral
	LightStone    = "#72726D" // neutral (WCAG AA)
	LightMist     = "#D9DAD6" // neutral
)

// Dark theme (AMOLED: pure black canvas, OLED pixel-off).
const (
	DarkBlack    = "#000000" // background
	DarkMoss     = "#49A863" // primary on black
	DarkPaper    = "#F7F7F4" // primary text
	DarkGold     = "#C9A855" // accent: labels, lines
	DarkBrass    = "#F5D96E" // spark: hover moments
	DarkStone    = "#888888" // tertiary, meta text (WCAG AA)
	DarkCharcoal = "#A8A8A4" // neutral
	DarkMist     = "#8A8A8A" // neutral
)

// Brand fonts, one per job (from the pptx type system).
const (
	FontDisplay = "Big Shoulders Inline Display" // hero names, tile titles
	FontHeading = "Cinzel"                       // logo, role labels, headings
	FontBody    = "Google Sans Flex"             // descriptions, body copy
	FontMeta    = "Google Sans Code"             // nav, labels, numbers
)

// Pairing declares a foreground/background text combination in use, with the
// minimum WCAG contrast ratio it must satisfy (4.5 body text, 3.0 large text
// and UI components).
type Pairing struct {
	Name string
	FG   string
	BG   string
	Min  float64
}

// LightPairings returns every text pairing used by the light theme.
func LightPairings() []Pairing {
	return []Pairing{
		{"body ink on paper", LightInk, LightPaper, 4.5},
		{"charcoal on paper", LightCharcoal, LightPaper, 4.5},
		{"stone on paper", LightStone, LightPaper, 4.5},
		{"pine on paper", LightPine, LightPaper, 4.5},
		{"gold meta on paper", LightGold, LightPaper, 4.5},
		{"paper on ink (selection/inverse)", LightPaper, LightInk, 4.5},
		{"paper on pine (tile hover)", LightPaper, LightPine, 4.5},
		{"mist on pine (tile hover desc)", LightMist, LightPine, 4.5},
		{"brass on pine (tile hover accent)", LightBrass, LightPine, 4.5},
		{"brass on ink (large spark)", LightBrass, LightInk, 3.0},
	}
}

// DarkPairings returns every text pairing used by the dark (AMOLED) theme.
func DarkPairings() []Pairing {
	return []Pairing{
		{"paper on black", DarkPaper, DarkBlack, 4.5},
		{"charcoal on black", DarkCharcoal, DarkBlack, 4.5},
		{"stone on black", DarkStone, DarkBlack, 4.5},
		{"moss on black", DarkMoss, DarkBlack, 4.5},
		{"gold meta on black", DarkGold, DarkBlack, 4.5},
		{"black on moss (selection/inverse)", DarkBlack, DarkMoss, 4.5},
		{"brass on black (hover spark)", DarkBrass, DarkBlack, 4.5},
		{"mist on black", DarkMist, DarkBlack, 4.5},
		{"black on moss (card/tile hover)", DarkBlack, DarkMoss, 4.5},
	}
}
