package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	ui "github.com/dreego-stack/dreego-ui"
)

func TestInstalledComponentsGenerateAndCompileWithDreego(t *testing.T) {
	directory := t.TempDir()
	writeFixtureFile(t, directory, "go.mod", `module componentfixture

go 1.22

require github.com/dreego-stack/dreego v0.6.4
`)
	writeFixtureFile(t, directory, "www/dreego.config.json", `{
    "logging": {"enabled": false},
    "redirects": [],
    "rewrites": []
}`)
	writeFixtureFile(t, directory, "www/routes/+page.dreego", `import Button "components/dreegoui/Button.dreego"
import Card "components/dreegoui/Card.dreego"
import CodeBox "components/dreegoui/CodeBox.dreego"
import Icon "components/dreegoui/Icon.dreego"
import ThemePicker "components/dreegoui/ThemePicker.dreego"
import Toggle "components/dreegoui/Toggle.dreego"

<head><title>Components</title></head>

<server>
    themes := []dreegoui.ThemeOption{
        {ID: "white", Name: "White"},
        {ID: "black", Name: "Black"},
        {ID: "ocean", Name: "Ocean"},
    }
</server>

<body>
    <@Button disabled={false}>Continue</@Button>
    <@Toggle id="updates" name="updates" label="Product updates" checked={true} disabled={false}/>
    <@Card>
        {#slot header}<h2>Account</h2>{/slot}
        <p>Private by default.</p>
        {#slot footer}<a href="/settings">Settings</a>{/slot}
    </@Card>
    <@CodeBox code="go test ./..." language="shell"/>
    <@Icon name="shield" label="Protected"/>
    <@ThemePicker id="theme-choice" themes={themes} csrfToken="token" selected="white" returnTo="/settings"/>
</body>`)
	writeFixtureFile(t, directory, "www/routes/render_test.go", `package routes

import (
    "net/http/httptest"
    "strings"
    "testing"

    dreego "github.com/dreego-stack/dreego/core"
)

func TestComponentPageMarkup(t *testing.T) {
    response := httptest.NewRecorder()
    request := httptest.NewRequest("GET", "/", nil)
    html, err := renderIndex(dreego.NewSSR(response, request))
    if err != nil {
        t.Fatal(err)
    }
    for _, want := range []string{
        "<button", "role=\"switch\"", "<article", "<figure", "<form", "<svg",
        "name=\"csrf_token\"", "value=\"black\"", "go test ./...",
    } {
        if !strings.Contains(html, want) {
            t.Errorf("rendered HTML does not contain %q", want)
        }
    }
}
`)

	_, err := ui.InstallComponents(ui.InstallOptions{
		Destination: filepath.Join(directory, "www", "components", "dreegoui"),
	})
	if err != nil {
		t.Fatal(err)
	}

	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = directory
	if output, err := tidy.CombinedOutput(); err != nil {
		t.Fatalf("resolve fixture module: %v\n%s", err, output)
	}
	generate := exec.Command("go", "run", "-mod=mod", "github.com/dreego-stack/dreego/cli/dreego", "generate", "--force")
	generate.Dir = directory
	if output, err := generate.CombinedOutput(); err != nil {
		t.Fatalf("dreego generate: %v\n%s", err, output)
	}
	check := exec.Command("go", "run", "-mod=mod", "github.com/dreego-stack/dreego/cli/dreego", "generate", "--check")
	check.Dir = directory
	if output, err := check.CombinedOutput(); err != nil {
		t.Fatalf("dreego generate check: %v\n%s", err, output)
	}
	compile := exec.Command("go", "test", "./...")
	compile.Dir = directory
	if output, err := compile.CombinedOutput(); err != nil {
		routeSource, _ := os.ReadFile(filepath.Join(directory, "www", "routes", "dree.go"))
		t.Fatalf("compile generated components: %v\n%s\n%s", err, output, routeSource)
	}

	generated, err := os.ReadFile(filepath.Join(directory, "www", "components", "dreegoui", "dree.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, function := range []string{"Button", "Toggle", "Card", "CodeBox", "Icon", "ThemePicker"} {
		if !strings.Contains(string(generated), "func "+function+"(") {
			t.Errorf("generated components do not contain %s", function)
		}
	}
}

func writeFixtureFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
