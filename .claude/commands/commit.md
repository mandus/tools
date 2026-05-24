Create a git commit for the staged changes (or stage and commit all modified tracked files if nothing is staged). Follow this process:

1. Run `git status` and `git diff --staged` (or `git diff HEAD` if nothing is staged) to understand what is changing.
2. Choose the single most appropriate gitmoji from the reference below.
3. Write a commit message:
   - First line: `<emoji> <subject>` — imperative mood, ≤72 chars, no period.
   - Blank line, then an optional body explaining *why* (not what).
4. Commit using the message.

If `$ARGUMENTS` is provided, treat it as additional context or a draft subject line.

---

## Gitmoji reference

| Emoji | Code | When to use |
|-------|------|-------------|
| ✨ | `:sparkles:` | Introduce a new feature |
| 🐛 | `:bug:` | Fix a bug |
| 🩹 | `:adhesive_bandage:` | Simple fix for a non-critical issue |
| 🚑️ | `:ambulance:` | Critical hotfix |
| ♻️ | `:recycle:` | Refactor without changing behaviour |
| ⚡️ | `:zap:` | Improve performance |
| 🔥 | `:fire:` | Remove code or files |
| 📝 | `:memo:` | Add or update documentation |
| ✅ | `:white_check_mark:` | Add, update, or pass tests |
| ⬆️ | `:arrow_up:` | Upgrade dependencies |
| ⬇️ | `:arrow_down:` | Downgrade dependencies |
| 📌 | `:pushpin:` | Pin dependencies to specific versions |
| 🔧 | `:wrench:` | Add or update configuration files |
| 🔨 | `:hammer:` | Add or update build/dev scripts |
| 🏗️ | `:building_construction:` | Make architectural changes |
| 💥 | `:boom:` | Introduce breaking changes |
| ⏪️ | `:rewind:` | Revert changes |
| 🙈 | `:see_no_evil:` | Add or update .gitignore |
| 👷 | `:construction_worker:` | Add or update CI/CD pipeline |
| 💚 | `:green_heart:` | Fix CI build |
| 🚚 | `:truck:` | Move or rename files/paths |
| 🎉 | `:tada:` | Begin a project |

Full reference: https://gitmoji.dev/
