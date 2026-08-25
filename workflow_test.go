package main

import (
	"os"
	"strings"
	"testing"
)

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// TestCIWorkflow pins the continuous-integration contract: vet, tests, the
// 100% coverage gate, and a build-only docker smoke check on every push/PR.
func TestCIWorkflow(t *testing.T) {
	ci := readFile(t, ".github/workflows/ci.yml")
	for _, want := range []string{
		"push:", "pull_request:",
		"go-version",
		"go vet ./...",
		"go test ./...",
		"scripts/check-coverage.sh",
		"docker build",
	} {
		if !strings.Contains(ci, want) {
			t.Errorf("ci.yml missing %q", want)
		}
	}
	// CI must never push images.
	if strings.Contains(ci, "docker push") || strings.Contains(ci, "push: true") {
		t.Error("ci.yml must not push images")
	}
}

// TestCDWorkflow pins the delivery contract: GHCR login with GITHUB_TOKEN,
// metadata-driven tags, push only from main/tags, SSH deploy with secrets.
func TestCDWorkflow(t *testing.T) {
	cd := readFile(t, ".github/workflows/cd.yml")
	for _, want := range []string{
		"packages: write",
		"ghcr.io",
		"GITHUB_TOKEN",
		"metadata-action",
		"build-push-action",
		"type=sha",
		"type=semver",
		"secrets.SERVER_HOST",
		"secrets.SERVER_USER",
		"secrets.SERVER_SSH_KEY",
		"docker pull",
		"127.0.0.1:8080:8080",
		"--restart unless-stopped",
	} {
		if !strings.Contains(cd, want) {
			t.Errorf("cd.yml missing %q", want)
		}
	}
	if !strings.Contains(cd, "main") {
		t.Error("cd.yml should trigger on main")
	}
	if !strings.Contains(cd, `tags:`) {
		t.Error("cd.yml should trigger on version tags")
	}
}

// TestWorkflowsContainNoHardcodedSecrets guards against credential leaks.
func TestWorkflowsContainNoHardcodedSecrets(t *testing.T) {
	for _, f := range []string{".github/workflows/ci.yml", ".github/workflows/cd.yml"} {
		content := readFile(t, f)
		if strings.Contains(content, "BEGIN OPENSSH PRIVATE KEY") ||
			strings.Contains(content, "BEGIN RSA PRIVATE KEY") {
			t.Errorf("%s contains a hardcoded private key", f)
		}
	}
}

// TestDeployAssets pins the server-side ops files.
func TestDeployAssets(t *testing.T) {
	caddy := readFile(t, "deploy/Caddyfile")
	for _, want := range []string{"adamndegwa.com", "reverse_proxy 127.0.0.1:8080", "www.adamndegwa.com", "redir"} {
		if !strings.Contains(caddy, want) {
			t.Errorf("deploy/Caddyfile missing %q", want)
		}
	}
	setup := readFile(t, "deploy/setup-server.sh")
	for _, want := range []string{"docker", "caddy", "ghcr.io"} {
		if !strings.Contains(setup, want) {
			t.Errorf("deploy/setup-server.sh missing %q", want)
		}
	}
}
