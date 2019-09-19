🥘 cafebot
====

Cafebot is a [telegram](https://telegram.org/) bot for making orders in a local cafe. 
Usefull if you don't want to call to make an order or wait inside for your
order to be ready.

### Usage instruction

1. Open [telegram](https://play.google.com/store/apps/details?id=org.telegram.messenger) app
2. Chat with @BotFather -> create your own bot
3. Copy your bot `API Token`
4. Clone this repo
```
$ git clone git@github.com:yarikbratashchuk/cafebot.git && cd cafebot
```
5. Build it from source
```
$ make install 
```
6. Check installation
```
$ ./cafebot --help
```
7. Edit `cafe.json` config file to match your cafe timetable and menu
8. Run app and play with it
```
$ ./cafebot --tg-api-token=YOUR_TELEGRAM_BOT_API_TOKEN
```
