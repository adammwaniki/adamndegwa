// Package markdown renders a small, deliberate subset of Markdown to HTML
// using only the standard library. Supported: #-#### headings, paragraphs,
// unordered lists, fenced code blocks, inline code, bold and links.
// All input text is HTML-escaped; the output is safe to inject as template.HTML.
package markdown

import (
	"html"
	"strings"
)

// Render converts markdown source to an HTML fragment.
func Render(src string) string {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	lines := strings.Split(src, "\n")

	var out strings.Builder
	var para []string
	inList := false
	inCode := false
	codeLang := ""

	flushPara := func() {
		if len(para) == 0 {
			return
		}
		out.WriteString("<p>")
		out.WriteString(inline(strings.Join(para, "\n")))
		out.WriteString("</p>\n")
		para = nil
	}
	closeList := func() {
		if inList {
			out.WriteString("</ul>\n")
			inList = false
		}
	}

	for _, line := range lines {
		// Fenced code block state machine.
		if inCode {
			if strings.HasPrefix(strings.TrimSpace(line), "```") {
				out.WriteString("</code></pre>\n")
				inCode = false
				continue
			}
			out.WriteString(html.EscapeString(line))
			out.WriteString("\n")
			continue
		}
		if trimmed := strings.TrimSpace(line); strings.HasPrefix(trimmed, "```") {
			flushPara()
			closeList()
			codeLang = strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
			if codeLang != "" {
				out.WriteString(`<pre><code class="language-` + html.EscapeString(codeLang) + `">`)
			} else {
				out.WriteString("<pre><code>")
			}
			inCode = true
			continue
		}

		if strings.TrimSpace(line) == "" {
			flushPara()
			closeList()
			continue
		}

		// Headings.
		if level, text, ok := heading(line); ok {
			flushPara()
			closeList()
			out.WriteString("<h")
			out.WriteString(string(rune('0' + level)))
			out.WriteString(">")
			out.WriteString(inline(text))
			out.WriteString("</h")
			out.WriteString(string(rune('0' + level)))
			out.WriteString(">\n")
			continue
		}

		// List items.
		if strings.HasPrefix(line, "- ") {
			flushPara()
			if !inList {
				out.WriteString("<ul>\n")
				inList = true
			}
			out.WriteString("<li>")
			out.WriteString(inline(strings.TrimPrefix(line, "- ")))
			out.WriteString("</li>\n")
			continue
		}

		para = append(para, line)
	}
	if inCode {
		out.WriteString("</code></pre>\n")
	}
	flushPara()
	closeList()
	return out.String()
}

// heading parses "# text" through "#### text".
func heading(line string) (level int, text string, ok bool) {
	for level = 0; level < len(line) && line[level] == '#'; level++ {
	}
	if level < 1 || level > 4 || level >= len(line) || line[level] != ' ' {
		return 0, "", false
	}
	return level, line[level+1:], true
}

// inline applies inline formatting to already-plain text: escape HTML first,
// then inline code, bold and links.
func inline(s string) string {
	s = html.EscapeString(s)
	s = replaceSpans(s, "`", func(inner string) string { return "<code>" + inner + "</code>" })
	s = replaceSpans(s, "**", func(inner string) string { return "<strong>" + inner + "</strong>" })
	s = replaceLinks(s)
	return s
}

// replaceSpans wraps text between pairs of delim using wrap.
func replaceSpans(s, delim string, wrap func(string) string) string {
	var out strings.Builder
	for {
		i := strings.Index(s, delim)
		if i < 0 {
			out.WriteString(s)
			break
		}
		j := strings.Index(s[i+len(delim):], delim)
		if j < 0 {
			out.WriteString(s)
			break
		}
		out.WriteString(s[:i])
		out.WriteString(wrap(s[i+len(delim) : i+len(delim)+j]))
		s = s[i+len(delim)+j+len(delim):]
	}
	return out.String()
}

// replaceLinks converts [text](url) (post-escaping) into anchors.
func replaceLinks(s string) string {
	var out strings.Builder
	for {
		open := strings.Index(s, "[")
		if open < 0 {
			out.WriteString(s)
			break
		}
		close := strings.Index(s[open:], "](")
		if close < 0 {
			out.WriteString(s)
			break
		}
		close += open
		end := strings.Index(s[close+2:], ")")
		if end < 0 {
			out.WriteString(s)
			break
		}
		end += close + 2
		text, url := s[open+1:close], s[close+2:end]
		out.WriteString(s[:open])
		out.WriteString(`<a href="` + url + `">` + text + "</a>")
		s = s[end+1:]
	}
	return out.String()
}
