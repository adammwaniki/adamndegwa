package main

import (
	"os"
	"strings"
	"testing"
)

// TestDockerfile pins the container contract: registry-qualified bases, a
// scratch final image, a static binary, and the expected runtime surface.
func TestDockerfile(t *testing.T) {
	data, err := os.ReadFile("Dockerfile")
	if err != nil {
		t.Fatal("Dockerfile missing:", err)
	}
	df := string(data)

	var froms []string
	for _, line := range strings.Split(df, "\n") {
		line = strings.TrimSpace(line)
		if rest, ok := strings.CutPrefix(line, "FROM "); ok {
			froms = append(froms, rest)
		}
	}
	if len(froms) < 2 {
		t.Fatalf("expected a multi-stage build, got FROM lines: %v", froms)
	}
	for _, f := range froms {
		if strings.HasPrefix(f, "scratch") {
			continue // scratch is a reserved base, not a registry image
		}
		if !strings.HasPrefix(f, "docker.io/") {
			t.Errorf("FROM %q is not fully qualified with docker.io/ (podman-compatible)", f)
		}
	}
	if froms[len(froms)-1] != "scratch" {
		t.Errorf("final stage = %q, want scratch", froms[len(froms)-1])
	}

	for _, want := range []string{
		"CGO_ENABLED=0", // static binary, required for scratch
		"EXPOSE 8080",
		"ENTRYPOINT",
		"/views", "/static", "/content", // runtime assets must be copied in
	} {
		if !strings.Contains(df, want) {
			t.Errorf("Dockerfile missing %q", want)
		}
	}
}

func TestDockerignore(t *testing.T) {
	data, err := os.ReadFile(".dockerignore")
	if err != nil {
		t.Fatal(".dockerignore missing:", err)
	}
	for _, want := range []string{".git", "portfolio-remix", "references", "scripts"} {
		if !strings.Contains(string(data), want) {
			t.Errorf(".dockerignore should exclude %q", want)
		}
	}
	// Content must NOT be ignored — the site reads it at runtime.
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "content" || strings.TrimSpace(line) == "content/" {
			t.Error(".dockerignore must not exclude content/")
		}
	}
}
