# Components

Dreego `v0.6.4` discovers `.dreego` components inside the application's local
website tree. Dreego UI therefore installs component sources explicitly instead
of relying on an implicit module namespace. This version also uses the current
[`<server>` and `<body>` semantic sections](https://github.com/dreego-stack/dreego/blob/v0.6.4/_docs/semantic-sections-migration.md).

## Install

Run the installer from the Dreego application root:

```sh
go run github.com/dreego-stack/dreego-ui/cmd/dreego-ui@latest install
```

The default destination is `www/components/dreegoui`. Existing files are never
overwritten. Use a custom website root when needed:

```sh
go run github.com/dreego-stack/dreego-ui/cmd/dreego-ui@latest install \
  -destination web/components/dreegoui
```

The `dreegoui` directory name is intentional. It gives generated Go code a
stable `dreegoui` package name under Dreego's current component importer.

## Button

```dreego
import Button "components/dreegoui/Button.dreego"

<body>
    <@Button disabled={false}>Save changes</@Button>
    <@Button disabled={false} variant="secondary">Cancel</@Button>
</body>
```

Variants are `primary`, `secondary`, `quiet`, and `danger`. Unknown variants
fall back to `primary`. `disabled` is required because Dreego v0.6.4 does not
support default values for Boolean component props.

## Toggle

```dreego
import Toggle "components/dreegoui/Toggle.dreego"

<body>
    <@Toggle
        id="product-updates"
        name="product_updates"
        label="Product updates"
        checked={true}
        disabled={false}
    />
</body>
```

Toggle uses a native checkbox associated with a visible label and exposes
switch semantics to assistive technology.

## Card

```dreego
import Card "components/dreegoui/Card.dreego"

<body>
    <@Card>
        {#slot header}<h2>Account</h2>{/slot}
        <p>Private by default.</p>
        {#slot footer}<a href="/settings">Settings</a>{/slot}
    </@Card>
</body>
```

Card uses named header and footer slots plus its default body slot.

## CodeBox

```dreego
import CodeBox "components/dreegoui/CodeBox.dreego"

<body>
    <@CodeBox code="go test ./..." language="shell" label="Test command"/>
</body>
```

Code content and labels are escaped by Dreego. The scrollable code region is
keyboard focusable.

## ThemePicker

`ThemeOption` is installed beside the component and belongs to the generated
`dreegoui` package.

```dreego
import ThemePicker "components/dreegoui/ThemePicker.dreego"

<server>
    themes := []dreegoui.ThemeOption{
        {ID: "white", Name: "White"},
        {ID: "black", Name: "Black"},
        {ID: "ocean", Name: "Ocean"},
    }
</server>

<body>
    <@ThemePicker
        id="theme-choice"
        themes={themes}
        csrfToken={c.CSRFToken()}
        selected="white"
        returnTo="/settings"
    />
</body>
```

ThemePicker submits a normal form to the plugin's selection endpoint. The CSRF
token is required by the component API so applications do not accidentally
omit Dreego's default protection. Pass an empty string only when CSRF is
explicitly disabled for the application.

## Dreego v0.6.4 compatibility

Dreego v0.6.4 emits scoped CSS attribute selectors without quotes. A numeric
scope hash is rejected by browsers, so the distributed sources contain stable
`data-dreego-ui` seed attributes and tests guard their hashes. Keep those
attributes when adapting a component.

The v0.6.4 formatter also inserts an additional blank line into component
declaration files on every run. Generation and compilation remain stable, but
the library intentionally does not use `dreego fmt --check` until that upstream
formatter behavior is fixed.
