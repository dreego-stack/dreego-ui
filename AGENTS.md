# Agent Instructions for Dreego UI

## Language

- User communication is German.
- Repository content is English.

## Project

- Module: `github.com/dreego-stack/dreego-ui`
- Dreego UI is an external Dreego plugin and component library.
- Dreego Core must never import this module.
- Keep the package SSR-first and usable without a client runtime.
- Prefer the Go standard library and native HTML elements.
- Do not add dependencies without an explicit supply-chain review.

## Development

- Run development commands through `smd`.
- Run Git commands on the host.
- Do not create binaries in the repository; use `/tmp` or `./tmp`.
- Use tests before implementation for behavioral changes.
- Keep handwritten files under 300 lines.
- Keep public APIs typed and documented.
- Validate user-provided theme values before generating CSS.
- Accessibility is a release gate: keyboard access, visible focus, semantic HTML,
  non-color state communication, reduced motion, and WCAG AA contrast are required.

## Design

- The visual language is calm, precise, private, and functional.
- Neuecast, Apple, Proton, and Signal are references, not templates to copy.
- Purple is not a default brand color.
- `BlackTheme` is a deliberate black design, not a generic dark-mode alias.
- Applications may register any number of custom themes.

## Git

- Use `codex/` branches for implementation.
- Keep commits small and independently testable.
- Do not create tags locally.

