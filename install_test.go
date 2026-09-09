package ui_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	ui "github.com/dreego-stack/dreego-ui"
)

func TestInstallComponentsWritesCompleteStableSet(t *testing.T) {
	destination := filepath.Join(t.TempDir(), "www", "components", "dreego-ui")
	installed, err := ui.InstallComponents(ui.InstallOptions{Destination: destination})
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"Button.dreego",
		"Card.dreego",
		"CodeBox.dreego",
		"Footer.dreego",
		"Icon.dreego",
		"Navbar.dreego",
		"PageShell.dreego",
		"PriceCard.dreego",
		"Sidebar.dreego",
		"ThemePicker.dreego",
		"Toggle.dreego",
		"theme_option.go",
	}
	got := make([]string, len(installed))
	for i, path := range installed {
		got[i] = filepath.Base(path)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("installed file %s: %v", path, err)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("installed files = %v, want %v", got, want)
	}
}

func TestInstallComponentsDoesNotOverwriteExistingFiles(t *testing.T) {
	destination := filepath.Join(t.TempDir(), "components")
	if err := os.MkdirAll(destination, 0o755); err != nil {
		t.Fatal(err)
	}
	existing := filepath.Join(destination, "Button.dreego")
	if err := os.WriteFile(existing, []byte("user content"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := ui.InstallComponents(ui.InstallOptions{Destination: destination})
	if !errors.Is(err, ui.ErrComponentExists) {
		t.Fatalf("error = %v, want ErrComponentExists", err)
	}
	content, readErr := os.ReadFile(existing)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(content) != "user content" {
		t.Fatalf("existing content changed to %q", content)
	}
	entries, readErr := os.ReadDir(destination)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if len(entries) != 1 {
		t.Fatalf("entry count = %d, want 1", len(entries))
	}
}

func TestInstallComponentsRejectsFileDestination(t *testing.T) {
	destination := filepath.Join(t.TempDir(), "components")
	if err := os.WriteFile(destination, []byte("file"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := ui.InstallComponents(ui.InstallOptions{Destination: destination}); err == nil {
		t.Fatal("InstallComponents error = nil, want error")
	}
}
