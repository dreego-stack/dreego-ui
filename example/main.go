package main

import (
	"log"
	"os"

	ui "github.com/dreego-stack/dreego-ui"
	"github.com/dreego-stack/dreego-ui/example/www"
	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/core/ssr"
)

func main() {
	app := dreego.New()
	store := dreego.NewCookieStore([]byte("dreego-ui-showcase-secret-32-bytes"))
	if err := app.SetSessionStore(store); err != nil {
		log.Fatal(err)
	}
	if _, err := ui.Register(app, ui.Options{
		Themes:     showcaseThemes(),
		Preference: ui.SessionPreference{},
	}); err != nil {
		log.Fatal(err)
	}
	if err := www.Register(app); err != nil {
		log.Fatal(err)
	}
	address := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		address = ":" + port
	}
	if err := ssr.Listen(app, address); err != nil {
		log.Fatal(err)
	}
}

func showcaseThemes() []ui.Theme {
	ocean := ui.WhiteTheme()
	ocean.ID = "ocean"
	ocean.Name = "Ocean"
	ocean.Colors.Accent = "#087e8b"
	ocean.Colors.AccentHover = "#066675"

	sand := ui.WhiteTheme()
	sand.ID = "sand"
	sand.Name = "Sand"
	sand.Colors.Canvas = "#fff8e8"
	sand.Colors.Surface = "#f6edd9"
	sand.Colors.Accent = "#855c18"
	sand.Colors.AccentHover = "#69470f"

	return []ui.Theme{ocean, sand}
}
