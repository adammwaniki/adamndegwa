package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"adamndegwa/internal/content"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	content.ContentDir = "../../content"
	t.Cleanup(func() { content.ContentDir = "content" })
	s, err := New("../../views", "../../static")
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func get(t *testing.T, s *Server, path string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	return rec
}

func TestHomeFullPage(t *testing.T) {
	s := newTestServer(t)
	rec := get(t, s, "/", nil)
	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		"<!DOCTYPE html>", `<html lang="en">`,
		"<title>Adam Ndegwa | Software Engineer</title>",
		`<link rel="canonical" href="https://adamndegwa.com/">`,
		`property="og:title"`, `name="twitter:card"`,
		`application/ld+json`, `"@type":"WebSite"`,
		"Adam", "Ndegwa", // hero
		"Technical Notes", "Projects", "Musings", "The Bullshitters", // tiles
		`href="/static/style.css"`,
		"/static/htmx.min.js",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("home page missing %q", want)
		}
	}
}

func TestHomeHTMXFragment(t *testing.T) {
	s := newTestServer(t)
	rec := get(t, s, "/", map[string]string{"HX-Request": "true"})
	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("HTMX request should receive a fragment, not the full layout")
	}
	if rec.Header().Get("HX-Title") != "Adam Ndegwa | Software Engineer" {
		t.Errorf("HX-Title = %q", rec.Header().Get("HX-Title"))
	}
	if !strings.Contains(body, "Adam") {
		t.Error("fragment should still contain page content")
	}
}

func TestHTMXHistoryRestoreGetsFullPage(t *testing.T) {
	s := newTestServer(t)
	rec := get(t, s, "/", map[string]string{
		"HX-Request":                "true",
		"HX-History-Restore-Request": "true",
	})
	if !strings.Contains(rec.Body.String(), "<!DOCTYPE html>") {
		t.Error("history restore should receive the full page")
	}
}

func TestSectionPage(t *testing.T) {
	s := newTestServer(t)
	rec := get(t, s, "/technical-notes", nil)
	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "<title>Technical Notes | Adam Ndegwa</title>") {
		t.Error("wrong title")
	}
	if !strings.Contains(body, "Building with Go") {
		t.Error("section should list its cards")
	}
	if !strings.Contains(body, `href="/technical-notes/go-std-lib"`) {
		t.Error("cards should link to detail pages")
	}
}

func TestUnknownSectionIs404(t *testing.T) {
	s := newTestServer(t)
	if rec := get(t, s, "/no-such-section", nil); rec.Code != 404 {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestDetailPage(t *testing.T) {
	s := newTestServer(t)
	// htmx-patterns is a middle card in technical-notes, so it has prev and next.
	rec := get(t, s, "/technical-notes/htmx-patterns", nil)
	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		"<title>HTMX Patterns for Server-Driven UIs | Adam Ndegwa</title>",
		`property="og:type" content="article"`,
		`"@type":"Article"`,
		`rel="prev"`, `rel="next"`,
		"Related",
		`href="/technical-notes"`, // back link
	} {
		if !strings.Contains(body, want) {
			t.Errorf("detail page missing %q", want)
		}
	}
}

func TestDetailHTMXFragment(t *testing.T) {
	s := newTestServer(t)
	rec := get(t, s, "/technical-notes/go-std-lib", map[string]string{"HX-Request": "true"})
	if strings.Contains(rec.Body.String(), "<!DOCTYPE html>") {
		t.Error("expected fragment")
	}
	if rec.Header().Get("HX-Title") == "" {
		t.Error("expected HX-Title header")
	}
}

func TestDetailEmptyCardRedirects(t *testing.T) {
	s := newTestServer(t)
	rec := get(t, s, "/technical-notes/", nil)
	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/technical-notes" {
		t.Errorf("Location = %q", loc)
	}
}

func TestUnknownCardIs404(t *testing.T) {
	s := newTestServer(t)
	if rec := get(t, s, "/technical-notes/no-such-card", nil); rec.Code != 404 {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestFirstCardHasNoPrev(t *testing.T) {
	s := newTestServer(t)
	sections := content.Sections()
	first := sections[0].Cards[0]
	rec := get(t, s, first.URL(sections[0].ID), nil)
	body := rec.Body.String()
	if !strings.Contains(body, `rel="next"`) {
		t.Error("first card should have a next link")
	}
	if strings.Contains(body, `rel="prev"`) {
		t.Error("first card should not have a prev link")
	}
}

func TestTagsPages(t *testing.T) {
	s := newTestServer(t)
	if rec := get(t, s, "/tags", nil); rec.Code != 200 || !strings.Contains(rec.Body.String(), "<title>Tags | Adam Ndegwa</title>") {
		t.Errorf("/tags status = %d", rec.Code)
	}
	rec := get(t, s, "/tags/go", nil)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Building with Go&#39;s Standard Library") && !strings.Contains(rec.Body.String(), "Building with Go's Standard Library") {
		t.Errorf("/tags/go should list Go articles, status %d", rec.Code)
	}
	if rec := get(t, s, "/tags/", nil); rec.Code != http.StatusFound || rec.Header().Get("Location") != "/tags" {
		t.Errorf("/tags/ should redirect to /tags, got %d", rec.Code)
	}
	if rec := get(t, s, "/tags/no-such-tag", nil); rec.Code != 404 {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestContactPage(t *testing.T) {
	s := newTestServer(t)
	rec := get(t, s, "/contact", nil)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "<title>Contact | Adam Ndegwa</title>") {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestSearchIndexEndpoint(t *testing.T) {
	s := newTestServer(t)
	rec := get(t, s, "/search-index.json", nil)
	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	var entries []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &entries); err != nil || len(entries) == 0 {
		t.Errorf("invalid or empty search index: %v", err)
	}
}

func TestSEOEndpoints(t *testing.T) {
	s := newTestServer(t)
	for path, ct := range map[string]string{
		"/sitemap.xml": "application/xml",
		"/robots.txt":  "text/plain",
		"/feed.xml":    "application/rss+xml",
	} {
		rec := get(t, s, path, nil)
		if rec.Code != 200 {
			t.Errorf("%s status = %d", path, rec.Code)
		}
		if got := rec.Header().Get("Content-Type"); !strings.HasPrefix(got, ct) {
			t.Errorf("%s Content-Type = %q, want prefix %q", path, got, ct)
		}
	}
	if !strings.Contains(get(t, s, "/sitemap.xml", nil).Body.String(), "<urlset") {
		t.Error("sitemap should contain urlset")
	}
}

func TestStaticFiles(t *testing.T) {
	s := newTestServer(t)
	rec := get(t, s, "/static/style.css", nil)
	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "--ink:") {
		t.Error("style.css should define brand tokens")
	}
}

func TestNewWithMissingViewsDir(t *testing.T) {
	if _, err := New(t.TempDir(), t.TempDir()); err == nil {
		t.Error("expected error when no templates exist")
	}
}

func TestRunUsesInjectedListener(t *testing.T) {
	content.ContentDir = "../../content"
	t.Cleanup(func() { content.ContentDir = "content" })
	called := false
	err := Run(":0", "../../views", "../../static", func(addr string, h http.Handler) error {
		called = true
		if addr != ":0" {
			t.Errorf("addr = %q", addr)
		}
		return nil
	})
	if err != nil || !called {
		t.Errorf("Run err = %v, called = %v", err, called)
	}

	sentinel := http.ErrServerClosed
	if err := Run(":0", "../../views", "../../static", func(string, http.Handler) error { return sentinel }); err != sentinel {
		t.Errorf("listener error should propagate, got %v", err)
	}
}
