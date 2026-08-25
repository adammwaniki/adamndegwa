package markdown

import (
	"strings"
	"testing"
)

func TestParagraphs(t *testing.T) {
	got := Render("Hello world.\n\nSecond paragraph.")
	want := "<p>Hello world.</p>\n<p>Second paragraph.</p>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHeadings(t *testing.T) {
	got := Render("# One\n\n## Two\n\n### Three\n\n#### Four")
	for _, h := range []string{"<h1>One</h1>", "<h2>Two</h2>", "<h3>Three</h3>", "<h4>Four</h4>"} {
		if !strings.Contains(got, h) {
			t.Errorf("output missing %s:\n%s", h, got)
		}
	}
}

func TestUnorderedList(t *testing.T) {
	got := Render("- first\n- second\n- third")
	want := "<ul>\n<li>first</li>\n<li>second</li>\n<li>third</li>\n</ul>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFencedCodeBlock(t *testing.T) {
	got := Render("```bash\necho hello\nrm -f x\n```")
	want := "<pre><code class=\"language-bash\">echo hello\nrm -f x\n</code></pre>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFencedCodeBlockWithoutLanguage(t *testing.T) {
	got := Render("```\nplain\n```")
	want := "<pre><code>plain\n</code></pre>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInlineCode(t *testing.T) {
	got := Render("Use `net/http` for routing.")
	want := "<p>Use <code>net/http</code> for routing.</p>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBold(t *testing.T) {
	got := Render("This is **important** text.")
	want := "<p>This is <strong>important</strong> text.</p>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLink(t *testing.T) {
	got := Render("See [the repo](https://example.com/repo) for details.")
	want := `<p>See <a href="https://example.com/repo">the repo</a> for details.</p>` + "\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHTMLEscaping(t *testing.T) {
	got := Render("A <script>alert('x')</script> and 1 < 2 & 3 > 2.")
	if strings.Contains(got, "<script>") {
		t.Errorf("raw HTML not escaped: %q", got)
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("expected escaped script tag: %q", got)
	}
}

func TestHTMLEscapingInCodeBlock(t *testing.T) {
	got := Render("```\nif (a < b) { return \"x\"; }\n```")
	if !strings.Contains(got, "a &lt; b") {
		t.Errorf("code block not escaped: %q", got)
	}
}

func TestMixedDocument(t *testing.T) {
	src := "## Heading\n\nIntro with `code`.\n\n- one\n- two\n\n```go\nfmt.Println()\n```\n\nOutro."
	got := Render(src)
	for _, want := range []string{
		"<h2>Heading</h2>",
		"<p>Intro with <code>code</code>.</p>",
		"<ul>\n<li>one</li>\n<li>two</li>\n</ul>",
		"<pre><code class=\"language-go\">fmt.Println()\n</code></pre>",
		"<p>Outro.</p>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q:\n%s", want, got)
		}
	}
}

func TestEmptyInput(t *testing.T) {
	if got := Render(""); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestLineBreaksInsideParagraphAreSpaces(t *testing.T) {
	got := Render("line one\nline two")
	want := "<p>line one\nline two</p>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
