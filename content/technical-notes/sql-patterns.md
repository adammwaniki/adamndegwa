---
title: SQL Patterns That Scale
tags: Databases, Performance
description: Indexing, query patterns and schema design lessons.
icon: "▦"
reading_time: 2 min
date: 2025-05-25
---

Writing SQL that works is easy. Writing SQL that scales is a discipline.

## Indexing Strategy

Indexes are not free. Every index speeds up reads but slows down writes. Index the columns in your WHERE clauses and JOIN conditions.

## Query Patterns

- Use EXISTS instead of IN for subqueries
- Avoid SELECT * in production
- Use EXPLAIN ANALYZE to understand query plans
- Paginate with keyset pagination for large datasets

## Schema Design

Normalise first, denormalise intentionally. Start with a clean relational model, then introduce denormalisation only with measured evidence.
