package ui

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	dreego "github.com/dreego-stack/dreego/core"
)

const DefaultThemeSessionKey = "dreego_ui_theme"

var ErrSessionStoreUnavailable = errors.New("dreego-ui: session store is unavailable")

var sessionKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,80}$`)

type SessionPreference struct {
	Key string
}

func (p SessionPreference) Load(r *http.Request) (string, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}
	store := dreego.StoreFromCtx(r.Context())
	if store == nil {
		return "", ErrSessionStoreUnavailable
	}
	themeID, err := store.Get(r, p.key())
	if err != nil {
		return "", fmt.Errorf("dreego-ui: load session theme preference: %w", err)
	}
	return themeID, nil
}

func (p SessionPreference) Save(w http.ResponseWriter, r *http.Request, themeID string) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if !themeIDPattern.MatchString(themeID) || isReservedThemeID(themeID) {
		return errors.New("dreego-ui: cannot save invalid theme ID")
	}
	store := dreego.StoreFromCtx(r.Context())
	if store == nil {
		return ErrSessionStoreUnavailable
	}
	if err := store.Set(w, r, p.key(), themeID, nil); err != nil {
		return fmt.Errorf("dreego-ui: save session theme preference: %w", err)
	}
	return nil
}

func (p SessionPreference) Validate() error {
	if !sessionKeyPattern.MatchString(p.key()) {
		return errors.New("dreego-ui: session key contains unsupported characters")
	}
	return nil
}

func (p SessionPreference) key() string {
	if p.Key == "" {
		return DefaultThemeSessionKey
	}
	return p.Key
}
