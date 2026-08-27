package content

import (
	"os"
	"path/filepath"
	"testing"
)

func writeCard(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

const sampleCard = `---
title: Test Post
tags: Go, Architecture
description: A test post.
icon: "{ }"
reading_time: 3 min
date: 2025-05-01
updated: 2025-06-15
demo_url: https://example.com/demo
---

Body **text** here.

## Heading

- one
- two
`

func TestLoadCardFromMarkdown(t *testing.T) {
	dir := t.TempDir()
	writeCard(t, dir, "test-post", sampleCard)

	card, err := LoadCardFromMarkdown(filepath.Join(dir, "test-post.md"))
	if err != nil {
		t.Fatal(err)
	}
	if card.ID != "test-post" {
		t.Errorf("ID = %q", card.ID)
	}
	if card.Title != "Test Post" {
		t.Errorf("Title = %q", card.Title)
	}
	if len(card.Tags) != 2 || card.Tags[0] != "Go" || card.Tags[1] != "Architecture" {
		t.Errorf("Tags = %v", card.Tags)
	}
	if card.Subtitle != "Go, Architecture" {
		t.Errorf("Subtitle = %q", card.Subtitle)
	}
	if card.Description != "A test post." || card.CardIcon != "{ }" ||
		card.ReadingTime != "3 min" || card.Date != "2025-05-01" ||
		card.Updated != "2025-06-15" || card.DemoURL != "https://example.com/demo" {
		t.Errorf("unexpected card: %+v", card)
	}
	for _, want := range []string{"<strong>text</strong>", "<h2>Heading</h2>", "<li>one</li>"} {
		if !contains(card.Detail, want) {
			t.Errorf("Detail missing %q:\n%s", want, card.Detail)
		}
	}
}

func TestLoadCardFromMarkdownMissingFile(t *testing.T) {
	if _, err := LoadCardFromMarkdown(filepath.Join(t.TempDir(), "nope.md")); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadCardFromMarkdownNoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	writeCard(t, dir, "bare", "Just a body.")
	card, err := LoadCardFromMarkdown(filepath.Join(dir, "bare.md"))
	if err != nil {
		t.Fatal(err)
	}
	if card.Title != "" || card.Detail != "<p>Just a body.</p>\n" {
		t.Errorf("unexpected card: %+v", card)
	}
}

func TestLoadCardsFromDirSortsByDateDescending(t *testing.T) {
	dir := t.TempDir()
	writeCard(t, dir, "newer", "---\ntitle: Newer\ndate: 2025-06-01\n---\n\nn")
	writeCard(t, dir, "older", "---\ntitle: Older\ndate: 2025-01-01\n---\n\no")
	writeCard(t, dir, "middle", "---\ntitle: Middle\ndate: 2025-03-01\n---\n\nm")

	cards, err := LoadCardsFromDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 3 || cards[0].ID != "newer" || cards[1].ID != "middle" || cards[2].ID != "older" {
		t.Errorf("order = %v", []string{cards[0].ID, cards[1].ID, cards[2].ID})
	}
}

func TestLoadCardsFromDirAutoNumbersIconlessCards(t *testing.T) {
	dir := t.TempDir()
	writeCard(t, dir, "newer", "---\ntitle: Newer\ndate: 2025-06-01\n---\n\nn")
	writeCard(t, dir, "older", "---\ntitle: Older\ndate: 2025-01-01\nicon: \"#\"\n---\n\no")

	cards, err := LoadCardsFromDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cards[0].CardIcon != "01" || cards[1].CardIcon != "#" {
		t.Errorf("icons = %q, %q; want positional 01 then explicit #", cards[0].CardIcon, cards[1].CardIcon)
	}
}

func TestLoadCardsFromDirMissingDir(t *testing.T) {
	if _, err := LoadCardsFromDir(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Error("expected error for missing dir")
	}
}

func TestLoadCardsFromDirSkipsNonMarkdown(t *testing.T) {
	dir := t.TempDir()
	writeCard(t, dir, "real", "---\ntitle: R\ndate: 2025-01-01\n---\n\nx")
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	cards, err := LoadCardsFromDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 1 {
		t.Errorf("len = %d, want 1", len(cards))
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}
