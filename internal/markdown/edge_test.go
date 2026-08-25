package markdown

import (
	"strings"
	"testing"
)

func TestUnclosedCodeFenceIsClosedAtEOF(t *testing.T) {
	got := Render("```\nopen block")
	want := "<pre><code>open block\n</code></pre>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestUnmatchedInlineDelimitersAreLeftAlone(t *testing.T) {
	got := Render("a `lonely tick and **lonely stars")
	if strings.Contains(got, "<code>") || strings.Contains(got, "<strong>") {
		t.Errorf("unmatched delimiters should not produce tags: %q", got)
	}
}

func TestUnmatchedLinkSyntaxIsLeftAlone(t *testing.T) {
	for _, src := range []string{"a [lonely bracket", "a [text](no-close", "a [text] (paren gap)"} {
		got := Render(src)
		if strings.Contains(got, "<a href") {
			t.Errorf("%q should not produce a link: %q", src, got)
		}
	}
}

func TestHeadingWithoutSpaceIsParagraph(t *testing.T) {
	got := Render("#hashtag")
	if strings.Contains(got, "<h1>") {
		t.Errorf("#hashtag should not be a heading: %q", got)
	}
}

func TestCRLFNormalised(t *testing.T) {
	got := Render("one\r\n\r\ntwo")
	want := "<p>one</p>\n<p>two</p>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHeadingWithInlineFormatting(t *testing.T) {
	got := Render("## Use `fmt` **well**")
	want := "<h2>Use <code>fmt</code> <strong>well</strong></h2>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
