package main

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"os"
	"testing"
)

// TestFaviconsAreWordmarkPNGs pins the favicon set: every icon must exist,
// decode as PNG, and have its documented dimensions.
func TestFaviconsAreWordmarkPNGs(t *testing.T) {
	sizes := map[string]int{
		"static/android-chrome-512x512.png": 512,
		"static/android-chrome-192x192.png": 192,
		"static/apple-touch-icon.png":       180,
		"static/favicon-32x32.png":          32,
		"static/favicon-16x16.png":          16,
	}
	for path, want := range sizes {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: %v", path, err)
			continue
		}
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			t.Errorf("%s: not a valid PNG: %v", path, err)
			continue
		}
		if got := img.Bounds().Dx(); got != want || img.Bounds().Dy() != want {
			t.Errorf("%s: %dx%d, want %dx%d", path, got, img.Bounds().Dy(), want, want)
		}
	}
}

// TestFaviconICO asserts favicon.ico is a valid ICO container carrying the
// 16x16 and 32x32 PNG images.
func TestFaviconICO(t *testing.T) {
	data, err := os.ReadFile("static/favicon.ico")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 6 {
		t.Fatal("favicon.ico too small")
	}
	// ICO header: reserved=0, type=1, count.
	if data[0] != 0 || data[1] != 0 || data[2] != 1 || data[3] != 0 {
		t.Fatalf("bad ICO header: %v", data[:4])
	}
	count := int(binary.LittleEndian.Uint16(data[4:6]))
	if count != 2 {
		t.Errorf("ICO image count = %d, want 2 (16x16 + 32x32)", count)
	}
	if len(data) < 6+16*count {
		t.Fatal("ICO directory truncated")
	}
	for i := 0; i < count; i++ {
		entry := data[6+16*i:]
		size := int(binary.LittleEndian.Uint32(entry[8:12]))
		offset := int(binary.LittleEndian.Uint32(entry[12:16]))
		if offset+size > len(data) {
			t.Fatalf("image %d out of bounds", i)
		}
		if _, err := png.Decode(bytes.NewReader(data[offset : offset+size])); err != nil {
			t.Errorf("image %d is not a PNG: %v", i, err)
		}
	}
}

// TestManifestIconsMatch keeps site.webmanifest pointing at real files.
func TestManifestIconsMatch(t *testing.T) {
	data, err := os.ReadFile("static/site.webmanifest")
	if err != nil {
		t.Fatal(err)
	}
	for _, icon := range []string{"android-chrome-192x192.png", "android-chrome-512x512.png"} {
		if !bytes.Contains(data, []byte(icon)) {
			t.Errorf("manifest missing %s", icon)
		}
	}
}
