# Dreego UI

Dreego UI is an accessible, SSR-first component and theme library for
[Dreego](https://github.com/dreego-stack/dreego).

The project is under active development. Its theme foundation provides an
extensible typed registry, `WhiteTheme`, `BlackTheme`, custom preference
storage, and a no-JavaScript selection flow. Composable `.dreego` components
will build on this foundation.

## Theme setup

```go
registry, err := ui.Register(app, ui.Options{
    DefaultTheme: "white",
    Themes:       []ui.Theme{myTheme},
})
```

Add `<link rel="stylesheet" href="/dreego-ui/theme.css">` to the application
layout. See [Themes](docs/themes.md) for custom themes and preference storage.

## Design language

Dreego UI is calm, precise, private, and functional. Neuecast, Apple, Proton,
and Signal inform the direction without being copied. Purple is not the default
accent, and `BlackTheme` is designed around a true black canvas rather than
being treated as a generic dark-mode switch.

## Development

```sh
smd make test
smd make check
```

See [the v0.1 specification](docs/spec.md) and
[implementation plan](tasks/plan.md).

## License

MPL-2.0
