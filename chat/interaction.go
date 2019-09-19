package chat

import (
	"strings"
	"sync"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/config"
	"github.com/yarikbratashchuk/cafebot/order"
)

type (
	// intrHandler is responsible for handling interactions via
	// bot message updates
	intrHandler struct {
		once *sync.Once

		data config.CafeConfig
		bot  *api.BotAPI

		handlers  map[intrEndpoint]intrFunc
		keyboards map[intrEndpoint]keyboardFunc
	}

	intrFunc     func(string, api.Update, order.Order)
	keyboardFunc func(intrData string, o order.Order) *api.InlineKeyboardMarkup

	intrEndpoint string
)

// intr is concurrent safe intrHandler singleton
var intr = &intrHandler{once: &sync.Once{}}

const (
	intrWhere intrEndpoint = "where"
	intrWhen  intrEndpoint = "when"
	intrWhat  intrEndpoint = "what"
)

// initIntrHandler initializes intr singleton
// It must be run once before calling intr.handle
func initIntrHandler(bot *api.BotAPI, cafeconf config.CafeConfig) error {
	log.Info("initializing interaction handler")

	if bot == nil {
		return errNoAPI
	}
	intr.once.Do(func() {
		intr.bot = bot
		intr.data = cafeconf
		intr.handlers = map[intrEndpoint]intrFunc{
			intrWhere: intr.where,
			intrWhen:  intr.when,
			intrWhat:  intr.what,
		}

		// NOTE: since intr.keyboards are called from within intr.handlers
		// all methods must have pointer receiver, otherwise
		// intr.keyboards must be set before intr.handlers
		intr.keyboards = map[intrEndpoint]keyboardFunc{
			intrWhere: whereKeyboardFactory(cafeconf),
			intrWhen:  whenKeyboardFactory(cafeconf),
			intrWhat:  whatKeyboardFactory(cafeconf),
		}

	})
	return nil
}

func (i *intrHandler) handle(
	reqdata string,
	update api.Update,
	order order.Order,
) {
	if i.handlers == nil {
		panic("interaction handler is not initialized")
	}

	parts := strings.Split(reqdata, "/")
	if len(parts) == 0 {
		return
	}
	h, ok := i.handlers[intrEndpoint(parts[0])]
	if !ok {
		return
	}

	var intrData string
	if len(parts) > 1 {
		intrData = parts[1]
	}

	h(intrData, update, order)
}

func (i *intrHandler) where(
	intrData string,
	update api.Update,
	order order.Order,
) {
	i.updateMessage(
		update.CallbackQuery,
		text["where?"],
		intrWhere,
		intrData,
		order,
	)
}

func (i *intrHandler) when(
	intrData string,
	update api.Update,
	order order.Order,
) {
	i.updateMessage(
		update.CallbackQuery,
		text["when?"],
		intrWhen,
		intrData,
		order,
	)
}

func (i *intrHandler) what(
	intrData string,
	update api.Update,
	order order.Order,
) {
	i.updateMessage(
		update.CallbackQuery,
		text["what?"],
		intrWhat,
		intrData,
		order,
	)
}

func (i *intrHandler) updateMessage(
	callbackQuery *api.CallbackQuery,
	msgText string,
	endpoint intrEndpoint,
	intrData string,
	o order.Order,
) {
	// remove spiner at the top right corner of the button
	i.bot.AnswerCallbackQuery(api.NewCallback(callbackQuery.ID, ""))

	msgInfo := callbackQuery.Message

	err := i.updateText(msgInfo, msgText)
	if err != nil {
		log.Errorf("updating text for %+v: %v", msgInfo, err)
		return
	}
	err = i.updateKeyboard(msgInfo, endpoint, intrData, o)
	if err != nil {
		log.Errorf("updating keyboard for %+v: %v", msgInfo, err)
	}
}

// updateText updates text in the last message (from the bot)
func (i *intrHandler) updateText(msgInfo *api.Message, text string) error {
	editText := api.EditMessageTextConfig{
		BaseEdit: api.BaseEdit{
			ChatID:    msgInfo.Chat.ID,
			MessageID: msgInfo.MessageID,
		},
		Text: text,
	}
	_, err := i.bot.Send(editText)
	return err
}

// updateKeyboard updates inline keyboard in the last message (from the bot)
func (i *intrHandler) updateKeyboard(
	msgInfo *api.Message,
	endpoint intrEndpoint,
	intrData string,
	o order.Order,
) error {
	editKeyboard := api.EditMessageReplyMarkupConfig{
		BaseEdit: api.BaseEdit{
			ChatID:      msgInfo.Chat.ID,
			MessageID:   msgInfo.MessageID,
			ReplyMarkup: i.keyboards[endpoint](intrData, o),
		},
	}
	_, err := i.bot.Send(editKeyboard)
	return err
}
