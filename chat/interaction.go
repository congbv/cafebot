package chat

import (
	"strings"
	"sync"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/config"
	"github.com/yarikbratashchuk/cafebot/order"
)

type (
	// intrHandler is responsible for handling interactions via
	// bot message updates
	intrHandler struct {
		once sync.Once

		data     config.CafeConfig
		bot      *api.BotAPI
		handlers map[string]intrFunc
	}

	intrFunc func(string, api.Update, order.Order)
)

// intr is concurrent safe intrHandler singleton
var intr intrHandler

const (
	intrWhereEndpoint = "where"
	intrWhenEndpoint  = "when"
	intrWhatEndpoint  = "what"
)

func initIntrHandler(cafedata config.CafeConfig, bot *api.BotAPI) error {
	if bot == nil {
		return errNoAPI
	}
	intr.once.Do(func() {
		intr.bot = bot
		intr.data = cafedata
		intr.handlers = map[string]intrFunc{
			intrWhereEndpoint: intr.where,
			intrWhenEndpoint:  intr.when,
			intrWhatEndpoint:  intr.what,
		}
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
	h, ok := i.handlers[parts[0]]
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

	editText := api.EditMessageTextConfig{
		BaseEdit: api.BaseEdit{
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
		},
		Text: message["where?"],
	}

	_, err := i.bot.Send(editText)
	if err != nil {
		log.Error(err)
	}

	keyboardMarkup := api.NewInlineKeyboardMarkup(
		api.NewInlineKeyboardRow(
			api.NewInlineKeyboardButtonData(
				buttonText["here"],
				"order_when?takeaway=0",
			),
			api.NewInlineKeyboardButtonData(
				buttonText["takeaway"],
				"order_when?takeaway=1",
			),
		),
	)

	editKeyboard := api.EditMessageReplyMarkupConfig{
		BaseEdit: api.BaseEdit{
			ChatID:      update.CallbackQuery.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.MessageID,
			ReplyMarkup: &keyboardMarkup,
		},
	}

	i.bot.Send(editKeyboard)
}

func (i intrHandler) when(
	reqdata string,
	update api.Update,
	order order.Order,
) {
	log.Debug("intr.when")

	editText := api.EditMessageTextConfig{
		BaseEdit: api.BaseEdit{
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
		},
		Text: message["when?"],
	}

	_, err := i.bot.Send(editText)
	if err != nil {
		log.Error(err)
	}

	// we need everything except hour and minute to be 0
	now, _ := time.Parse("15:04", time.Now().Format("15:04"))

	timeSlots := generateTimeSlots(
		now,
		i.data.TimeSlotInterval,
		i.data.TimeSlotNumber,
		time.Time(i.data.OpenTime),
		time.Time(i.data.CloseTime),
	)
	timeRows := generateTimeSlotsKeyboard(
		timeSlots,
		i.data.TimeSlotNumInRow,
		intrWhatEndpoint,
	)

	keyboardRows := append(timeRows, backKeyboardButton("order_where"))

	keyboardMarkup := api.NewInlineKeyboardMarkup(keyboardRows...)

	editKeyboard := api.EditMessageReplyMarkupConfig{
		BaseEdit: api.BaseEdit{
			ChatID:      update.CallbackQuery.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.MessageID,
			ReplyMarkup: &keyboardMarkup,
		},
	}

	i.bot.Send(editKeyboard)
}

func (i intrHandler) what(
	reqdata string,
	update api.Update,
	order order.Order,
) {
	log.Debug("intr.what")

	// TODO: add stuff here
}
