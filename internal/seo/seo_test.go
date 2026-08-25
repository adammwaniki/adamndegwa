package seo

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"strings"
	"testing"
	"time"

	"adamndegwa/internal/content"
)

func TestMain(m *testing.M) {
	content.ContentDir = "../../content"
	os.Exit(m.Run())
}

func TestSitemap(t *testing.T) {
	out := Sitemap(content.Sections(), time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC))

	var parsed struct {
		XMLName xml.Name `xml:"urlset"`
		URLs    []struct {
			Loc      string  `xml:"loc"`
			LastMod  string  `xml:"lastmod"`
			Priority float64 `xml:"priority"`
		} `xml:"url"`
	}
	if err := xml.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("sitemap is not valid XML: %v", err)
	}
	if len(parsed.URLs) == 0 {
		t.Fatal("sitemap has no URLs")
	}
	if parsed.URLs[0].Loc != content.SiteURL+"/" {
		t.Errorf("first URL should be home, got %q", parsed.URLs[0].Loc)
	}

	locSet := make(map[string]bool)
	for _, u := range parsed.URLs {
		locSet[u.Loc] = true
		if u.LastMod == "" {
			t.Errorf("%s missing lastmod", u.Loc)
		}
	}
	// Every section and card URL must appear.
	for _, s := range content.Sections() {
		if !locSet[content.SiteURL+"/"+s.ID] {
			t.Errorf("sitemap missing section %s", s.ID)
		}
		for _, c := range s.Cards {
			if !locSet[content.SiteURL+"/"+s.ID+"/"+c.ID] {
				t.Errorf("sitemap missing card %s/%s", s.ID, c.ID)
			}
		}
	}
	// Tags and contact pages.
	for _, path := range []string{"/tags", "/contact"} {
		if !locSet[content.SiteURL+path] {
			t.Errorf("sitemap missing %s", path)
		}
	}
}

func TestSitemapUsesUpdatedAsLastmod(t *testing.T) {
	sections := []content.Section{{
		ID:    "s",
		Title: "S",
		Cards: []content.Card{
			{ID: "with-update", Title: "T", Date: "2025-01-01", Updated: "2025-05-05"},
			{ID: "no-update", Title: "T2", Date: "2025-02-02"},
		},
	}}
	out := Sitemap(sections, time.Now())
	if !strings.Contains(out, "/s/with-update</loc><lastmod>2025-05-05") {
		t.Error("updated date should be used as lastmod")
	}
	if !strings.Contains(out, "/s/no-update</loc><lastmod>2025-02-02") {
		t.Error("date should be used as lastmod when no updated")
	}
}

func TestRobots(t *testing.T) {
	out := Robots()
	for _, want := range []string{"User-agent: *", "Allow: /", "Sitemap: " + content.SiteURL + "/sitemap.xml"} {
		if !strings.Contains(out, want) {
			t.Errorf("robots.txt missing %q:\n%s", want, out)
		}
	}
}

func TestRSS(t *testing.T) {
	out := RSS(content.Sections(), time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC))

	var parsed struct {
		XMLName xml.Name `xml:"rss"`
		Channel struct {
			Title string `xml:"title"`
			Link  string `xml:"link"`
			Items []struct {
				Title   string `xml:"title"`
				Link    string `xml:"link"`
				PubDate string `xml:"pubDate"`
			} `xml:"item"`
		} `xml:"channel"`
	}
	if err := xml.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("feed is not valid XML: %v", err)
	}
	if parsed.Channel.Title != "Adam Ndegwa" {
		t.Errorf("channel title = %q", parsed.Channel.Title)
	}
	totalCards := 0
	for _, s := range content.Sections() {
		totalCards += len(s.Cards)
	}
	if len(parsed.Channel.Items) != totalCards {
		t.Errorf("items = %d, want %d", len(parsed.Channel.Items), totalCards)
	}
	// Newest first.
	if len(parsed.Channel.Items) >= 2 {
		first, err1 := time.Parse(time.RFC1123Z, parsed.Channel.Items[0].PubDate)
		second, err2 := time.Parse(time.RFC1123Z, parsed.Channel.Items[1].PubDate)
		if err1 != nil || err2 != nil {
			t.Fatalf("unparseable pubDates: %q, %q", parsed.Channel.Items[0].PubDate, parsed.Channel.Items[1].PubDate)
		}
		if first.Before(second) {
			t.Errorf("items not sorted newest-first: %q then %q",
				parsed.Channel.Items[0].PubDate, parsed.Channel.Items[1].PubDate)
		}
	}
}

func TestRSSSortingAndEscaping(t *testing.T) {
	sections := []content.Section{{
		ID:    "s",
		Title: "S & Co",
		Cards: []content.Card{
			{ID: "old", Title: "Old <Post>", Date: "2025-01-01", Description: "old"},
			{ID: "new", Title: "New & Shiny", Date: "2025-06-01", Description: "new"},
			{ID: "bad-date", Title: "Bad", Date: "n/a", Description: "bad"},
		},
	}}
	out := RSS(sections, time.Now())
	if strings.Contains(out, "New & Shiny") || strings.Contains(out, "Old <Post>") {
		t.Error("RSS output must XML-escape titles")
	}
	if !strings.Contains(out, "New &amp; Shiny") {
		t.Error("expected escaped title")
	}
	newIdx := strings.Index(out, "/s/new")
	oldIdx := strings.Index(out, "/s/old")
	if newIdx > oldIdx {
		t.Error("newest item should appear first")
	}
	if !strings.Contains(out, "<pubDate>n/a</pubDate>") {
		t.Error("unparseable date should pass through")
	}
}

func TestSearchIndex(t *testing.T) {
	data, err := SearchIndex(content.Sections())
	if err != nil {
		t.Fatal(err)
	}
	var entries []map[string]any
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("search index is not valid JSON: %v", err)
	}
	totalCards := 0
	for _, s := range content.Sections() {
		totalCards += len(s.Cards)
	}
	if len(entries) != totalCards {
		t.Errorf("entries = %d, want %d", len(entries), totalCards)
	}
	for _, e := range entries {
		for _, key := range []string{"title", "description", "tags", "url", "section"} {
			if _, ok := e[key]; !ok {
				t.Errorf("entry missing key %q: %v", key, e)
			}
		}
	}
}
