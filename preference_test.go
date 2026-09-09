package ui_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	ui "github.com/dreego-stack/dreego-ui"
	dreego "github.com/dreego-stack/dreego/core"
)

func TestCookiePreferenceUsesSecureDefaults(t *testing.T) {
	preference := ui.CookiePreference{}
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	response := httptest.NewRecorder()

	if err := preference.Save(response, request, "black"); err != nil {
		t.Fatal(err)
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookie count = %d, want 1", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != ui.DefaultThemeCookieName || cookie.Value != "black" {
		t.Fatalf("cookie = %s=%s", cookie.Name, cookie.Value)
	}
	if !cookie.HttpOnly {
		t.Fatal("theme cookie must be HttpOnly")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("SameSite = %v, want Lax", cookie.SameSite)
	}
	if cookie.Path != "/" || cookie.MaxAge <= 0 {
		t.Fatalf("cookie path/max-age = %q/%d", cookie.Path, cookie.MaxAge)
	}
}

func TestCookiePreferenceLoadsSavedValue(t *testing.T) {
	preference := ui.CookiePreference{Name: "site_theme"}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: "site_theme", Value: "ocean"})

	got, err := preference.Load(request)
	if err != nil {
		t.Fatal(err)
	}
	if got != "ocean" {
		t.Fatalf("theme ID = %q, want ocean", got)
	}
}

func TestCookiePreferenceValidation(t *testing.T) {
	tests := []ui.CookiePreference{
		{Name: "bad-name"},
		{MaxAge: time.Nanosecond},
		{MaxAge: -time.Second},
	}
	for _, preference := range tests {
		if err := preference.Validate(); err == nil {
			t.Fatalf("Validate(%+v) error = nil, want error", preference)
		}
	}
}

func TestCookiePreferenceRejectsUnsafeThemeID(t *testing.T) {
	preference := ui.CookiePreference{}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", nil)

	if err := preference.Save(response, request, "bad theme; Path=/admin"); err == nil {
		t.Fatal("Save error = nil, want error")
	}
	if got := len(response.Result().Cookies()); got != 0 {
		t.Fatalf("cookie count = %d, want 0", got)
	}
}

func TestRegisterRejectsInvalidCookiePreference(t *testing.T) {
	app := testApp(t)
	_, err := ui.Register(app, ui.Options{
		Preference: ui.CookiePreference{Name: "bad-name"},
	})
	if err == nil {
		t.Fatal("Register error = nil, want error")
	}
}

func TestPreferenceErrorsReturnInternalServerError(t *testing.T) {
	preference := failingPreference{}
	app := testApp(t)
	if _, err := ui.Register(app, ui.Options{Preference: preference}); err != nil {
		t.Fatal(err)
	}

	loadRequest := httptest.NewRequest(http.MethodGet, ui.DefaultStylesheetPath, nil)
	loadResponse := httptest.NewRecorder()
	app.ServeHTTP(loadResponse, loadRequest)
	if loadResponse.Code != http.StatusInternalServerError {
		t.Fatalf("load status = %d, want 500", loadResponse.Code)
	}

	saveRequest := httptest.NewRequest(http.MethodPost, ui.DefaultSelectionPath, nil)
	saveRequest.Form = map[string][]string{"theme": {"black"}, "return_to": {"/"}}
	saveResponse := httptest.NewRecorder()
	app.ServeHTTP(saveResponse, saveRequest)
	if saveResponse.Code != http.StatusInternalServerError {
		t.Fatalf("save status = %d, want 500", saveResponse.Code)
	}
}

func TestSessionPreferencePersistsTheme(t *testing.T) {
	preference := ui.SessionPreference{}
	store := dreego.NewCookieStore([]byte("dreego-ui-test-secret-at-least-32-bytes"))
	app := testApp(t)
	if err := app.SetSessionStore(store); err != nil {
		t.Fatal(err)
	}
	if _, err := ui.Register(app, ui.Options{Preference: preference}); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(app.Handler())
	defer server.Close()
	cookie := selectTheme(t, server.URL, "black", "/")
	stylesheet := getStylesheet(t, server.URL, cookie)
	if !strings.HasPrefix(stylesheet, `:root{--dreego-color-canvas:#000000;`) {
		t.Fatal("session preference did not select BlackTheme")
	}
}

type failingPreference struct{}

func (failingPreference) Load(*http.Request) (string, error) {
	return "", errors.New("load failed")
}

func (failingPreference) Save(http.ResponseWriter, *http.Request, string) error {
	return errors.New("save failed")
}
