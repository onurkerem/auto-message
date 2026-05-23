# auto-message

Send Telegram bot notifications from terminal commands and automation scripts.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/onurkerem/auto-message/main/install.sh | sh
```

Run the same command again to update to the latest version.

## Quick Start

Add your first bot profile:

```bash
auto-message config add mybot --token "123456:ABC-DEF" --chat-id "999888777" --default
```

Send a message:

```bash
auto-message send "Deploy complete"
```

## Commands

### `auto-message send <message>`

Send a message via Telegram.

```bash
# Uses default profile
auto-message send "Build passed"

# Uses specific profile
auto-message send -c prod "Deployed to production"

# Uses environment variables (no config needed)
AUTO_MESSAGE_TOKEN="123456:ABC-DEF" AUTO_MESSAGE_CHAT_ID="999888777" auto-message send "CI build passed"
```

Resolution order: `--config` flag > default profile > `AUTO_MESSAGE_TOKEN` + `AUTO_MESSAGE_CHAT_ID` env vars.

### `auto-message config add <name>`

Add a named bot configuration profile.

```bash
auto-message config add mybot --token "123456:ABC-DEF" --chat-id "999888777" --default
auto-message config add work --token "456789:XYZ" --chat-id "111222333"
```

### `auto-message config list`

List all saved profiles. Tokens are masked by default.

```bash
auto-message config list
auto-message config list --show-tokens
```

### `auto-message config default <name>`

Set a profile as the default.

```bash
auto-message config default mybot
```

### `auto-message config remove <name>`

Remove a profile.

```bash
auto-message config remove work
```

## CI/CD Usage

For automation without a config file, use environment variables:

```bash
export AUTO_MESSAGE_TOKEN="123456:ABC-DEF"
export AUTO_MESSAGE_CHAT_ID="999888777"
auto-message send "Deploy tamamlandı"
```

## Configuration

Profiles are stored at `~/.config/auto-message/config.json` (respects `XDG_CONFIG_HOME`).

## Uninstall

```bash
rm /usr/local/bin/auto-message
rm -rf ~/.config/auto-message
```

## License

MIT
