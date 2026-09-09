package ui_test

import (
	"testing"
	"time"

	ui "github.com/dreego-stack/dreego-ui"
)

func TestWhiteTheme(t *testing.T) {
	theme := ui.WhiteTheme()

	if theme.ID != "white" {
		t.Fatalf("ID = %q, want white", theme.ID)
	}
	if theme.Colors.Canvas != "#ffffff" {
		t.Fatalf("canvas = %q, want #ffffff", theme.Colors.Canvas)
	}
	if err := theme.Validate(); err != nil {
		t.Fatalf("WhiteTheme validation: %v", err)
	}
}

func TestBlackThemeUsesTrueBlackCanvas(t *testing.T) {
	theme := ui.BlackTheme()

	if theme.ID != "black" {
		t.Fatalf("ID = %q, want black", theme.ID)
	}
	if theme.Colors.Canvas != "#000000" {
		t.Fatalf("canvas = %q, want #000000", theme.Colors.Canvas)
	}
	if err := theme.Validate(); err != nil {
		t.Fatalf("BlackTheme validation: %v", err)
	}
}

func TestBuiltInThemesReturnFreshValues(t *testing.T) {
	first := ui.WhiteTheme()
	first.Colors.Accent = "#000000"
	second := ui.WhiteTheme()

	if second.Colors.Accent == first.Colors.Accent {
		t.Fatal("WhiteTheme returned shared mutable state")
	}
}

func TestThemeValidationRejectsUnsafeAndIncompleteValues(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*ui.Theme)
	}{
		{name: "empty ID", mutate: func(theme *ui.Theme) { theme.ID = "" }},
		{name: "unsafe ID", mutate: func(theme *ui.Theme) { theme.ID = `x"]{color:red}` }},
		{name: "missing name", mutate: func(theme *ui.Theme) { theme.Name = "" }},
		{name: "invalid color", mutate: func(theme *ui.Theme) { theme.Colors.Accent = "purple" }},
		{name: "unsafe length", mutate: func(theme *ui.Theme) { theme.Shape.RadiusMedium = "1rem;display:none" }},
		{name: "unsafe font", mutate: func(theme *ui.Theme) { theme.Typography.FontSans = "sans-serif;}body{" }},
		{name: "sub-millisecond motion", mutate: func(theme *ui.Theme) { theme.Motion.Fast = time.Nanosecond }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := ui.WhiteTheme()
			tt.mutate(&theme)
			if err := theme.Validate(); err == nil {
				t.Fatal("Validate() error = nil, want error")
			}
		})
	}
}

func TestCustomThemeCanAdaptBuiltInTheme(t *testing.T) {
	theme := ui.WhiteTheme()
	theme.ID = "ocean"
	theme.Name = "Ocean"
	theme.Colors.Accent = "#087e8b"
	theme.Colors.AccentHover = "#066675"

	if err := theme.Validate(); err != nil {
		t.Fatalf("custom theme validation: %v", err)
	}
}

func TestThemeValidationRejectsInsufficientTextContrast(t *testing.T) {
	theme := ui.WhiteTheme()
	theme.Colors.Text = "#fefefe"

	if err := theme.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want contrast error")
	}
}

func TestThemeValidationRejectsInsufficientAccentContrast(t *testing.T) {
	theme := ui.WhiteTheme()
	theme.Colors.Accent = "#ffffff"

	if err := theme.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want contrast error")
	}
}
