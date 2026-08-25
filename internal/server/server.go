// Package server wires HTTP routing, template rendering and HTMX partial
// responses for the site.
package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"adamndegwa/internal/content"
	"adamndegwa/internal/seo"
)

// Server is the site's HTTP handler.
type Server struct {
	mux *http.ServeMux
	tpl *template.Template
}

// New parses the views under viewsDir and registers all routes, serving
// static assets from staticDir.
func New(viewsDir, staticDir string) (*Server, error) {
	tpl, err := template.New("").Funcs(template.FuncMap{
		"raw":     func(s string) template.HTML { return template.HTML(s) },
		"dict":    dict,
		"tagSlug": content.TagSlug,
		"join":    strings.Join,
		"lower":   strings.ToLower,
	}).ParseGlob(filepath.Join(viewsDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("parsing views: %w", err)
	}
	if tpl, err = tpl.ParseGlob(filepath.Join(viewsDir, "partials", "*.html")); err != nil {
		return nil, fmt.Errorf("parsing partials: %w", err)
	}

	s := &Server{mux: http.NewServeMux(), tpl: tpl}

	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	s.mux.HandleFunc("/", s.handleHome)
	s.mux.HandleFunc("/tags", s.handleTags)
	s.mux.HandleFunc("/tags/", s.handleTagsBySlug)
	s.mux.HandleFunc("/contact", s.handleContact)
	s.mux.HandleFunc("/search-index.json", s.handleSearchIndex)
	s.mux.HandleFunc("/sitemap.xml", s.handleSitemap)
	s.mux.HandleFunc("/robots.txt", s.handleRobots)
	s.mux.HandleFunc("/feed.xml", s.handleFeed)

	for _, sec := range content.Sections() {
		s.mux.HandleFunc("/"+sec.ID, s.handleSection)
		s.mux.HandleFunc("/"+sec.ID+"/", s.handleDetail)
	}

	return s, nil
}

// Run builds the server and serves it with the given listen function
// (http.ListenAndServe in production).
func Run(addr, viewsDir, staticDir string, listen func(string, http.Handler) error) error {
	s, err := New(viewsDir, staticDir)
	if err != nil {
		return err
	}
	return listen(addr, s)
}

// RunFromEnv serves on the PORT environment variable (default 8080).
func RunFromEnv() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Run(":"+port, "views", "static", listenAndServe)
}

// listenAndServe is swappable in tests.
var listenAndServe = http.ListenAndServe

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func dict(values ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i < len(values)-1; i += 2 {
		key, _ := values[i].(string)
		m[key] = values[i+1]
	}
	return m
}

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// renderPage answers HTMX requests with the page-content fragment (plus an
// HX-Title header) and everything else with the full layout.
func (s *Server) renderPage(w http.ResponseWriter, r *http.Request, data map[string]any) {
	if isHTMX(r) && r.Header.Get("HX-History-Restore-Request") != "true" {
		if title, ok := data["Title"].(string); ok {
			w.Header().Set("HX-Title", title)
		}
		if err := s.tpl.ExecuteTemplate(w, "page-content", data); err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	if err := s.tpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func baseData(page, title, description, canonical, ogType string) map[string]any {
	return map[string]any{
		"Page":         page,
		"Title":        title,
		"Description":  description,
		"CanonicalURL": canonical,
		"OGType":       ogType,
		"Sections":     content.Sections(),
		"IsDetail":     false,
	}
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	s.renderPage(w, r, baseData("home",
		"Adam Ndegwa | Software Engineer",
		"Adam Ndegwa | Software Engineer. I build systems that are simple to use, easy to maintain, and quick to scale.",
		content.SiteURL+"/", "website"))
}

func (s *Server) handleSection(w http.ResponseWriter, r *http.Request) {
	sectionID := strings.TrimPrefix(r.URL.Path, "/")
	section, ok := content.SectionByID(sectionID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	data := baseData(sectionID,
		fmt.Sprintf("%s | Adam Ndegwa", section.Title),
		section.Subtitle,
		content.SiteURL+"/"+sectionID, "website")
	data["Section"] = section
	s.renderPage(w, r, data)
}

func (s *Server) handleDetail(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)
	if len(parts) != 2 || parts[1] == "" {
		http.Redirect(w, r, "/"+parts[0], http.StatusFound)
		return
	}

	sectionID, cardID := parts[0], parts[1]
	section, card, ok := content.CardByID(sectionID, cardID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	cardIndex := 0
	var nextCard, prevCard map[string]any
	for i, c := range section.Cards {
		if c.ID == cardID {
			cardIndex = i
			if i > 0 {
				prevCard = map[string]any{
					"Title": section.Cards[i-1].Title,
					"URL":   "/" + sectionID + "/" + section.Cards[i-1].ID,
				}
			}
			if i < len(section.Cards)-1 {
				nextCard = map[string]any{
					"Title": section.Cards[i+1].Title,
					"URL":   "/" + sectionID + "/" + section.Cards[i+1].ID,
				}
			}
			break
		}
	}

	data := baseData(sectionID,
		fmt.Sprintf("%s | Adam Ndegwa", card.Title),
		card.Description,
		content.SiteURL+"/"+sectionID+"/"+cardID, "article")
	data["Section"] = section
	data["Card"] = card
	data["CardIndex"] = cardIndex + 1
	data["NextCard"] = nextCard
	data["PrevCard"] = prevCard
	data["RelatedCards"] = content.RelatedCards(sectionID, cardID, 3)
	data["IsDetail"] = true
	data["BackURL"] = "/" + sectionID
	s.renderPage(w, r, data)
}

func (s *Server) handleTags(w http.ResponseWriter, r *http.Request) {
	data := baseData("tags",
		"Tags | Adam Ndegwa",
		"Browse articles by topic across all sections.",
		content.SiteURL+"/tags", "website")
	data["AllTags"] = content.AllTags()
	s.renderPage(w, r, data)
}

func (s *Server) handleTagsBySlug(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/tags/")
	if slug == "" {
		http.Redirect(w, r, "/tags", http.StatusFound)
		return
	}

	tag := content.TagFromSlug(slug)
	cards := content.CardsByTag(tag)
	if len(cards) == 0 {
		http.NotFound(w, r)
		return
	}

	data := baseData("tags",
		fmt.Sprintf("%s | Adam Ndegwa", tag),
		fmt.Sprintf("Articles tagged with %q on Adam Ndegwa's portfolio.", tag),
		content.SiteURL+"/tags/"+slug, "website")
	data["Tag"] = tag
	data["TagCards"] = cards
	s.renderPage(w, r, data)
}

func (s *Server) handleContact(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, r, baseData("contact",
		"Contact | Adam Ndegwa",
		"Get in touch with Adam Ndegwa | Software Engineer.",
		content.SiteURL+"/contact", "website"))
}

func (s *Server) handleSearchIndex(w http.ResponseWriter, r *http.Request) {
	// SearchIndex only fails if json.Marshal fails, which is impossible for
	// the string-only structure it builds.
	data, _ := seo.SearchIndex(content.Sections())
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *Server) handleSitemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(seo.Sitemap(content.Sections(), time.Now())))
}

func (s *Server) handleRobots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(seo.Robots()))
}

func (s *Server) handleFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Write([]byte(seo.RSS(content.Sections(), time.Now())))
}
