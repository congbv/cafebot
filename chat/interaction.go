package chat

import (
	"fmt"
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
	keyboardFunc func(intrData string, o order.Order) api.InlineKeyboardMarkup

	intrEndpoint string
)

// intr is concurrent safe intrHandler singleton
var intr = &intrHandler{once: &sync.Once{}}

const (
	intrWhere   intrEndpoint = "where"
	intrWhen    intrEndpoint = "when"
	intrWhat    intrEndpoint = "what"
	intrPreview intrEndpoint = "preview"
	intrFinish  intrEndpoint = "finish"
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
			intrWhere:   intr.where,
			intrWhen:    intr.when,
			intrWhat:    intr.what,
			intrPreview: intr.preview,
		}

		// NOTE: since intr.keyboards are called from within intr.handlers
		// all methods must have pointer receiver, otherwise
		// intr.keyboards must be set before intr.handlers
		intr.keyboards = map[intrEndpoint]keyboardFunc{
			intrWhere:   whereKeyboardFactory(cafeconf),
			intrWhen:    whenKeyboardFactory(cafeconf),
			intrWhat:    whatKeyboardFactory(cafeconf),
			intrPreview: previewKeyboardFactory(cafeconf),
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

	endpoint, data := splitEndpointIntrData(reqdata)

	h, ok := i.handlers[endpoint]
	if !ok {
		return
	}

	h(data, update, order)
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

func (i *intrHandler) preview(
	intrData string,
	update api.Update,
	order order.Order,
) {
	caption := text["preview"]
	p := generatePreviewText(order)

	i.updateMessage(
		update.CallbackQuery,
		fmt.Sprintf("%s\n\n%s", caption, p),
		intrPreview,
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
	msgInfo := callbackQuery.Message

	newText := i.prepareUpdateText(msgInfo, msgText)
	newKeyboard := i.prepareUpdateKeyboard(msgInfo, endpoint, intrData, o)

	var err error
	for _, f := range []func(){
		func() { _, err = i.bot.Send(newText) },
		func() { _, err = i.bot.Send(newKeyboard) },
	} {
		f()
		if err != nil {
			log.Errorf("sending update: %+v", err)
		}
	}

	// remove spiner at the top right corner of the button
	_, _ = i.bot.AnswerCallbackQuery(api.NewCallback(callbackQuery.ID, ""))
}

// prepareUpdateText prepares text update for the last message from the bot
func (i *intrHandler) prepareUpdateText(
	msgInfo *api.Message,
	text string,
) api.EditMessageTextConfig {
	return api.EditMessageTextConfig{
		BaseEdit: api.BaseEdit{
			ChatID:    msgInfo.Chat.ID,
			MessageID: msgInfo.MessageID,
		},
		Text: text,
	}
}

// prepareUpdateKeyboard prepares new keyboard for the last message from the bot
func (i *intrHandler) prepareUpdateKeyboard(
	msgInfo *api.Message,
	endpoint intrEndpoint,
	intrData string,
	o order.Order,
) api.EditMessageReplyMarkupConfig {
	replyMarkup := i.keyboards[endpoint](intrData, o)

	// TODO (yb): refactor this logic somehow
	var lastButtonRow []api.InlineKeyboardButton
	var backButton *api.InlineKeyboardButton
	if endpoint == intrWhen {
		backButton = backKeyboardButton(intrWhere)
	} else if endpoint == intrWhat {
		if intrData == "" {
			backButton = backKeyboardButton(intrWhen)
		} else {
			backButton = backKeyboardButton(intrWhat)
		}
	} else if endpoint == intrPreview {
		backButton = backKeyboardButton(intrWhat)
	}
	if backButton != nil {
		lastButtonRow = append(lastButtonRow, *backButton)
	}

	if o.IsReady() && endpoint != intrPreview {
		lastButtonRow = append(lastButtonRow, previewButton())
	}

	replyMarkup.InlineKeyboard = append(
		replyMarkup.InlineKeyboard,
		lastButtonRow,
	)

	return api.EditMessageReplyMarkupConfig{
		BaseEdit: api.BaseEdit{
			ChatID:      msgInfo.Chat.ID,
			MessageID:   msgInfo.MessageID,
			ReplyMarkup: &replyMarkup,
		},
	}
}
