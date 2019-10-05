Work in progress

Telegram bot written in go

The goal is a bot which sends new notifications as telegram messages

To get the bot working you need to have an `.env` file with `TELEGRAM_API_TOKEN` and `GITHUB_API_TOKEN` set

TODO
- error handling
- check on start if all env are set
- check for chat token so that only the correct person gets the notifications
- chat token could be set in .env or passed in via command line 
- track which updates are already sent
- send message on new updates
- define struct and maybe send more information like repo name
- dockerize
