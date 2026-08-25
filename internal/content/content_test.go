package content

import (
	"path/filepath"
	"testing"
)

func setupSections(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	writeCard(t, filepath.Join(dir, "technical-notes"), "go-note",
		"---\ntitle: Go Note\ntags: Go, Architecture\ndescription: d\ndate: 2025-01-01\n---\n\nbody")
	writeCard(t, filepath.Join(dir, "projects"), "proj-a",
		"---\ntitle: Proj A\ntags: Go, CLI\ndescription: d\ndate: 2025-02-01\n---\n\nbody")
	writeCard(t, filepath.Join(dir, "musings"), "muse-a",
		"---\ntitle: Muse A\ntags: Design, Philosophy\ndescription: d\ndate: 2025-03-01\n---\n\nbody")
	return dir
}

func withContentDir(t *testing.T, dir string) {
	t.Helper()
	old := ContentDir
	ContentDir = dir
	t.Cleanup(func() { ContentDir = old })
}

func TestSectionsLoadsFourSectionsWithCards(t *testing.T) {
	withContentDir(t, setupSections(t))
	sections := Sections()
	if len(sections) != 4 {
		t.Fatalf("len(Sections()) = %d, want 4", len(sections))
	}
	wantIDs := []string{"technical-notes", "projects", "musings", "the-bullshitters"}
	for i, id := range wantIDs {
		if sections[i].ID != id {
			t.Errorf("sections[%d].ID = %q, want %q", i, sections[i].ID, id)
		}
	}
	if len(sections[0].Cards) != 1 || sections[0].Cards[0].Title != "Go Note" {
		t.Errorf("technical-notes cards = %+v", sections[0].Cards)
	}
	// A section with no markdown dir still appears, with no cards.
	if len(sections[3].Cards) != 0 {
		t.Errorf("the-bullshitters should have 0 cards, got %d", len(sections[3].Cards))
	}
	// Subtitle derived from tags.
	if sections[0].Cards[0].Subtitle != "Go · Architecture" {
		t.Errorf("Subtitle = %q", sections[0].Cards[0].Subtitle)
	}
}

func TestSectionsMissingContentDir(t *testing.T) {
	withContentDir(t, filepath.Join(t.TempDir(), "absent"))
	if len(Sections()) != 4 {
		t.Error("sections metadata should survive a missing content dir")
	}
}

func TestSectionByID(t *testing.T) {
	withContentDir(t, setupSections(t))
	s, ok := SectionByID("projects")
	if !ok || s.Title != "Projects" {
		t.Errorf("SectionByID(projects) = %q, %v", s.Title, ok)
	}
	if _, ok := SectionByID("nope"); ok {
		t.Error("expected !ok for unknown section")
	}
}

func TestCardByID(t *testing.T) {
	withContentDir(t, setupSections(t))
	_, card, ok := CardByID("musings", "muse-a")
	if !ok || card.Title != "Muse A" {
		t.Errorf("CardByID = %q, %v", card.Title, ok)
	}
	if _, _, ok := CardByID("musings", "nope"); ok {
		t.Error("expected !ok for unknown card")
	}
	if _, _, ok := CardByID("nope", "muse-a"); ok {
		t.Error("expected !ok for unknown section")
	}
}

func TestAllTags(t *testing.T) {
	withContentDir(t, setupSections(t))
	tags := AllTags()
	want := []string{"Architecture", "CLI", "Design", "Go", "Philosophy"}
	if len(tags) != len(want) {
		t.Fatalf("AllTags = %v", tags)
	}
	for i, w := range want {
		if tags[i] != w {
			t.Errorf("tags[%d] = %q, want %q", i, tags[i], w)
		}
	}
}

func TestTagSlugAndBack(t *testing.T) {
	cases := map[string]string{
		"Go":           "go",
		"Full-Stack":   "full-stack",
		"C++":          "c++",
		"Node.js":      "nodejs",
		"Public Infra": "public-infra",
	}
	for tag, slug := range cases {
		if got := TagSlug(tag); got != slug {
			t.Errorf("TagSlug(%q) = %q, want %q", tag, got, slug)
		}
	}

	withContentDir(t, setupSections(t))
	if got := TagFromSlug("go"); got != "Go" {
		t.Errorf("TagFromSlug(go) = %q", got)
	}
	if got := TagFromSlug("unknown"); got != "unknown" {
		t.Errorf("TagFromSlug(unknown) = %q, want fallback to slug", got)
	}
}

func TestCardsByTag(t *testing.T) {
	withContentDir(t, setupSections(t))
	refs := CardsByTag("go") // case-insensitive
	if len(refs) != 2 {
		t.Fatalf("len = %d, want 2", len(refs))
	}
	if refs[0].SectionID == "" || refs[0].SectionTitle == "" {
		t.Errorf("refs should carry section context: %+v", refs[0])
	}
	if got := CardsByTag("Nonexistent"); len(got) != 0 {
		t.Errorf("expected no refs, got %v", got)
	}
}

func TestRelatedCards(t *testing.T) {
	dir := t.TempDir()
	writeCard(t, filepath.Join(dir, "technical-notes"), "base",
		"---\ntitle: Base\ntags: Go, Architecture\ndate: 2025-01-01\n---\n\nb")
	writeCard(t, filepath.Join(dir, "technical-notes"), "one-tag",
		"---\ntitle: OneTag\ntags: Go\ndate: 2025-01-02\n---\n\nb")
	writeCard(t, filepath.Join(dir, "technical-notes"), "two-tags",
		"---\ntitle: TwoTags\ntags: Go, Architecture\ndate: 2025-01-03\n---\n\nb")
	writeCard(t, filepath.Join(dir, "projects"), "unrelated",
		"---\ntitle: Unrelated\ntags: Cooking\ndate: 2025-01-04\n---\n\nb")
	withContentDir(t, dir)

	related := RelatedCards("technical-notes", "base", 5)
	if len(related) != 2 {
		t.Fatalf("len = %d, want 2 (unrelated excluded)", len(related))
	}
	if related[0].Card.Title != "TwoTags" {
		t.Errorf("highest-scoring related card should come first, got %q", related[0].Card.Title)
	}
	if got := RelatedCards("technical-notes", "base", 1); len(got) != 1 {
		t.Errorf("limit not respected: %d", len(got))
	}
	if got := RelatedCards("nope", "base", 5); got != nil {
		t.Errorf("expected nil for unknown card, got %v", got)
	}
}

func TestAllCards(t *testing.T) {
	withContentDir(t, setupSections(t))
	if got := AllCards(); len(got) != 3 {
		t.Errorf("len(AllCards()) = %d, want 3", len(got))
	}
}

func TestCardHelpers(t *testing.T) {
	c := Card{ID: "x", Subtitle: "Go · Architecture", Date: "2025-06-01", Updated: "2025-07-04"}
	if got := c.CommaSubtitle(); got != "Go, Architecture" {
		t.Errorf("CommaSubtitle = %q", got)
	}
	if got := c.FormattedDate(); got != "June 1, 2025" {
		t.Errorf("FormattedDate = %q", got)
	}
	if got := c.FormattedUpdated(); got != "July 4, 2025" {
		t.Errorf("FormattedUpdated = %q", got)
	}
	if got := c.URL("technical-notes"); got != "/technical-notes/x" {
		t.Errorf("URL = %q", got)
	}

	bad := Card{Date: "not-a-date", Updated: "also-bad"}
	if got := bad.FormattedDate(); got != "not-a-date" {
		t.Errorf("invalid date should pass through, got %q", got)
	}
	if got := bad.FormattedUpdated(); got != "also-bad" {
		t.Errorf("invalid updated should pass through, got %q", got)
	}
	if got := (Card{}).FormattedUpdated(); got != "" {
		t.Errorf("empty updated should be empty, got %q", got)
	}
}
