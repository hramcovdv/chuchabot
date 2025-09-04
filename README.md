# ChuChaBot

Telegram бот для перенаправления сообщений в формате JSON на HTTP-endpoint.

### Как пользоваться:
```
git clone https://github.com/hramcovdv/chuchabot.git
cd chuchabot
```

Далее правим файл **compose.yml**, указываем переменные окружения *TELEGRAM_BOT_API_TOKEN* и *TELEGRAM_BOT_SEND_URL*. И запускаем проект:
```
docker compose up -d
```

`TELEGRAM_BOT_API_TOKEN` - ключ можно получить у [@BotFather](https://t.me/BotFather).

`TELEGRAM_BOT_SEND_URL` - терминальная точка куда будут отправляться JSON-запросы.
