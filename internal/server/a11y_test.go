package server

import (
	"os"
	"strings"
	"testing"
)

// a11yPaths covers every page type the site renders.
var a11yPaths = []string{
	"/",
	"/technical-notes",
	"/projects",
	"/musings",
	"/the-bullshitters",
	"/technical-notes/htmx-patterns",
	"/tags",
	"/tags/go",
	"/contact",
}

// TestAccessibilityLandmarks runs WCAG-relevant structural checks against
// every rendered page.
func TestAccessibilityLandmarks(t *testing.T) {
	s := newTestServer(t)
	for _, path := range a11yPaths {
		rec := get(t, s, path, nil)
		if rec.Code != 200 {
			t.Fatalf("%s: status = %d", path, rec.Code)
		}
		body := rec.Body.String()

		if !strings.Contains(body, `<html lang="en">`) {
			t.Errorf("%s: missing lang attribute", path)
		}
		if n := strings.Count(body, "<h1"); n != 1 {
			t.Errorf("%s: found %d h1 elements, want exactly 1", path, n)
		}
		for _, want := range []string{
			`class="skip-link" href="#page-wrapper"`, // skip to content
			`<nav aria-label="Main">`,                // labelled nav landmark
			`<main id="page-wrapper">`,               // main landmark
			`<footer>`,                               // footer landmark
			`aria-label="Toggle dark mode"`,          // named icon buttons
			`aria-label="Open menu"`,
			`aria-expanded="false"`, // disclosure state wired
			`aria-hidden="true"`,    // decorative SVGs hidden from AT
		} {
			if !strings.Contains(body, want) {
				t.Errorf("%s: missing %s", path, want)
			}
		}
		// No placeholder links.
		if strings.Contains(body, `href="#"`) {
			t.Errorf("%s: contains placeholder href=# link", path)
		}
	}
}

// TestLayoutHeadA11y checks document-level features on the home page.
func TestLayoutHeadA11y(t *testing.T) {
	s := newTestServer(t)
	body := get(t, s, "/", nil).Body.String()
	for _, want := range []string{
		`name="viewport" content="width=device-width, initial-scale=1.0"`, // zoom not disabled
		"display=swap",               // fonts don't block invisible text
		`prefers-color-scheme: dark`, // respects OS theme
	} {
		if !strings.Contains(body, want) {
			t.Errorf("layout missing %q", want)
		}
	}
}

// TestThemeScript verifies automatic dark/light mode: initial
// prefers-color-scheme detection guarded against localStorage failures, and a
// live listener for OS theme changes when the user has not chosen explicitly.
func TestThemeScript(t *testing.T) {
	s := newTestServer(t)
	body := get(t, s, "/", nil).Body.String()
	for _, want := range []string{
		"prefers-color-scheme: dark",
		"addEventListener('change'", // live OS theme tracking
		"try {",                     // localStorage can throw (privacy modes); must not kill detection
		"localStorage.getItem('theme')",
		"localStorage.setItem('theme'",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("layout theme script missing %q", want)
		}
	}
}

// TestFooterPinnedToBottom asserts the CSS that keeps the footer at the bottom
// of the viewport on short pages: main must grow (flex:1) and be a flex column
// so the footer's margin-top:auto works.
func TestFooterPinnedToBottom(t *testing.T) {
	css, err := os.ReadFile("../../static/style.css")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"#page-wrapper", "flex:1", "flex-direction:column", "margin-top:auto",
	} {
		if !strings.Contains(string(css), want) {
			t.Errorf("style.css missing %q for sticky footer", want)
		}
	}
}

// TestMobileNavToggleGrouping: on small screens the dark-mode toggle and the
// hamburger must sit adjacent at the right edge (dark toggle pushed right),
// not spread apart by space-between.
func TestMobileNavToggleGrouping(t *testing.T) {
	css, err := os.ReadFile("../../static/style.css")
	if err != nil {
		t.Fatal(err)
	}
	mediaIdx := strings.LastIndex(string(css), "@media (max-width:900px)")
	if mediaIdx < 0 {
		t.Fatal("mobile media query missing")
	}
	block := string(css)[mediaIdx:]
	if !strings.Contains(block, ".dark-toggle{margin-left:auto}") {
		t.Error("mobile media block should pin .dark-toggle next to .nav-toggle with margin-left:auto")
	}
}

// TestContactValueFont: contact row values (email, GitHub handle, LinkedIn
// name) use Google Sans Flex, not the Cinzel display font.
func TestContactValueFont(t *testing.T) {
	css, err := os.ReadFile("../../static/style.css")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(css), ".contact-value{\n  font-family:'Google Sans Flex'") {
		t.Error(".contact-value should be set in Google Sans Flex")
	}
}

// TestDarkHoverTextIsReadable: on hover, cards and tiles invert to the accent
// background (--pine, which is Moss in dark mode). Mist/Brass text fails WCAG
// on that background; dark mode must switch hover text to the dark theme's
// black (--paper) instead.
func TestDarkHoverTextIsReadable(t *testing.T) {
	css, err := os.ReadFile("../../static/style.css")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`[data-theme="dark"] .tile:hover p{color:var(--paper)}`,
		`[data-theme="dark"] .card:hover .card-desc{color:var(--paper)}`,
		`[data-theme="dark"] .card:hover .card-meta{color:var(--paper)}`,
	} {
		if !strings.Contains(string(css), want) {
			t.Errorf("style.css missing dark hover override %q", want)
		}
	}
	// The old failing override must be gone.
	if strings.Contains(string(css), `[data-theme="dark"] .tile:hover p{color:#D9DAD6}`) {
		t.Error("old mist-on-moss dark hover override should be removed")
	}
}

// TestMobileNavFadesAway: the mobile menu must fade (transition on
// opacity/visibility rather than display toggling) and must close when a menu
// link is tapped, revealing the swapped page underneath.
func TestMobileNavFadesAway(t *testing.T) {
	css, err := os.ReadFile("../../static/style.css")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"opacity:0", "visibility:hidden", "transition:opacity", // fade-capable hidden state
		"opacity:1", "visibility:visible", // open state
	} {
		if !strings.Contains(string(css), want) {
			t.Errorf("style.css missing %q for fading mobile nav", want)
		}
	}
	// display:none on .mobile-nav would kill the fade — the base rule must
	// not use it (the .mobile-nav.open{display:flex} pattern is gone).
	if strings.Contains(string(css), ".mobile-nav.open{display:flex}") {
		t.Error("mobile nav must fade via opacity/visibility, not display toggling")
	}

	s := newTestServer(t)
	body := get(t, s, "/", nil).Body.String()
	if !strings.Contains(body, "mobileNav.classList.remove('open')") {
		t.Error("layout script must close the mobile nav")
	}
	// A tap on any menu link (not just Escape or the toggle) closes the menu.
	if !strings.Contains(body, "mobileNav.querySelectorAll('a')") {
		t.Error("menu links must be wired to close the menu on tap")
	}
}

// TestLogoTapClosesMobileNav: the nav logo stays visible above the open menu;
// tapping it must also fade the menu away while navigating home.
func TestLogoTapClosesMobileNav(t *testing.T) {
	s := newTestServer(t)
	body := get(t, s, "/", nil).Body.String()
	if !strings.Contains(body, "querySelector('.nav-logo')") {
		t.Error("layout script must select the nav logo")
	}
	if !strings.Contains(body, "navLogo.addEventListener('click', closeMenu)") {
		t.Error("logo tap must close the mobile menu")
	}
}
