package ui_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	ui "github.com/dreego-stack/dreego-ui"
	dreego "github.com/dreego-stack/dreego/core"
)

func TestRegisterServesDefaultAndSelectedTheme(t *testing.T) {
	app := testApp(t)
	registry, err := ui.Register(app, ui.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if registry.Default().ID != "white" {
		t.Fatalf("default theme = %q, want white", registry.Default().ID)
	}

	server := httptest.NewServer(app.Handler())
	defer server.Close()

	defaultCSS := getStylesheet(t, server.URL, nil)
	if !strings.HasPrefix(defaultCSS, `:root{--dreego-color-canvas:#ffffff;`) {
		t.Fatal("default stylesheet does not start with WhiteTheme")
	}

	cookie := selectTheme(t, server.URL, "black", "/settings")
	selectedCSS := getStylesheet(t, server.URL, cookie)
	if !strings.HasPrefix(selectedCSS, `:root{--dreego-color-canvas:#000000;`) {
		t.Fatal("selected stylesheet does not start with BlackTheme")
	}
}

func TestRegisterSupportsCustomPreferenceStore(t *testing.T) {
	preference := &memoryPreference{id: "black"}
	app := testApp(t)
	_, err := ui.Register(app, ui.Options{Preference: preference})
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodGet, "/dreego-ui/theme.css", nil)
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	if !strings.HasPrefix(response.Body.String(), `:root{--dreego-color-canvas:#000000;`) {
		t.Fatal("custom preference was not used")
	}
}

func TestSelectionRejectsUnknownTheme(t *testing.T) {
	app := testApp(t)
	if _, err := ui.Register(app, ui.Options{}); err != nil {
		t.Fatal(err)
	}

	form := url.Values{"theme": {"missing"}, "return_to": {"/settings"}}
	request := httptest.NewRequest(http.MethodPost, "/dreego-ui/theme", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", response.Code)
	}
}

func TestSelectionRejectsExternalReturnURL(t *testing.T) {
	app := testApp(t)
	if _, err := ui.Register(app, ui.Options{}); err != nil {
		t.Fatal(err)
	}

	form := url.Values{"theme": {"black"}, "return_to": {"https://example.com"}}
	request := httptest.NewRequest(http.MethodPost, "/dreego-ui/theme", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", response.Code)
	}
}

func TestInvalidStoredThemeFallsBackSafely(t *testing.T) {
	app := testApp(t)
	if _, err := ui.Register(app, ui.Options{}); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(app.Handler())
	defer server.Close()
	body := getStylesheet(t, server.URL, &http.Cookie{Name: "dreego_ui_theme", Value: `x"]{color:red}`})
	if !strings.HasPrefix(body, `:root{--dreego-color-canvas:#ffffff;`) {
		t.Fatal("invalid stored theme did not fall back to WhiteTheme")
	}
}

func TestRegisterAfterBuildReturnsDreegoError(t *testing.T) {
	app := testApp(t)
	if err := app.Build(); err != nil {
		t.Fatal(err)
	}
	_, err := ui.Register(app, ui.Options{})
	if !errors.Is(err, dreego.ErrAppBuilt) {
		t.Fatalf("error = %v, want ErrAppBuilt", err)
	}
}

func testApp(t *testing.T) *dreego.App {
	t.Helper()
	app := dreego.New()
	if err := app.SetLogging(false); err != nil {
		t.Fatal(err)
	}
	if err := app.SetCSRF(false); err != nil {
		t.Fatal(err)
	}
	return app
}

func getStylesheet(t *testing.T, baseURL string, cookie *http.Cookie) string {
	t.Helper()
	request, err := http.NewRequest(http.MethodGet, baseURL+"/dreego-ui/theme.css", nil)
	if err != nil {
		t.Fatal(err)
	}
	if cookie != nil {
		request.AddCookie(cookie)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("stylesheet status = %d, want 200", response.StatusCode)
	}
	return string(body)
}

func selectTheme(t *testing.T, baseURL, themeID, returnTo string) *http.Cookie {
	t.Helper()
	client := &http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	response, err := client.PostForm(baseURL+"/dreego-ui/theme", url.Values{
		"theme": {themeID}, "return_to": {returnTo},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusSeeOther {
		t.Fatalf("selection status = %d, want 303", response.StatusCode)
	}
	if location := response.Header.Get("Location"); location != returnTo {
		t.Fatalf("location = %q, want %q", location, returnTo)
	}
	cookies := response.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookie count = %d, want 1", len(cookies))
	}
	return cookies[0]
}

type memoryPreference struct {
	id string
}

func (p *memoryPreference) Load(*http.Request) (string, error) {
	return p.id, nil
}

func (p *memoryPreference) Save(_ http.ResponseWriter, _ *http.Request, id string) error {
	p.id = id
	return nil
}
