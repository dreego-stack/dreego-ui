package ui_test

import (
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
		{"Navbar.dreego", []string{"<header", "<nav", `aria-label=`, "{#slot brand}", "{#slot actions}"}},
		{"Sidebar.dreego", []string{"<aside", "<nav", `aria-label=`, "{#slot footer}"}},
		{"Footer.dreego", []string{"<footer", "{#slot meta}"}},
		{"PriceCard.dreego", []string{"<article", "<h2", "<a", "{#slot features}", `aria-labelledby=`}},
		{"PageShell.dreego", []string{"<main", "{#slot header}", "{#slot sidebar}", "{#slot footer}"}},
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
		})
	}
}
