# Themes

Dreego UI registers `WhiteTheme` and `BlackTheme` by default. Applications may
add any number of custom themes by adapting a built-in value or constructing a
complete `Theme`.

## Register the plugin

```go
package main

import (
    "log"

    ui "github.com/dreego-stack/dreego-ui"
    dreego "github.com/dreego-stack/dreego/core"
)

func main() {
    ocean := ui.WhiteTheme()
    ocean.ID = "ocean"
    ocean.Name = "Ocean"
    ocean.Colors.Accent = "#087e8b"
    ocean.Colors.AccentHover = "#066675"

    app := dreego.New()
    _, err := ui.Register(app, ui.Options{
        DefaultTheme: "white",
        Themes:       []ui.Theme{ocean},
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

Include the registered stylesheet in the application layout:

```html
<link rel="stylesheet" href="/dreego-ui/theme.css">
```

The stylesheet resolves the current preference on every request and emits the
selected theme as `:root`. It also includes attribute selectors for every
registered theme, so progressive enhancement can preview a theme immediately:

```html
<html data-dreego-theme="black">
```

## No-JavaScript selection

The default selection endpoint accepts an ordinary form post:

```html
<form method="post" action="/dreego-ui/theme">
    <label for="theme">Theme</label>
    <select id="theme" name="theme">
        <option value="white">White</option>
        <option value="black">Black</option>
        <option value="ocean">Ocean</option>
    </select>
    <input type="hidden" name="return_to" value="/settings">
    <button type="submit">Save theme</button>
</form>
```

Unknown themes and external return URLs are rejected. The default preference
store uses an HttpOnly, SameSite=Lax cookie. Applications can provide account-
level storage by implementing `PreferenceStore`.

Applications that need the selected value inside `.dreego` templates can use
the Dreego session store:

```go
if err := app.SetSessionStore(
    dreego.NewCookieStore([]byte("replace-with-a-secret")),
); err != nil {
    log.Fatal(err)
}

_, err := ui.Register(app, ui.Options{
    Preference: ui.SessionPreference{},
})
```

The selected ID is then available as
`c.SessionVal(ui.DefaultThemeSessionKey)`. A custom session key can be supplied
with `ui.SessionPreference{Key: "appearance_theme"}`.

## Validation

Theme registration happens before the Dreego app is built. Registration fails
when a theme contains:

- an empty, reserved, duplicate, or unsafe ID;
- incomplete or non-hex colors;
- text or focus colors below the required contrast ratio;
- unsupported CSS lengths, fonts, shadows, or easing values;
- unordered or sub-millisecond motion durations.

Palette colors are opaque six-digit hex values. Shadow colors may additionally
use eight-digit hex values for alpha transparency. Raw CSS declarations are not
accepted as theme data.

## Black is a theme

`BlackTheme` uses a true `#000000` canvas with raised near-black surfaces. It is
not a dark-mode alias. A preference such as "follow the system" can choose
between registered themes, but it is a selection strategy rather than another
built-in theme.
