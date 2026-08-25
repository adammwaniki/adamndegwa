---
title: adamndegwa.com
tags: Go, HTMX, CSS
description: This site. Server-rendered with Go's standard library, HTMX and pure CSS — WCAG2 tested.
icon: "05"
reading_time: 2 min
date: 2025-06-10
---

This portfolio is built with Go's standard library, HTMX and pure CSS. No frameworks, no bundlers, no build step — not even a markdown library.

## Design Principles

- Minimalism — every element earns its place
- Reusability — components that work across pages
- Inclusivity — accessible by default, proven by automated WCAG2 contrast and landmark tests
- Clarity — easy to understand over clever

## Architecture

A single Go binary using `net/http` for routing and `html/template` for rendering. A small stdlib-only markdown renderer turns the content files into HTML. HTMX handles partial page swaps. CSS handles all layout, animation and theming, driven by a brand palette shared between the stylesheet and the test suite.
