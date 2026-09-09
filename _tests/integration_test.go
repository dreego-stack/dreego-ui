package integration_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	ui "github.com/dreego-stack/dreego-ui"
	dreego "github.com/dreego-stack/dreego/core"
)

func TestApplicationRegistersAndSelectsManyThemes(t *testing.T) {
	customThemes := make([]ui.Theme, 3)
	for i, id := range []string{"ocean", "forest", "sand"} {
		customThemes[i] = ui.WhiteTheme()
		customThemes[i].ID = id
		customThemes[i].Name = id
	}
	customThemes[1].Colors.Canvas = "#f3faf5"
	customThemes[1].Colors.Surface = "#e9f4ec"
	customThemes[2].Colors.Canvas = "#fff8e8"
	customThemes[2].Colors.Surface = "#f6edd9"

	app := dreego.New()
	if err := app.SetLogging(false); err != nil {
		t.Fatal(err)
	}
	if err := app.SetCSRF(false); err != nil {
		t.Fatal(err)
	}
	registry, err := ui.Register(app, ui.Options{
		DefaultTheme: "forest",
		Themes:       customThemes,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(registry.Themes()); got != 5 {
		t.Fatalf("theme count = %d, want 5", got)
	}

	server := httptest.NewServer(app.Handler())
	defer server.Close()
	defaultCSS := requestCSS(t, server.URL, nil)
	if !strings.HasPrefix(defaultCSS, `:root{--dreego-color-canvas:#f3faf5;`) {
		t.Fatal("custom default theme did not render")
	}
	for _, id := range []string{"white", "black", "ocean", "forest", "sand"} {
		selector := `[data-dreego-theme="` + id + `"]`
		if !strings.Contains(defaultCSS, selector) {
			t.Fatalf("stylesheet does not contain %s", selector)
		}
	}

	cookie := postTheme(t, server.URL, "sand")
	selectedCSS := requestCSS(t, server.URL, cookie)
	if !strings.HasPrefix(selectedCSS, `:root{--dreego-color-canvas:#fff8e8;`) {
		t.Fatal("selected custom theme did not render")
	}
}

func requestCSS(t *testing.T, baseURL string, cookie *http.Cookie) string {
	t.Helper()
	request, err := http.NewRequest(http.MethodGet, baseURL+ui.DefaultStylesheetPath, nil)
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
	return string(body)
}

func postTheme(t *testing.T, baseURL, themeID string) *http.Cookie {
	t.Helper()
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	response, err := client.PostForm(baseURL+ui.DefaultSelectionPath, url.Values{
		"theme": {themeID}, "return_to": {"/settings"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", response.StatusCode)
	}
	return response.Cookies()[0]
}
