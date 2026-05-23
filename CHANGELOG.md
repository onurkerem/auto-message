# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-05-23

### Added

- `auto-message config add <name> --token <token> --chat-id <id> [--default]` to add a named bot profile
- `auto-message config list [--show-tokens]` to list all profiles (tokens masked by default)
- `auto-message config default <name>` to set the default profile
- `auto-message config remove <name>` to remove a profile
- `auto-message send <message> [--config <name>]` to send a Telegram message
- Environment variable fallback: `AUTO_MESSAGE_TOKEN` and `AUTO_MESSAGE_CHAT_ID`
- Configuration stored at `~/.config/auto-message/config.json` (respects `XDG_CONFIG_HOME`)
- Cross-platform binaries for macOS, Linux, and Windows (amd64, arm64)
- One-command install/update via `curl -fsSL https://raw.githubusercontent.com/onurkerem/auto-message/main/install.sh | sh`
