# Telegram Bot Setup Guide

[中文版本](./telegram_bot_setup.md)

### 1. Create a Telegram Bot

1. First, add [@BotFather](https://t.me/BotFather) in Telegram
2. Send the `/newbot` command to BotFather
3. Follow the prompts to set up your bot's name and username
4. Upon completion, BotFather will provide an API Token (format: `123456789:ABCdefGHIjklmNOPQrstUVwxyz`)
5. Save this Token securely as it will be needed for configuration

### 2. Get the Group Chat ID

1. Add your newly created bot to your target Telegram group
2. Send a message in the group (this step is crucial for obtaining the Chat ID)
3. Open the following URL in your browser (replace `YOUR_BOT_TOKEN` with your Bot API Token):
   ```
   https://api.telegram.org/botYOUR_BOT_TOKEN/getUpdates
   ```
4. In the returned JSON data, locate the `id` field within the `chat` object - this is your group's Chat ID
   - Group Chat IDs are typically negative numbers, e.g., `-1001234567890`

### 3. Configuration Instructions

In your configuration file, you'll need to provide:
- `BOT_TOKEN`: The API Token received from BotFather
- `CHAT_ID`: The group ID obtained from the getUpdates endpoint

### Important Notes

- Please keep your Bot Token secure and never share it publicly
- If you receive an empty response when getting the Chat ID, ensure you've sent a message in the group
- It's recommended to test the configuration in a test group before applying it to the production group 