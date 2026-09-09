# Dreego UI v0.1 Implementation Plan

## Dependency Order

```text
repository scaffold
        |
theme model and validation
        |
registry and deterministic CSS
        |
application registration and selection
        |
component foundations and icons
        |
component families
        |
showcase and browser accessibility review
        |
release readiness review
```

## Phase 1: Foundation

Create the Go module, development commands, CI-ready checks, and the smallest
theme model. Prove `WhiteTheme`, `BlackTheme`, custom themes, validation, and
deterministic CSS before component work begins.

Checkpoint: unit tests and build pass in `smd`; the public theme API is reviewed.

## Phase 2: Dreego Integration

Register theme assets against an owning `dreego.App`, resolve a selected theme,
and validate the component distribution path against the real Dreego CLI.

Checkpoint: a Dreego integration fixture renders three or more themes without a
client runtime.

## Phase 3: Component Foundations

Implement typography, layout primitives, states, essential icons, Button,
Toggle, Card, CodeBox, and ThemePicker as the first vertical component slice.

Checkpoint: responsive and keyboard checks pass in the showcase.

## Phase 4: Component Families

Add form, feedback, navigation, surface, layout, and data components in small
independent slices. Each slice includes tests and documentation.

Checkpoint: every component in the v0.1 specification is represented in the
showcase and has semantic markup tests.

## Phase 5: Quality Gate

Review API consistency, accessibility, CSS safety, responsive behavior,
documentation, file size, and release metadata.

Checkpoint: `smd make check` and the full browser review pass with no known
release-blocking issue.

## Risks and Mitigations

- Component module imports may not yet be implemented by Dreego. Validate this
  early and use an explicit non-overwriting installer if required.
- User-defined CSS tokens can become an injection boundary. Permit only typed,
  grammar-validated token values.
- A large first component catalog can hide weak foundations. Ship vertical
  slices and review the showcase between component families.
- Theme flexibility can reduce accessibility. Validate built-ins strictly and
  provide contrast diagnostics for custom themes without making false claims.

