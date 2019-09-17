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
	keyboardFunc func(order.Order) *api.InlineKeyboardMarkup

	intrEndpoint string
)

// intr is concurrent safe intrHandler singleton
var intr = intrHandler{once: &sync.Once{}}

const (
	intrWhere intrEndpoint = "where"
	intrWhen  intrEndpoint = "when"
	intrWhat  intrEndpoint = "what"
)

func initIntrHandler(bot *api.BotAPI, cafeconf config.CafeConfig) error {
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
		intr.keyboards = initKeyboards(cafeconf)
	})
	return nil
}

func (i intrHandler) handle(
	reqdata string,
	update api.Update,
	order order.Order,
) {
	if i.handlers == nil {
		panic("intr is not initialized")
	}

	parts := strings.Split(reqdata, "/")
	if len(parts) == 0 {
		return
	}
	h, ok := i.handlers[intrEndpoint(parts[0])]
	if !ok {
		return
	}

	h(reqdata, update, order)
}

func (i intrHandler) where(
	reqdata string,
	update api.Update,
	order order.Order,
) {
	log.Debug("intr.where")

	msgInfo := update.CallbackQuery.Message
	text := text["where?"]

	err := i.updateText(msgInfo, text)
	if err != nil {
		log.Errorf("update text %#v: %v", msgInfo, err)
		return
	}

	err = i.updateKeyboard(msgInfo, intrWhere, order)
	if err != nil {
		log.Errorf("update keyboard %#v: %v", msgInfo, err)
	}
}

func (i intrHandler) when(
	reqdata string,
	update api.Update,
	order order.Order,
) {
	log.Debug("intr.when")

	msgInfo := update.CallbackQuery.Message
	text := text["when?"]

	err := i.updateText(msgInfo, text)
	if err != nil {
		log.Errorf("update text %#v: %v", msgInfo, err)
		return
	}

	err = i.updateKeyboard(msgInfo, intrWhen, order)
	if err != nil {
		log.Errorf("update keyboard %#v: %v", msgInfo, err)
	}
}

func (i intrHandler) what(
	reqdata string,
	update api.Update,
	order order.Order,
) {
	log.Debug("intr.what")

	// TODO: add stuff here
}

func (i intrHandler) updateText(msgInfo *api.Message, text string) error {
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

func (i intrHandler) updateKeyboard(
	msgInfo *api.Message,
	endpoint intrEndpoint,
	o order.Order,
) error {
	editKeyboard := api.EditMessageReplyMarkupConfig{
		BaseEdit: api.BaseEdit{
			ChatID:      msgInfo.Chat.ID,
			MessageID:   msgInfo.MessageID,
			ReplyMarkup: i.keyboards[endpoint](o),
		},
	}
	_, err := i.bot.Send(editKeyboard)
	return err
}
