# Telegram Bot 配置指南

[English Version](./telegram_bot_setup_en.md)

### 1. 创建 Telegram 机器人

1. 首先在 Telegram 中添加 [@BotFather](https://t.me/BotFather)
2. 向 BotFather 发送 `/newbot` 命令
3. 按照提示设置机器人名称和用户名
4. 完成后，BotFather 会提供一个 API Token（格式如：`123456789:ABCdefGHIjklmNOPQrstUVwxyz`）
5. 请妥善保存此 Token，它将用于后续配置

### 2. 获取群组 Chat ID

1. 将创建好的机器人添加到目标 Telegram 群组中
2. 在群组中发送一条消息（这一步很重要，否则可能获取不到 Chat ID）
3. 在浏览器中打开以下 URL（替换 `YOUR_BOT_TOKEN` 为你的 Bot API Token）：
   ```
   https://api.telegram.org/botYOUR_BOT_TOKEN/getUpdates
   ```
4. 在返回的 JSON 数据中，找到 `chat` 对象中的 `id` 字段，这就是群组的 Chat ID
   - 群组的 Chat ID 通常是一个负数，例如：`-1001234567890`

### 3. 配置说明

在配置文件中，你需要填入：
- `BOT_TOKEN`：从 BotFather 获取的 API Token
- `CHAT_ID`：从 getUpdates 接口获取的群组 ID

### 注意事项

- 请确保妥善保管你的 Bot Token，不要分享给他人或暴露在公开场合
- 如果在获取 Chat ID 时遇到空白响应，请确保已在群组中发送过消息
- 建议在测试群组中进行配置测试，确认无误后再应用到正式群组 