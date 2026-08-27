package content

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"adamndegwa/internal/markdown"
)

// LoadCardFromMarkdown reads a markdown file and returns a Card.
// Frontmatter is parsed as simple "key: value" pairs between --- delimiters.
func LoadCardFromMarkdown(path string) (Card, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Card{}, err
	}

	front, body := parseFrontmatter(string(data))

	return Card{
		ID:          strings.TrimSuffix(filepath.Base(path), ".md"),
		Title:       front["title"],
		Subtitle:    front["tags"],
		Tags:        parseTags(front["tags"]),
		Description: front["description"],
		CardIcon:    front["icon"],
		ReadingTime: front["reading_time"],
		Date:        front["date"],
		Updated:     front["updated"],
		Detail:      markdown.Render(body),
		DemoURL:     front["demo_url"],
	}, nil
}

// LoadCardsFromDir reads all .md files from a directory, sorted by date
// descending — most recent first (cards without dates sort last, by filename).
func LoadCardsFromDir(dir string) ([]Card, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var cards []Card
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		card, err := LoadCardFromMarkdown(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	sort.SliceStable(cards, func(i, j int) bool {
		return cards[i].Date > cards[j].Date
	})
	// Cards without an icon in frontmatter are numbered by position in the
	// sorted list (01 = most recent), so numbering stays correct as cards
	// are added or removed.
	for i := range cards {
		if cards[i].CardIcon == "" {
			cards[i].CardIcon = fmt.Sprintf("%02d", i+1)
		}
	}
	return cards, nil
}

// parseFrontmatter splits --- delimited frontmatter from the body.
func parseFrontmatter(content string) (map[string]string, string) {
	front := make(map[string]string)

	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "---") {
		return front, content
	}

	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return front, content
	}

	frontBlock := rest[:idx]
	body := strings.TrimSpace(rest[idx+4:])

	for _, line := range strings.Split(frontBlock, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		colon := strings.Index(line, ":")
		if colon < 0 {
			continue
		}
		key := strings.TrimSpace(line[:colon])
		value := strings.TrimSpace(line[colon+1:])
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}
		front[key] = value
	}

	return front, body
}

// parseTags splits a comma-or-dot-separated tag string into a slice.
func parseTags(s string) []string {
	s = strings.ReplaceAll(s, " · ", ", ")
	parts := strings.Split(s, ",")
	var tags []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			tags = append(tags, p)
		}
	}
	return tags
}
