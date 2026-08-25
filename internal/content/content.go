// Package content loads the site's sections and cards from markdown files.
package content

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// SiteURL is the canonical base URL of the site.
const SiteURL = "https://adamndegwa.com"

// Section represents a navigable section of the portfolio.
type Section struct {
	ID          string // URL slug: "technical-notes", "projects", etc.
	Title       string // Full title: "Technical Notes"
	NavTitle    string // Short nav label: "Notes"
	Subtitle    string // Description for section header
	TileDesc    string // Short description for home page tiles
	Label       string // Section label: "01", "02", "03", "04"
	AccentLabel string // Accent label: "notes", "projects", "musings", "bs"
	Cards       []Card
}

// Card represents a single content card within a section.
type Card struct {
	ID          string // URL slug within section
	Title       string
	Subtitle    string   // Display tags: "Go · Architecture"
	Tags        []string // Parsed tags: ["Go", "Architecture"]
	Description string   // Short description for card listing
	CardIcon    string   // Text icon: "{ }", "AI", etc.
	ReadingTime string   // "3 min"
	Date        string   // ISO date: "2025-06-01"
	Updated     string   // Optional ISO date for last update
	Detail      string   // HTML content for detail view
	DemoURL     string   // Optional external demo/repo URL
}

// CardRef is a reference to a card with its parent section context.
type CardRef struct {
	SectionID    string
	SectionTitle string
	Card         Card
}

// CommaSubtitle returns the subtitle with · replaced by commas for article meta.
func (c Card) CommaSubtitle() string {
	return strings.ReplaceAll(c.Subtitle, " · ", ", ")
}

// FormattedDate returns the date as "January 2, 2006".
func (c Card) FormattedDate() string {
	t, err := time.Parse("2006-01-02", c.Date)
	if err != nil {
		return c.Date
	}
	return t.Format("January 2, 2006")
}

// FormattedUpdated returns the updated date as "January 2, 2006", or empty.
func (c Card) FormattedUpdated() string {
	if c.Updated == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", c.Updated)
	if err != nil {
		return c.Updated
	}
	return t.Format("January 2, 2006")
}

// URL returns the full path for this card within a section.
func (c Card) URL(sectionID string) string {
	return "/" + sectionID + "/" + c.ID
}

// ContentDir is the directory to load markdown content from.
var ContentDir = "content"

// sectionMeta defines the site's fixed sections and their display metadata,
// in nav order. Cards come from markdown files under ContentDir/<ID>/.
var sectionMeta = []Section{
	{
		ID:          "technical-notes",
		Title:       "Technical Notes",
		NavTitle:    "Notes",
		Subtitle:    "Lessons learned, patterns documented and things worth remembering.",
		TileDesc:    "Patterns documented, lessons learned.",
		Label:       "01",
		AccentLabel: "notes",
	},
	{
		ID:          "projects",
		Title:       "Projects",
		NavTitle:    "Projects",
		Subtitle:    "Things I've built, shipped, and maintained.",
		TileDesc:    "Built, shipped, maintained.",
		Label:       "02",
		AccentLabel: "projects",
	},
	{
		ID:          "musings",
		Title:       "Musings",
		NavTitle:    "Musings",
		Subtitle:    "Thoughts on software, design, work and everything in between.",
		TileDesc:    "Software, design, everything between.",
		Label:       "03",
		AccentLabel: "musings",
	},
	{
		ID:          "the-bullshitters",
		Title:       "The Bullshitters",
		NavTitle:    "Bullshitters",
		Subtitle:    "Absurdities — ideas so strange you have to stop and actually think about them.",
		TileDesc:    "Thoughts I should've kept to myself.",
		Label:       "04",
		AccentLabel: "bs",
	},
}

// Sections returns all portfolio sections with cards loaded from markdown.
func Sections() []Section {
	sections := make([]Section, len(sectionMeta))
	copy(sections, sectionMeta)

	if _, err := os.Stat(ContentDir); err != nil {
		return sections
	}

	for i := range sections {
		dir := filepath.Join(ContentDir, sections[i].ID)
		if _, err := os.Stat(dir); err != nil {
			continue
		}
		cards, err := LoadCardsFromDir(dir)
		if err != nil {
			continue
		}
		// Normalise Subtitle to the "Tag1 · Tag2" display format.
		for j := range cards {
			if cards[j].Subtitle == "" || cards[j].Subtitle == strings.Join(cards[j].Tags, ", ") {
				cards[j].Subtitle = strings.Join(cards[j].Tags, " · ")
			}
		}
		sections[i].Cards = cards
	}

	return sections
}

// SectionByID looks up a section by its slug.
func SectionByID(id string) (Section, bool) {
	for _, s := range Sections() {
		if s.ID == id {
			return s, true
		}
	}
	return Section{}, false
}

// CardByID looks up a card within a section.
func CardByID(sectionID, cardID string) (Section, Card, bool) {
	s, ok := SectionByID(sectionID)
	if !ok {
		return Section{}, Card{}, false
	}
	for _, c := range s.Cards {
		if c.ID == cardID {
			return s, c, true
		}
	}
	return s, Card{}, false
}

// AllTags returns all unique tags across all sections, sorted alphabetically.
func AllTags() []string {
	seen := make(map[string]bool)
	for _, s := range Sections() {
		for _, c := range s.Cards {
			for _, t := range c.Tags {
				seen[t] = true
			}
		}
	}
	tags := make([]string, 0, len(seen))
	for t := range seen {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	return tags
}

// TagSlug converts a tag name to a URL-safe slug.
func TagSlug(tag string) string {
	s := strings.ToLower(tag)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, ".", "")
	return s
}

// TagFromSlug finds the original tag name from a URL slug.
func TagFromSlug(slug string) string {
	for _, s := range Sections() {
		for _, c := range s.Cards {
			for _, t := range c.Tags {
				if TagSlug(t) == slug {
					return t
				}
			}
		}
	}
	return slug
}

// CardsByTag returns all cards matching a tag (case-insensitive).
func CardsByTag(tag string) []CardRef {
	var refs []CardRef
	for _, s := range Sections() {
		for _, c := range s.Cards {
			for _, t := range c.Tags {
				if strings.EqualFold(t, tag) {
					refs = append(refs, CardRef{SectionID: s.ID, SectionTitle: s.Title, Card: c})
					break
				}
			}
		}
	}
	return refs
}

// RelatedCards finds cards related to the given card by shared tags.
func RelatedCards(sectionID, cardID string, limit int) []CardRef {
	_, card, ok := CardByID(sectionID, cardID)
	if !ok {
		return nil
	}

	tagSet := make(map[string]bool)
	for _, t := range card.Tags {
		tagSet[t] = true
	}

	type scored struct {
		ref   CardRef
		score int
	}

	var candidates []scored
	for _, s := range Sections() {
		for _, c := range s.Cards {
			if s.ID == sectionID && c.ID == cardID {
				continue
			}
			score := 0
			for _, t := range c.Tags {
				if tagSet[t] {
					score++
				}
			}
			if score > 0 {
				candidates = append(candidates, scored{
					ref:   CardRef{SectionID: s.ID, SectionTitle: s.Title, Card: c},
					score: score,
				})
			}
		}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	result := make([]CardRef, 0, limit)
	for i := 0; i < len(candidates) && i < limit; i++ {
		result = append(result, candidates[i].ref)
	}
	return result
}

// AllCards returns every card across all sections.
func AllCards() []CardRef {
	var refs []CardRef
	for _, s := range Sections() {
		for _, c := range s.Cards {
			refs = append(refs, CardRef{SectionID: s.ID, SectionTitle: s.Title, Card: c})
		}
	}
	return refs
}
