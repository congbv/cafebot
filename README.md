# 🥘 Cafebot
====

Cafebot is a [telegram](https://telegram.org/) bot for making orders in a local cafe. 
Its handy if you don't want to make an order my the call or sit and wait for your
order to be ready.

## Usage instruction

1. Open [telegram](https://play.google.com/store/apps/details?id=org.telegram.messenger) app
2. Chat with @BotFather -> create your own bot
3. Copy your bot `API Token`
4. Clone this repo
```
$ git clone git@github.com:yarikbratashchuk/cafebot.git && cd cafebot
```
5. Build it from source (requires go toolchain to be installed)
```
$ make install 
```
6. Check installation
```
$ cafebot --help
Usage:
  cafebot [OPTIONS]

  Application Options:
    -p, --port=     port to listen on (default: 8080)
        --tg-api-token= telegram bot api token [$TG_API_TOKEN]
        --log-level=    log level for all subsystems {trace, debug, info, error, critical} (default: info)
        --cafe-config=  cafe config file path (default: cafe.json)

  Help Options:
    -h, --help      Show this help message

```
7. Edit `cafe.json` config file to match your cafe timetable and menu
8. Run app and play with it
```
$ cafebot --tg-api-token=YOUR_TELEGRAM_BOT_API_TOKEN
```
