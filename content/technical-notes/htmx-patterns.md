---
title: HTMX Patterns for Server-Driven UIs
tags: HTMX, Frontend
description: Partial swaps, lazy loading and progressive enhancement.
icon: "</>"
reading_time: 2 min
date: 2025-05-10
---

HTMX lets you build dynamic interfaces by returning HTML from the server instead of JSON. It's a return to the architecture the web was designed for, with modern capabilities layered on top.

## Core Patterns

- Use `hx-get` and `hx-swap` for partial page updates
- Use `hx-trigger="revealed"` for lazy loading
- Use `hx-push-url` to maintain browser history
- Use `hx-indicator` for loading states

## Progressive Enhancement

The best HTMX applications work without JavaScript enabled. Every link and form should function as standard HTML first. HTMX then enhances the experience. If the JavaScript fails to load, the page still works.
