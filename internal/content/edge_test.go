package content

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCardsFromDirUnreadableFile(t *testing.T) {
	dir := t.TempDir()
	writeCard(t, dir, "bad", "---\ntitle: B\ndate: 2025-01-01\n---\n\nb")
	if err := os.Chmod(filepath.Join(dir, "bad.md"), 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(filepath.Join(dir, "bad.md"), 0o644) })

	if _, err := LoadCardsFromDir(dir); err == nil {
		t.Error("expected error for unreadable file")
	}
}

func TestSectionsSkipsSectionWithBrokenCard(t *testing.T) {
	dir := t.TempDir()
	writeCard(t, filepath.Join(dir, "projects"), "bad", "---\ntitle: B\ndate: 2025-01-01\n---\n\nb")
	if err := os.Chmod(filepath.Join(dir, "projects", "bad.md"), 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(filepath.Join(dir, "projects", "bad.md"), 0o644) })
	withContentDir(t, dir)

	sections := Sections()
	if len(sections[1].Cards) != 0 {
		t.Errorf("broken cards should leave the section empty, got %+v", sections[1].Cards)
	}
}

func TestParseFrontmatterEdgeCases(t *testing.T) {
	// Opening --- with no closing delimiter: whole input is body.
	front, body := parseFrontmatter("---\ntitle: dangling\n\nno closing fence")
	if len(front) != 0 {
		t.Errorf("front = %v, want empty", front)
	}
	if body == "" {
		t.Error("body should be the full input")
	}

	// Blank and colon-less lines are skipped; quoted values are unquoted.
	front, body = parseFrontmatter("---\n\ntitle: \"Quoted Title\"\nno-colon-line\ndate: 2025-01-01\n---\n\nbody text")
	if front["title"] != "Quoted Title" || front["date"] != "2025-01-01" || len(front) != 2 {
		t.Errorf("front = %v", front)
	}
	if body != "body text" {
		t.Errorf("body = %q", body)
	}
}

func TestParseTagsEdgeCases(t *testing.T) {
	if got := parseTags(""); got != nil {
		t.Errorf("empty tags = %v, want nil", got)
	}
	got := parseTags("Go · Architecture, , CLI")
	if len(got) != 3 || got[0] != "Go" || got[1] != "Architecture" || got[2] != "CLI" {
		t.Errorf("got %v", got)
	}
}
