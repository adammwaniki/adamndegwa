---
title: Kubernetes from the Ground Up
tags: Cloud, Infrastructure
description: Container orchestration fundamentals demystified.
icon: "◎"
reading_time: 3 min
date: 2025-05-18
---

These are my notes from studying for the Kubernetes Cloud Native Associate certification — written for clarity, not certification prep.

## The Mental Model

Think of Kubernetes as an operating system for your infrastructure. You declare the desired state, and Kubernetes works to make reality match.

## Key Concepts

- Pods — one or more containers sharing a network namespace
- Deployments — manage scaling, rolling updates, rollbacks
- Services — stable network endpoints for your pods
- ConfigMaps and Secrets — separate configuration from code

## What I Learned

The biggest lesson was understanding when Kubernetes is the right tool and when it's over-engineering. Not every application needs container orchestration.
