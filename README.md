# adamndegwa.com

Personal portfolio of Adam Ndegwa | Software Engineer. Built to the brand system (Ink / Pine /
Paper / Gold palette; Cinzel, Big Shoulders Inline, Google Sans Flex and Google
Sans Code type; light and AMOLED dark themes).

## Stack

- **Go standard library only** — `net/http` routing, `html/template` rendering,
  and a small in-house markdown renderer (`internal/markdown`). No third-party
  Go modules.
- **HTMX** (vendored at `static/htmx.min.js`) for partial page swaps.
- **Hand-written CSS** (`static/style.css`).

## Writing content

Drop a markdown file into the matching `content/<section>/` directory:

```markdown
---
title: My Post
tags: Go, Architecture
description: One-line summary for listings and SEO.
icon: "{ }"
reading_time: 3 min
date: 2025-05-01
updated: 2025-06-15        # optional
demo_url: https://...      # optional, projects
---

Body markdown here. Supports `#`–`####` headings, `-` lists, `inline code`,
fenced code blocks, `**bold**` and `[links](https://…)`.
```

Cards sort by date ascending. Sections: `technical-notes`, `projects`,
`musings`, `the-bullshitters`.

## Running

```sh
go run .          # serves on :8080 (PORT env var overrides)
```

## Container

Multi-stage build ending in a `scratch` image (binary + views/static/content
only, non-root, ~10 MB). All base images are registry-qualified with
`docker.io/` so both Docker and Podman work:

```sh
docker build -t adamndegwa .        # or: podman build -t adamndegwa .
docker run -p 8080:8080 adamndegwa  # or: podman run -p 8080:8080 adamndegwa
```

## Testing

Everything is test-driven; internal packages hold 100% unit-test coverage:

```sh
go test ./...              # full suite
go vet ./...
scripts/check-coverage.sh  # coverage gate
```

WCAG2 AA compliance is enforced by tests, not by inspection:

- `internal/wcag` computes WCAG contrast ratios for every brand text pairing in
  both themes and fails below AA (4.5:1 body, 3.0:1 large/UI). It also
  cross-checks that `static/style.css` tokens match `internal/brand`.
- `internal/server/a11y_test.go` renders every page type and asserts landmarks,
  `lang`, exactly one `<h1>`, skip link, named icon buttons, `aria-expanded`
  wiring and no placeholder links.

## SEO

Canonical URLs, Open Graph, Twitter cards and JSON-LD in the layout;
`/sitemap.xml`, `/robots.txt`, `/feed.xml` and `/search-index.json` served from
`internal/seo`.

## Layout

```txt
main.go             thin entry point
internal/brand      palette + fonts (single source of truth)
internal/markdown   stdlib-only markdown renderer
internal/content    section/card model + markdown loading
internal/seo        sitemap, robots, RSS, search index
internal/server     routing, handlers, HTMX partials
internal/wcag       WCAG contrast math + compliance proofs
views/              html/template pages and partials
static/             style.css, htmx.min.js, favicons
content/            markdown content (one file per card)
```
