# RSS to Telegram

This is simple RSS to Telegram Channel

## How to configure?

1. Copy .env.example to .env
2. Edit the .env file like this:
```
# Example RSS URL https://stackoverflow.blog/newsletter/feed
RSS_URL=""
TELEGRAM_BOT_TOKEN=""
# Use @ at CHANNEL_USERNAME like this @BotFather
CHANNEL_USERNAME=""
```
3. Don't forget to run `go mod tidy`
4. Last run `go run main.go`
5. Happy using ~^-^~