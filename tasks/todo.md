# Dreego UI v0.1 Tasks

- [x] Create the repository scaffold and development commands.
  - Acceptance: module metadata, license, README, Makefile, and `smd.toml` exist.
  - Verify: `smd make check` runs from the repository root.
  - Files: root configuration and documentation.

- [x] Implement theme types and built-in themes using TDD.
  - Acceptance: fresh valid `WhiteTheme` and `BlackTheme` values are returned.
  - Verify: focused Go unit tests pass.
  - Files: `theme.go`, `theme_test.go`, internal token files.

- [x] Implement theme validation using TDD.
  - Acceptance: malformed IDs, colors, lengths, and unsafe CSS values fail with
    actionable errors.
  - Verify: validation unit and fuzz tests pass.
  - Files: internal validation package and tests.

- [x] Implement the immutable registry and CSS generator using TDD.
  - Acceptance: arbitrary themes register in stable order; duplicates fail;
    generated CSS is deterministic.
  - Verify: registry tests and stylesheet golden tests pass.
  - Files: `registry.go`, registry tests, internal CSS package and tests.

- [x] Implement App registration and theme resolution using TDD.
  - Acceptance: assets register before build; late or conflicting registration
    returns wrapped Dreego errors; invalid selections fall back safely.
  - Verify: unit and Dreego integration tests pass.
  - Files: `plugin.go`, selection files, tests, one integration fixture.

- [ ] Validate and implement component distribution.
  - Acceptance: application imports are explicit and installation never
    overwrites user files silently.
  - Verify: a clean fixture generates and builds twice reproducibly.
  - Files: installer or import metadata, fixture, documentation.

- [ ] Build the first component vertical slice.
  - Acceptance: Button, Toggle, Card, CodeBox, ThemePicker, essential icons, and
    layout foundations render accessibly in all registered themes.
  - Verify: markup tests, Dreego build, and responsive browser review pass.
  - Files: component sources, focused tests, showcase pages, documentation.

- [ ] Add remaining component families incrementally.
  - Acceptance: every v0.1 component has semantic markup, docs, and showcase
    coverage.
  - Verify: focused tests pass after each family.
  - Files: no task slice changes more than five files.

- [ ] Complete the release quality review.
  - Acceptance: API, accessibility, security, responsive behavior, docs, and
    file-size rules have no known release blocker.
  - Verify: full `smd make check` and browser checklist pass.
  - Files: tests, docs, and release metadata only.
