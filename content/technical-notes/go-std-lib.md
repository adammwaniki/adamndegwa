---
title: Building with Go's Standard Library
tags: Go, Architecture
description: Why Go's standard library is enough for most web apps.
icon: "{ }"
reading_time: 3 min
date: 2025-05-01
---

Go's standard library is one of the most underrated tools in a developer's arsenal. While the ecosystem is full of frameworks, the standard library alone provides everything you need to build production-grade web applications.

## Why Standard Library?

The `net/http` package gives you a fully capable HTTP server. The `html/template` package provides safe, composable HTML templating. The `encoding/json` package handles serialisation. The `database/sql` package provides a clean interface to any SQL database.

## Project Structure

- A `cmd/` directory for entry points
- An `internal/` directory for business logic
- A `views/` directory for templates
- A `static/` directory for assets

## The Trade-Off

You write slightly more boilerplate upfront. In exchange, you get zero hidden magic, complete control over your request lifecycle, and a codebase any Go developer can read without learning a framework first.
