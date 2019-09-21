package chat

import (
	"errors"
	"sync"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

type (
	// cmdHandler is responsible for handling chat commands
	cmdHandler struct {
		once *sync.Once

		bot      *api.BotAPI
		commands map[string]cmdFunc
	}

	cmdFunc func(api.Update) error
)

var (
	// cmd is concurrent safe cmdHandler singleton
	cmd = &cmdHandler{once: &sync.Once{}}

	errNoAPI = errors.New("bot api is nil")
)

const (
	cmdStartEndpoint = "start"
	cmdHelpEndpoint  = "help"
)

// initCmdHandler initializes cmd singleton
// It must be run once before calling cmd.handle
func initCmdHandler(bot *api.BotAPI) error {
	log.Info("initializing cmd handler")

	if bot == nil {
		return errNoAPI
	}
	cmd.once.Do(func() {
		cmd.bot = bot
		cmd.commands = map[string]cmdFunc{
			cmdStartEndpoint: cmd.start,
			cmdHelpEndpoint:  cmd.help,
		}
	})
	return nil
}

func (c *cmdHandler) handle(command string, update api.Update) error {
	if c.commands == nil {
		panic("command handler is not initialized")
	}

	h, ok := c.commands[command]
	if !ok {
		h = c.wrong
	}

	return h(update)
}

func (c *cmdHandler) help(update api.Update) error {
	_, err := c.bot.Send(
		newMessage(
			update.Message.Chat.ID,
			text["help"],
			false,
		),
	)
	return err
}

func (c *cmdHandler) start(update api.Update) error {
	msg := newMessage(update.Message.Chat.ID, text["start"], true)
	msg.ReplyMarkup = api.NewInlineKeyboardMarkup(
		api.NewInlineKeyboardRow(
			api.NewInlineKeyboardButtonData(
				buttonText["new"],
				newIntrData(intrWhere, "", opInit, "1"),
			),
		),
	)

	_, err := c.bot.Send(msg)
	return err
}

func (c *cmdHandler) wrong(update api.Update) error {
	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	}

	_, err := c.bot.Send(newMessage(chatID, text["wrong"], false))
	return err
}

func newMessage(chatID int64, text string, bold bool) *api.MessageConfig {
	if bold {
		text = boldText(text)
	}
	msg := api.NewMessage(chatID, text)
	msg.ParseMode = api.ModeHTML
	return &msg
}
