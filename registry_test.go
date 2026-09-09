package ui_test

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"testing"

	ui "github.com/dreego-stack/dreego-ui"
)

func TestRegistryIncludesBuiltInsAndCustomThemesInStableOrder(t *testing.T) {
	ocean := ui.WhiteTheme()
	ocean.ID = "ocean"
	ocean.Name = "Ocean"
	ocean.Colors.Accent = "#087e8b"
	ocean.Colors.AccentHover = "#066675"

	registry, err := ui.NewRegistry(ui.Options{
		DefaultTheme: "ocean",
		Themes:       []ui.Theme{ocean},
	})
	if err != nil {
		t.Fatal(err)
	}

	themes := registry.Themes()
	if len(themes) != 3 {
		t.Fatalf("theme count = %d, want 3", len(themes))
	}
	for i, want := range []string{"white", "black", "ocean"} {
		if themes[i].ID != want {
			t.Fatalf("theme %d ID = %q, want %q", i, themes[i].ID, want)
		}
	}
	if registry.Default().ID != "ocean" {
		t.Fatalf("default ID = %q, want ocean", registry.Default().ID)
	}
}

func TestRegistrySupportsManyCustomThemes(t *testing.T) {
	themes := make([]ui.Theme, 3)
	for i, id := range []string{"ocean", "forest", "sand"} {
		themes[i] = ui.WhiteTheme()
		themes[i].ID = id
		themes[i].Name = id
	}

	registry, err := ui.NewRegistry(ui.Options{Themes: themes})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(registry.Themes()); got != 5 {
		t.Fatalf("theme count = %d, want 5", got)
	}
}

func TestRegistryRejectsDuplicateTheme(t *testing.T) {
	_, err := ui.NewRegistry(ui.Options{Themes: []ui.Theme{ui.WhiteTheme()}})
	if !errors.Is(err, ui.ErrDuplicateTheme) {
		t.Fatalf("error = %v, want ErrDuplicateTheme", err)
	}
}

func TestRegistryRejectsUnknownDefault(t *testing.T) {
	_, err := ui.NewRegistry(ui.Options{DefaultTheme: "missing"})
	if !errors.Is(err, ui.ErrUnknownDefaultTheme) {
		t.Fatalf("error = %v, want ErrUnknownDefaultTheme", err)
	}
}

func TestRegistryCanDisableBuiltInThemes(t *testing.T) {
	ocean := ui.WhiteTheme()
	ocean.ID = "ocean"
	ocean.Name = "Ocean"

	registry, err := ui.NewRegistry(ui.Options{
		DefaultTheme:         "ocean",
		Themes:               []ui.Theme{ocean},
		DisableBuiltInThemes: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(registry.Themes()); got != 1 {
		t.Fatalf("theme count = %d, want 1", got)
	}
}

func TestRegistryReturnsCopies(t *testing.T) {
	registry, err := ui.NewRegistry(ui.Options{})
	if err != nil {
		t.Fatal(err)
	}

	themes := registry.Themes()
	themes[0].Name = "Changed"
	if registry.Themes()[0].Name == "Changed" {
		t.Fatal("Themes returned mutable registry state")
	}
}

func TestResolveFallsBackToDefault(t *testing.T) {
	registry, err := ui.NewRegistry(ui.Options{DefaultTheme: "black"})
	if err != nil {
		t.Fatal(err)
	}

	if got := registry.Resolve("missing"); got.ID != "black" {
		t.Fatalf("resolved ID = %q, want black", got.ID)
	}
	if got := registry.Resolve("white"); got.ID != "white" {
		t.Fatalf("resolved ID = %q, want white", got.ID)
	}
}

func TestStylesheetIsDeterministicAndNamespaced(t *testing.T) {
	registry, err := ui.NewRegistry(ui.Options{DefaultTheme: "black"})
	if err != nil {
		t.Fatal(err)
	}

	first := registry.Stylesheet()
	second := registry.Stylesheet()
	if first != second {
		t.Fatal("Stylesheet output is not deterministic")
	}

	wants := []string{
		`:root{--dreego-color-canvas:#000000;`,
		`[data-dreego-theme="white"]{`,
		`[data-dreego-theme="black"]{`,
		`--dreego-color-accent:`,
		`--dreego-shape-radius-medium:`,
		`@media (prefers-reduced-motion:reduce)`,
	}
	for _, want := range wants {
		if !strings.Contains(first, want) {
			t.Errorf("stylesheet does not contain %q", want)
		}
	}
}

func TestStylesheetGolden(t *testing.T) {
	registry, err := ui.NewRegistry(ui.Options{DefaultTheme: "black"})
	if err != nil {
		t.Fatal(err)
	}

	got := fmt.Sprintf("%x", sha256.Sum256([]byte(registry.Stylesheet())))
	const want = "c87be978b7d94f4063a87dd43e822faa4de92ced3ecb37fc0c04a1a408f3c935"
	if got != want {
		t.Fatalf("stylesheet hash = %s, want %s", got, want)
	}
}
