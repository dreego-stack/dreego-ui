package ui

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

var ErrComponentExists = errors.New("dreego-ui: component already exists")

const defaultComponentDestination = "www/components/dreegoui"

//go:embed components/dreegoui/*
var componentSources embed.FS

type InstallOptions struct {
	Destination string
}

func InstallComponents(options InstallOptions) ([]string, error) {
	destination := options.Destination
	if destination == "" {
		destination = defaultComponentDestination
	}
	info, err := os.Stat(destination)
	if err == nil && !info.IsDir() {
		return nil, fmt.Errorf("dreego-ui: component destination is not a directory: %s", destination)
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("dreego-ui: inspect component destination: %w", err)
	}

	entries, err := fs.ReadDir(componentSources, "components/dreegoui")
	if err != nil {
		return nil, fmt.Errorf("dreego-ui: read embedded components: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, entry := range entries {
		target := filepath.Join(destination, entry.Name())
		if _, err := os.Lstat(target); err == nil {
			return nil, fmt.Errorf("%w: %s", ErrComponentExists, target)
		} else if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("dreego-ui: inspect component target: %w", err)
		}
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return nil, fmt.Errorf("dreego-ui: create component destination: %w", err)
	}

	installed := make([]string, 0, len(entries))
	for _, entry := range entries {
		sourcePath := "components/dreegoui/" + entry.Name()
		content, err := componentSources.ReadFile(sourcePath)
		if err != nil {
			return nil, fmt.Errorf("dreego-ui: read component %s: %w", entry.Name(), err)
		}
		target := filepath.Join(destination, entry.Name())
		if err := os.WriteFile(target, content, 0o644); err != nil {
			return nil, fmt.Errorf("dreego-ui: write component %s: %w", entry.Name(), err)
		}
		installed = append(installed, target)
	}
	return installed, nil
}
