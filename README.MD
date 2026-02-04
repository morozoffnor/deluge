# Телеграм бот для загрузки торрентов через Deluge

## Настройка

### бот
1. создать бота в BotFather
2. создать файл `config.json`
3. заполнить файл `config.json` данными
```json
{
  "telegram_token": "<token>",
  "telegram_chat_id": <chat_id>, // без кавычек, должен быть тип int64
  "films_path": "<путь до папки с торрент-файлами фильмов>", // не забудьте / в конце
  "shows_path": "<путь до папки с торрент-файлами сериалов>", // не забудьте / в конце
  "api_port": "80" // если это менять, то нужно поменять и в скриптах
} 
```
4. добавить бота как сервис в систему
   1. `touch /etc/systemd/system/torrent_bot.service`
   2. `nano /etc/systemd/system/torrent_bot.service`
   3. вставить в файл текст:
   ```
    [Unit]
    Description=Telegram Bot for Downloading torrents
    After=network.target
    
    [Service]
    Type=simple
    ExecStart=<абсолютный путь до исполняемого файла бота>
    Restart=always
    Environment="CONFIG_FILE=<абсолютный путь до файла конфигурации>"
   
    [Install]
    WantedBy=multi-user.target
    ```
5. `sudo chmod +x <путь до исполняемого файла бота>`
6. `systemctl daemon-reload`
7. `systemctl enable torrent_bot.service`
8. `systemctl start torrent_bot.service`

### Deluge
1. включить плагины AutoAdd и Execute (перезапустить Deluge)
2. для AutoAdd добавить две папки (фильмы и сериалы) откуда будут браться торренты
3. в плагине Execute добавить скрипты `torrent_added.sh` и `torrent_completed.sh`
из репозитория, указать события сооответственно