# AGENTS.md

## Project

This is a public Go CLI project for sending Telegram bot notifications from terminal commands and automation scripts.

The tool must be simple, stable, and easy to install or update with a single command.

## Repository Structure

```
packages/
  cli/        — Go CLI source, go.mod, goreleaser config
  website/    — Astro marketing site (Tailwind 4, TypeScript)
```

## Tech Stack

- CLI: Go, Cobra
- Website: Astro 6, Tailwind 4, TypeScript
- Config: environment variables and optional config file
- Distribution: GitHub Releases
- Versioning: Semantic Versioning

## Product Requirements

Read when needed: @docs/product-goal.md

## Installation Requirement

Users must be able to install or update the tool with one command.

Preferred installation pattern:

```bash
curl -fsSL https://raw.githubusercontent.com/onurkerem/auto-message/main/install.sh | sh
```

If the tool is already installed, the same command must update it to the latest stable release.

## Release Requirements

- Use Git tags like `v1.0.0`
- Build binaries for macOS, Linux, and Windows
- Publish binaries through GitHub Releases
- Keep `CHANGELOG.md` updated
- Do not break existing CLI commands without a major version bump

## Development Rules

- Keep the code small and readable
- Prefer standard library unless a dependency is clearly useful
- Do not hardcode secrets
- Do not commit tokens, chat IDs, or local config files
- Add tests for command parsing and Telegram API client behavior
- Mock external HTTP calls in tests

## CLI UX Rules

Commands should be predictable and friendly to agents powered by a LLM.

Error messages and `help` command should be explanatory.

## Security

- Never log bot tokens
- Mask sensitive values in debug output
- Use HTTPS only
- Validate required config before making API calls

## Agent Guidance

When changing this project:

- Preserve backward compatibility
- Update README examples when CLI behavior changes
- Update CHANGELOG for user-facing changes
- Prefer small, focused pull requests
- Do not introduce unnecessary frameworks
- Website changes are part of every feature — any change to CLI flags, config structure, or documented behavior must be reflected in the website
