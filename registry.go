package ui

import (
	"errors"
	"fmt"
)

var ErrDuplicateTheme = errors.New("dreego-ui: duplicate theme")
var ErrUnknownDefaultTheme = errors.New("dreego-ui: unknown default theme")

type Options struct {
	DefaultTheme         string
	Themes               []Theme
	DisableBuiltInThemes bool
	StylesheetPath       string
	SelectionPath        string
	Preference           PreferenceStore
}

type Registry struct {
	themes    map[string]Theme
	order     []string
	defaultID string
}

func NewRegistry(options Options) (*Registry, error) {
	themes := options.Themes
	if !options.DisableBuiltInThemes {
		themes = append([]Theme{WhiteTheme(), BlackTheme()}, themes...)
	}
	if len(themes) == 0 {
		return nil, fmt.Errorf("%w: at least one theme is required", ErrInvalidTheme)
	}

	registry := &Registry{
		themes: make(map[string]Theme, len(themes)),
		order:  make([]string, 0, len(themes)),
	}
	for _, theme := range themes {
		if err := theme.Validate(); err != nil {
			return nil, fmt.Errorf("theme %q: %w", theme.ID, err)
		}
		if _, exists := registry.themes[theme.ID]; exists {
			return nil, fmt.Errorf("%w: %q", ErrDuplicateTheme, theme.ID)
		}
		registry.themes[theme.ID] = theme
		registry.order = append(registry.order, theme.ID)
	}

	defaultID := options.DefaultTheme
	if defaultID == "" {
		defaultID = registry.order[0]
	}
	if _, exists := registry.themes[defaultID]; !exists {
		return nil, fmt.Errorf("%w: %q", ErrUnknownDefaultTheme, defaultID)
	}
	registry.defaultID = defaultID
	return registry, nil
}

func (r *Registry) Themes() []Theme {
	if r == nil {
		return nil
	}
	themes := make([]Theme, 0, len(r.order))
	for _, id := range r.order {
		themes = append(themes, r.themes[id])
	}
	return themes
}

func (r *Registry) Default() Theme {
	if r == nil {
		return Theme{}
	}
	return r.themes[r.defaultID]
}

func (r *Registry) Resolve(id string) Theme {
	if r == nil {
		return Theme{}
	}
	if theme, exists := r.themes[id]; exists {
		return theme
	}
	return r.Default()
}

func (r *Registry) Theme(id string) (Theme, bool) {
	if r == nil {
		return Theme{}, false
	}
	theme, exists := r.themes[id]
	return theme, exists
}
