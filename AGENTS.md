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

## Commit Messages

Every commit message must start with a [gitmoji](https://gitmoji.dev/) that
matches the nature of the change, followed by a concise subject line written
in the imperative mood.

**Format:**
```
<emoji> <subject>

<body — optional, explains why not what>
```

**Example:**
```
✨ Add gitprompt binary for fast shell prompt

Reads .git/ internals directly to eliminate git subprocess overhead
on Windows/WSL where process creation costs 50–150 ms per prompt.
```

**Common gitmojis for this repo:**

| Emoji | When to use |
|-------|-------------|
| ✨ | Introduce a new feature |
| 🐛 | Fix a bug |
| 🩹 | Simple fix for a non-critical issue |
| 🚑️ | Critical hotfix |
| ♻️ | Refactor code without changing behaviour |
| ⚡️ | Improve performance |
| 🔥 | Remove code or files |
| 📝 | Add or update documentation |
| ✅ | Add, update, or pass tests |
| ⬆️ | Upgrade dependencies |
| ⬇️ | Downgrade dependencies |
| 📌 | Pin dependencies to specific versions |
| 🔧 | Add or update configuration files |
| 🔨 | Add or update build/dev scripts |
| 🏗️ | Make architectural changes |
| 💥 | Introduce breaking changes |
| ⏪️ | Revert changes |
| 🙈 | Add or update .gitignore |
| 👷 | Add or update CI/CD pipeline |
| 💚 | Fix CI build |
| 🚚 | Move or rename files/paths |
| 🎉 | Begin a project |

Full reference: <https://gitmoji.dev/>

## Spec-Driven Development

Follow spec-driven development practices: write or consult a spec before implementing. Prefer small, focused changes over large sweeping refactors.

## Testing

Test both the happy path and uninstallation/teardown for any tooling changes. Verify files land in the configured locations after setup.

## Common Pitfalls

- Do not use shorthand keys for CLI tools — use the full executable name.
- Use the correct argument placeholder format for the agent type.
- Always test installation and uninstallation end-to-end.
