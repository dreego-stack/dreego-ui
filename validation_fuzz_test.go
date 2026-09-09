package ui_test

import (
	"strings"
	"testing"

	ui "github.com/dreego-stack/dreego-ui"
)

func FuzzThemeColorValidation(f *testing.F) {
	for _, seed := range []string{"#006d77", "purple", "#fff", "#00000000", "red;}body{"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, value string) {
		theme := ui.WhiteTheme()
		theme.Colors.Accent = value
		err := theme.Validate()
		if err == nil && !isOpaqueHex(value) {
			t.Fatalf("Validate accepted unsupported color %q", value)
		}
	})
}

func FuzzThemeIDValidation(f *testing.F) {
	for _, seed := range []string{"ocean", "black-2", "", "System", `x"]{color:red}`} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, value string) {
		theme := ui.WhiteTheme()
		theme.ID = value
		err := theme.Validate()
		if err == nil && !isSafeThemeID(value) {
			t.Fatalf("Validate accepted unsupported ID %q", value)
		}
	})
}

func isOpaqueHex(value string) bool {
	if len(value) != 7 || value[0] != '#' {
		return false
	}
	for _, char := range value[1:] {
		if !strings.ContainsRune("0123456789abcdefABCDEF", char) {
			return false
		}
	}
	return true
}

func isSafeThemeID(value string) bool {
	if value == "auto" || value == "default" || value == "system" {
		return false
	}
	if len(value) == 0 || len(value) > 63 || value[0] < 'a' || value[0] > 'z' {
		return false
	}
	for _, char := range value[1:] {
		if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
			return false
		}
	}
	return true
}
