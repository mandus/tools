# AGENTS.md

Always follow the [GitHub Spec Kit](https://github.com/github/spec-kit) guidelines when working in this repository.

## Branch Naming

All branches must follow: `<type>/<number>-<short-slug>`

| Prefix | Use Case | Example |
|---|---|---|
| `feat/` | New features | `feat/42-add-git-alias` |
| `fix/` | Bug fixes | `fix/7-prompt-escaping` |
| `docs/` | Documentation | `docs/12-update-readme` |
| `chore/` | Maintenance/tooling | `chore/3-editorconfig` |

- Include the issue/PR number immediately after the prefix.
- Use kebab-case for slugs; keep them short and identifiable.

## Spec-Driven Development

Follow spec-driven development practices: write or consult a spec before implementing. Prefer small, focused changes over large sweeping refactors.

## Testing

Test both the happy path and uninstallation/teardown for any tooling changes. Verify files land in the configured locations after setup.

## Common Pitfalls

- Do not use shorthand keys for CLI tools — use the full executable name.
- Use the correct argument placeholder format for the agent type.
- Always test installation and uninstallation end-to-end.
