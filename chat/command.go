package chat

import (
	"errors"
	"sync"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

// cmdHandler is responsible for handling chat commands
type cmdHandler struct {
	once sync.Once

	bot      *api.BotAPI
	commands map[string]func(api.Update)
}

var (
	// cmd is concurrent safe cmdHandler singleton
	cmd cmdHandler

	errNoAPI = errors.New("bot api is nil")
)

const (
	cmdStartEndpoint = "start"
	cmdHelpEndpoint  = "help"
)

func initCmdHandler(bot *api.BotAPI) error {
	if bot == nil {
		return errNoAPI
	}
	cmd.once.Do(func() {
		cmd.bot = bot
		cmd.commands = map[string]func(api.Update){
			cmdStartEndpoint: cmd.start,
			cmdHelpEndpoint:  cmd.help,
		}
	})
	return nil
}

func (c cmdHandler) handle(command string, update api.Update) {
	if c.commands == nil {
		panic("cmd is not initialized")
	}

	h, ok := c.commands[command]
	if !ok {
		c.wrong(update)
		return
	}

	h(update)
}

func (c cmdHandler) help(update api.Update) {
	log.Debug("cmd.help")

	msg := api.NewMessage(
		update.Message.Chat.ID,
		message["help"],
	)

	c.bot.Send(msg)
}

func (c cmdHandler) start(update api.Update) {
	log.Debug("cmd.start")

	msg := api.NewMessage(
		update.Message.Chat.ID,
		message["start"],
	)
	msg.ReplyMarkup = api.NewInlineKeyboardMarkup(
		api.NewInlineKeyboardRow(
			api.NewInlineKeyboardButtonData(
				buttonText["new_order"],
				intrWhereEndpoint,
			),
		),
	)

	c.bot.Send(msg)
}

func (c cmdHandler) wrong(update api.Update) {
	log.Debug("cmd.wrong")

	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	}

	msg := api.NewMessage(
		chatID,
		message["wrong"],
	)

	c.bot.Send(msg)
}
