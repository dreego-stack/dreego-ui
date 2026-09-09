# Change files

Every pull request adds exactly one Markdown file in this directory.

```markdown
---
version: patch
---

- Feat: add X
```

Use `version: none` only when the release tag must not change. The release
workflow combines pending files into the changelog and removes them after a
release.
