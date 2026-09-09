package ui

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

const DefaultThemeCookieName = "dreego_ui_theme"

type PreferenceStore interface {
	Load(r *http.Request) (string, error)
	Save(w http.ResponseWriter, r *http.Request, themeID string) error
}

type CookiePreference struct {
	Name   string
	MaxAge time.Duration
	Secure bool
}

var cookieNamePattern = regexp.MustCompile(`^[A-Za-z0-9_]{1,64}$`)

func (p CookiePreference) Load(r *http.Request) (string, error) {
	cookie, err := r.Cookie(p.name())
	if errors.Is(err, http.ErrNoCookie) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("dreego-ui: load theme preference: %w", err)
	}
	return cookie.Value, nil
}

func (p CookiePreference) Save(w http.ResponseWriter, _ *http.Request, themeID string) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if !themeIDPattern.MatchString(themeID) || isReservedThemeID(themeID) {
		return errors.New("dreego-ui: cannot save invalid theme ID")
	}
	maxAge := p.MaxAge
	if maxAge == 0 {
		maxAge = 365 * 24 * time.Hour
	}
	http.SetCookie(w, &http.Cookie{
		Name:     p.name(),
		Value:    themeID,
		Path:     "/",
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: true,
		Secure:   p.Secure,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

func (p CookiePreference) Validate() error {
	if !cookieNamePattern.MatchString(p.name()) {
		return errors.New("dreego-ui: cookie name must contain only letters, numbers, or underscores")
	}
	if p.MaxAge < 0 {
		return errors.New("dreego-ui: cookie max age must not be negative")
	}
	if p.MaxAge > 0 && p.MaxAge < time.Second {
		return errors.New("dreego-ui: cookie max age must be at least one second")
	}
	return nil
}

func (p CookiePreference) name() string {
	if p.Name == "" {
		return DefaultThemeCookieName
	}
	return p.Name
}
