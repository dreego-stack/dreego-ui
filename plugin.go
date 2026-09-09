package ui

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	dreego "github.com/dreego-stack/dreego/core"
)

const DefaultStylesheetPath = "/dreego-ui/theme.css"
const DefaultSelectionPath = "/dreego-ui/theme"

func Register(app *dreego.App, options Options) (*Registry, error) {
	if app == nil {
		return nil, errors.New("dreego-ui: app is nil")
	}
	registry, err := NewRegistry(options)
	if err != nil {
		return nil, err
	}
	preference := options.Preference
	if preference == nil {
		preference = CookiePreference{}
	}
	if validator, ok := preference.(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return nil, err
		}
	}

	stylesheetPath := options.StylesheetPath
	if stylesheetPath == "" {
		stylesheetPath = DefaultStylesheetPath
	}
	selectionPath := options.SelectionPath
	if selectionPath == "" {
		selectionPath = DefaultSelectionPath
	}
	if stylesheetPath == selectionPath {
		return nil, errors.New("dreego-ui: stylesheet and selection paths must differ")
	}

	if err := app.Register(http.MethodGet, stylesheetPath, stylesheetHandler(registry, preference)); err != nil {
		return nil, fmt.Errorf("dreego-ui: register stylesheet: %w", err)
	}
	if err := app.Register(http.MethodPost, selectionPath, selectionHandler(registry, preference)); err != nil {
		return nil, fmt.Errorf("dreego-ui: register theme selection: %w", err)
	}
	return registry, nil
}

func stylesheetHandler(registry *Registry, preference PreferenceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		themeID, err := preference.Load(r)
		if err != nil {
			http.Error(w, "unable to load theme", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "private, no-cache")
		w.Header().Set("Vary", "Cookie")
		_, _ = w.Write([]byte(registry.StylesheetFor(themeID)))
	}
}

func selectionHandler(registry *Registry, preference PreferenceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid theme selection", http.StatusBadRequest)
			return
		}
		themeID := r.Form.Get("theme")
		if _, exists := registry.Theme(themeID); !exists {
			http.Error(w, "unknown theme", http.StatusBadRequest)
			return
		}
		returnTo := r.Form.Get("return_to")
		if returnTo == "" {
			returnTo = "/"
		}
		if !isLocalReturnURL(returnTo) {
			http.Error(w, "invalid return URL", http.StatusBadRequest)
			return
		}
		if err := preference.Save(w, r, themeID); err != nil {
			http.Error(w, "unable to save theme", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, returnTo, http.StatusSeeOther)
	}
}

func isLocalReturnURL(value string) bool {
	if !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") || strings.Contains(value, `\`) {
		return false
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}
	return parsed.Scheme == "" && parsed.Host == ""
}
