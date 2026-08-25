package server

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestNewFailsWithoutPartials(t *testing.T) {
	views := t.TempDir()
	if err := os.WriteFile(views+"/layout.html", []byte(`{{define "layout.html"}}x{{end}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := New(views, t.TempDir()); err == nil {
		t.Error("expected error when partials glob matches nothing")
	}
}

func TestRunPropagatesNewError(t *testing.T) {
	if err := Run(":0", t.TempDir(), t.TempDir(), func(string, http.Handler) error {
		t.Error("listener should not be called when New fails")
		return nil
	}); err == nil {
		t.Error("expected error from Run")
	}
}

func TestRunFromEnv(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir("../.."); err != nil { // repo root: views/, static/, content/
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(cwd) })

	old := listenAndServe
	t.Cleanup(func() { listenAndServe = old })

	var addr string
	listenAndServe = func(a string, h http.Handler) error { addr = a; return nil }

	t.Setenv("PORT", "9090")
	if err := RunFromEnv(); err != nil {
		t.Fatal(err)
	}
	if addr != ":9090" {
		t.Errorf("addr = %q, want :9090", addr)
	}

	t.Setenv("PORT", "")
	os.Unsetenv("PORT")
	addr = ""
	if err := RunFromEnv(); err != nil {
		t.Fatal(err)
	}
	if addr != ":8080" {
		t.Errorf("addr = %q, want default :8080", addr)
	}
}

func TestRenderPageTemplateErrors(t *testing.T) {
	// A Server whose templates lack the expected names exercises the error
	// branches of renderPage.
	broken := &Server{tpl: template.New("empty")}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	broken.renderPage(rec, req, map[string]any{"Title": "t"})
	if rec.Code != 500 {
		t.Errorf("full render with broken template: status = %d, want 500", rec.Code)
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("HX-Request", "true")
	broken.renderPage(rec, req, map[string]any{"Title": "t"})
	if rec.Code != 500 {
		t.Errorf("fragment render with broken template: status = %d, want 500", rec.Code)
	}
	if rec.Header().Get("HX-Title") != "t" {
		t.Errorf("HX-Title should be set before the error, got %q", rec.Header().Get("HX-Title"))
	}
}

func TestHandleSectionUnknownDirectly(t *testing.T) {
	s := newTestServer(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/not-a-section", nil)
	s.handleSection(rec, req)
	if rec.Code != 404 {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestDictHelper(t *testing.T) {
	m := dict("a", 1, "b", 2, "dangling")
	if m["a"] != 1 || m["b"] != 2 || len(m) != 2 {
		t.Errorf("dict = %v", m)
	}
	if dict() == nil {
		t.Error("dict() should return an empty map")
	}
	// Non-string key is coerced to "".
	m = dict(42, "x")
	if _, ok := m[""]; !ok {
		t.Errorf("non-string key should map to empty string, got %v", m)
	}
}

func TestHTMXFragmentContentType(t *testing.T) {
	s := newTestServer(t)
	rec := get(t, s, "/contact", map[string]string{"HX-Request": "true"})
	if strings.Contains(rec.Body.String(), "<nav") {
		t.Error("fragment should not contain the nav (it lives in the layout)")
	}
	if rec.Header().Get("HX-Title") != "Contact | Adam Ndegwa" {
		t.Errorf("HX-Title = %q", rec.Header().Get("HX-Title"))
	}
}
