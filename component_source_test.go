package ui_test

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComponentSourcesUseThemeTokensAndNativeControls(t *testing.T) {
	tests := []struct {
		file  string
		wants []string
	}{
		{"Button.dreego", []string{"<button", ":focus-visible", "--dreego-color-accent"}},
		{"Toggle.dreego", []string{`type="checkbox"`, `role="switch"`, "<label"}},
		{"Card.dreego", []string{"<article", "{#slot header}", "{#slot footer}"}},
		{"CodeBox.dreego", []string{"<figure", "<figcaption", "<pre", "<code"}},
		{"Icon.dreego", []string{"<svg", "<title", `aria-hidden="true"`, `role="img"`}},
		{"ThemePicker.dreego", []string{"<form", "<label", "<select", `method="post"`, `name="csrf_token"`}},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join("components", "dreegoui", tt.file))
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range tt.wants {
				if !strings.Contains(string(content), want) {
					t.Errorf("%s does not contain %q", tt.file, want)
				}
			}
			hash := sha256.Sum256(content)
			if hash[0]>>4 < 10 {
				t.Fatal("component scope hash starts with a digit, which Dreego v0.6.4 emits as invalid unquoted CSS")
			}
		})
	}
}

func TestShowcaseScopeHashStartsWithLetter(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("example", "www", "routes", "+page.dreego"))
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(content)
	if hash[0]>>4 < 10 {
		t.Fatal("showcase scope hash starts with a digit, which Dreego v0.6.4 emits as invalid unquoted CSS")
	}
}
