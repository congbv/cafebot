package chat

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

func initBotAPI(tgToken string) (*api.BotAPI, error) {
	bot, err := api.NewBotAPI(tgToken)
	if err != nil {
		return nil, err
	}

	//bot.Debug = true
	log.Debugf("Authorized on account %s", bot.Self.UserName)

	return bot, nil
}
