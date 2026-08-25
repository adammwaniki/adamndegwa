// Package seo builds the site's search-engine artefacts: sitemap, robots.txt,
// RSS feed and the client-side search index.
package seo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"adamndegwa/internal/content"
)

// Sitemap returns a sitemap.xml document covering home, all sections and
// cards, plus the tags and contact pages.
func Sitemap(sections []content.Section, now time.Time) string {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	buf.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	today := now.Format("2006-01-02")

	buf.WriteString(fmt.Sprintf("  <url><loc>%s/</loc><lastmod>%s</lastmod><changefreq>weekly</changefreq><priority>1.0</priority></url>\n", content.SiteURL, today))

	for _, s := range sections {
		buf.WriteString(fmt.Sprintf("  <url><loc>%s/%s</loc><lastmod>%s</lastmod><changefreq>weekly</changefreq><priority>0.8</priority></url>\n", content.SiteURL, s.ID, today))
		for _, c := range s.Cards {
			lastmod := c.Date
			if c.Updated != "" {
				lastmod = c.Updated
			}
			buf.WriteString(fmt.Sprintf("  <url><loc>%s/%s/%s</loc><lastmod>%s</lastmod><changefreq>monthly</changefreq><priority>0.6</priority></url>\n", content.SiteURL, s.ID, c.ID, lastmod))
		}
	}

	buf.WriteString(fmt.Sprintf("  <url><loc>%s/tags</loc><lastmod>%s</lastmod><changefreq>weekly</changefreq><priority>0.5</priority></url>\n", content.SiteURL, today))
	buf.WriteString(fmt.Sprintf("  <url><loc>%s/contact</loc><lastmod>%s</lastmod><changefreq>monthly</changefreq><priority>0.5</priority></url>\n", content.SiteURL, today))

	buf.WriteString("</urlset>\n")
	return buf.String()
}

// Robots returns the robots.txt body.
func Robots() string {
	return fmt.Sprintf("User-agent: *\nAllow: /\n\nSitemap: %s/sitemap.xml\n", content.SiteURL)
}

// RSS returns an RSS 2.0 feed of all cards, newest first.
func RSS(sections []content.Section, now time.Time) string {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	buf.WriteString(`<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">` + "\n")
	buf.WriteString("<channel>\n")
	buf.WriteString("  <title>Adam Ndegwa</title>\n")
	buf.WriteString(fmt.Sprintf("  <link>%s</link>\n", content.SiteURL))
	buf.WriteString("  <description>Software Engineer. I build systems that are simple to use, easy to maintain, and quick to scale.</description>\n")
	buf.WriteString("  <language>en</language>\n")
	buf.WriteString(fmt.Sprintf("  <atom:link href=\"%s/feed.xml\" rel=\"self\" type=\"application/rss+xml\"/>\n", content.SiteURL))
	buf.WriteString(fmt.Sprintf("  <lastBuildDate>%s</lastBuildDate>\n", now.Format(time.RFC1123Z)))

	type entry struct {
		section content.Section
		card    content.Card
	}
	var entries []entry
	for _, s := range sections {
		for _, c := range s.Cards {
			entries = append(entries, entry{s, c})
		}
	}
	// Newest first (stable insertion sort keeps order deterministic).
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].card.Date > entries[i].card.Date {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	for _, e := range entries {
		pubDate := e.card.Date
		if t, err := time.Parse("2006-01-02", pubDate); err == nil {
			pubDate = t.Format(time.RFC1123Z)
		}
		link := fmt.Sprintf("%s/%s/%s", content.SiteURL, e.section.ID, e.card.ID)
		buf.WriteString("  <item>\n")
		buf.WriteString(fmt.Sprintf("    <title>%s</title>\n", escapeXML(e.card.Title)))
		buf.WriteString(fmt.Sprintf("    <link>%s</link>\n", link))
		buf.WriteString(fmt.Sprintf("    <guid>%s</guid>\n", link))
		buf.WriteString(fmt.Sprintf("    <pubDate>%s</pubDate>\n", pubDate))
		buf.WriteString(fmt.Sprintf("    <description>%s</description>\n", escapeXML(e.card.Description)))
		buf.WriteString(fmt.Sprintf("    <category>%s</category>\n", escapeXML(e.section.Title)))
		buf.WriteString("  </item>\n")
	}

	buf.WriteString("</channel>\n</rss>\n")
	return buf.String()
}

// SearchIndex returns the JSON search index consumed by client-side search.
func SearchIndex(sections []content.Section) ([]byte, error) {
	var index []map[string]any
	for _, s := range sections {
		for _, c := range s.Cards {
			index = append(index, map[string]any{
				"title":       c.Title,
				"description": c.Description,
				"tags":        c.Tags,
				"url":         "/" + s.ID + "/" + c.ID,
				"section":     s.Title,
			})
		}
	}
	return json.Marshal(index)
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
