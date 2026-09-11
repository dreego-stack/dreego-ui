# Spec: Dreego UI v0.1

## Assumptions

1. Dreego UI is a web component library for Dreego SSR applications, not a
   native iOS library.
2. The Go module path is `github.com/dreego-stack/dreego-ui`.
3. Dreego UI remains an external module and Dreego Core never imports it.
4. JavaScript is optional. Components must render useful, accessible HTML
   without a client runtime.
5. Applications own persistence of account-level preferences. The plugin may
   offer cookie-based theme selection as an opt-in convenience.
6. Neuecast, Apple, Proton, and Signal are visual references only.

## Objective

Build an accessible, SSR-first UI library for Dreego applications with a
recognizable design language and a typed, extensible theme system.

The library ships `WhiteTheme` and `BlackTheme`. Applications can register any
number of additional themes, choose their default theme, render a theme picker,
and persist the selected theme without modifying Dreego UI internals.

## Product Principles

- Calm: content has priority over decoration.
- Precise: spacing, typography, states, and motion have clear purpose.
- Private: interfaces communicate trust without security theater.
- Functional: native semantics and progressive enhancement come first.
- Adaptable: applications control color and may define many themes.
- Accessible: information never depends on color, vision, or pointer input.

`BlackTheme` is an intentional black interface with a `#000000` canvas and
raised near-black surfaces. It is not named or modeled as generic dark mode.
Purple is not used as the default accent color.

## Tech Stack

- Go 1.27 or newer
- `github.com/dreego-stack/dreego/core`
- Dreego `.dreego` components
- CSS custom properties
- Native HTML and optional minimal plain JavaScript
- Go standard library tests
- `smd` for all development commands

No frontend framework, Node.js runtime, CSS framework, or icon dependency is
required by the library.

## Public Theme Contract

Themes are data, so the primary extension contract is a struct rather than an
interface. Interfaces are reserved for behavior such as preference storage.

```go
type Theme struct {
    ID          string
    Name        string
    Description string
    Colors      Colors
    Typography  Typography
    Shape       Shape
    Elevation   Elevation
    Motion      Motion
}

type Options struct {
    DefaultTheme string
    Themes       []Theme
    Preference   PreferenceStore
}

func Register(app *dreego.App, options Options) (*Registry, error)
```

The registry:

- contains immutable copies of registered themes;
- rejects empty, reserved, or duplicate IDs;
- rejects an unknown default theme;
- validates every value before CSS generation;
- returns themes in registration order for stable settings UIs;
- provides lookup and stylesheet generation without exposing mutable state.

`WhiteTheme()` and `BlackTheme()` return fresh values that callers may copy and
adapt. Built-in themes are registered by default unless explicitly disabled.

## Theme Selection

Rendered pages select a theme using a stable HTML attribute:

```html
<html data-dreego-theme="black">
```

The stylesheet contains selectors for every registered theme and component
styles consume namespaced custom properties such as `--dreego-color-canvas`.

The resolution order is:

1. valid explicit theme selected by the application;
2. valid stored preference when a preference store is configured;
3. configured default theme;
4. `white`.

An invalid or removed stored theme falls back safely and never emits arbitrary
content into HTML or CSS. A system-color preference is a selection strategy,
not a third built-in theme.

## Component Scope

### Foundations

- Theme registry and CSS tokens
- Typography, spacing, shape, elevation, and motion
- Focus, disabled, invalid, loading, and reduced-motion states
- Small, replaceable icon contract and essential icon set

### Components

- Button and IconButton
- Input, Select, Checkbox, Radio, Form, and Toggle
- Badge, Alert, Toast, and EmptyState
- Card, PriceCard, CodeBox, and Dialog
- Navbar, Sidebar, Tabs, and Footer
- PageShell, Container, Stack, and Grid
- Table
- ThemePicker

Components favor composition and slots over large configuration structs.
Interactive controls use native elements whenever possible.

## Distribution

The Go module owns theme registration, validation, preference helpers, CSS
generation, and embedded assets. Dreego component sources are distributed from
the module with documented explicit import paths or a deterministic installer,
depending on the import behavior validated against the current Dreego CLI.

Generated application files are never silently overwritten. Installation and
updates must be explicit and reproducible.

## Project Structure

```text
dreego-ui/
├── docs/                 design, API, accessibility, and usage documentation
├── tasks/                implementation plan and task state
├── components/           distributable .dreego components
├── icons/                essential replaceable icons
├── internal/             validation and CSS implementation
├── example/              visual component and theme showcase
├── _tests/               end-to-end Dreego integration tests
├── theme.go              public theme types and built-ins
├── registry.go           public registry API
├── plugin.go             explicit App registration
├── go.mod
└── smd.toml
```

## Commands

The exact commands are finalized with the initial scaffold. The required
interface is:

```sh
smd make test
smd make check
smd make example
```

## Testing Strategy

- Unit tests cover theme validation, duplicate handling, fallback resolution,
  deterministic CSS, immutability, and preference behavior.
- Golden tests cover semantic component markup and generated styles.
- Integration tests build a real Dreego fixture with custom themes.
- Browser checks cover 320, 768, 1024, and 1440 pixel widths.
- Accessibility checks cover keyboard order, visible focus, labels, semantics,
  non-color states, reduced motion, and WCAG AA contrast.
- Tests for behavioral work are written failing before implementation.

## Boundaries

### Always

- Keep Dreego UI outside Dreego Core.
- Keep SSR output useful without JavaScript.
- Escape HTML content and validate CSS values at system boundaries.
- Preserve application ownership of account settings and persistence.
- Support an arbitrary number of application-defined themes.

### Ask First

- Add a third-party dependency.
- Add a client runtime or mandatory JavaScript.
- Add network-backed theme storage.
- Expand the stable public API beyond the approved v0.1 contract.
- Publish a repository, tag, package, or release.

### Never

- Import Dreego UI from Dreego Core.
- Call `BlackTheme` dark mode in the public API or documentation.
- Use purple as the default accent.
- Accept raw unvalidated CSS declarations from theme data.
- Make security or accessibility claims that the library cannot guarantee for
  arbitrary application content.
- Depend on color alone to communicate state.

## Success Criteria

- `WhiteTheme` and `BlackTheme` pass token and contrast validation.
- At least three additional custom themes can be registered and rendered in one
  application without code changes to Dreego UI.
- Duplicate and malformed themes fail with specific errors before the app is
  built.
- Theme output is deterministic and safe for untrusted theme configuration.
- The settings example switches among all registered themes and has a no-JS
  server fallback.
- Every listed component has documented semantics and a tested example.
- The full test and check commands pass in `smd`.
- No handwritten file exceeds 300 lines.

## Open Questions

- Confirm whether Dreego's current component importer can consume `.dreego`
  files directly from a Go module. If not, use an explicit, non-overwriting
  installer for v0.1.
- Confirm the final neutral/accent palette visually in the showcase before
  treating built-in theme colors as stable.
- Decide whether cookie preference storage belongs in v0.1 or follows after the
  registry and ThemePicker vertical slice.
