# Product Goal

This page is explains all requirements for the project. Coding agents must plan and implement code to satisfy all needs. 

## Configure

We need cli commands to save Telegram Bot Token and Telegram Chat ID. We should be able to name them and set default. The data should be saved at `~/.config/auto-message`. 

## Send Message

We need cli command to send message via Telegram. The program should use the endpoint under the hood:

```
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage" \
  -d chat_id="$TELEGRAM_CHAT_ID" \
  -d text="Deploy tamamlandı"
```