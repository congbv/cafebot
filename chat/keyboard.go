package chat

import (
	"fmt"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/config"
	"github.com/yarikbratashchuk/cafebot/order"
)

func whereKeyboardFactory(conf config.CafeConfig) keyboardFunc {
	return func(o order.Order) *api.InlineKeyboardMarkup {
		nextIntr := intrWhen

		hereButton := newIntrButton(
			buttonText["here"],
			newIntrData(nextIntr, opWhere, opWhereHere),
			o.Takeaway != nil && *o.Takeaway == false,
		)
		takeawayButton := newIntrButton(
			buttonText["takeaway"],
			newIntrData(nextIntr, opWhere, opWhereTakeaway),
			o.Takeaway != nil && *o.Takeaway == true,
		)

		keyboard := api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				hereButton,
				takeawayButton,
			),
		)

		return &keyboard
	}
}

func whenKeyboardFactory(conf config.CafeConfig) keyboardFunc {
	return func(o order.Order) *api.InlineKeyboardMarkup {
		nextIntr := intrWhat
		prevIntr := intrWhere

		// we need everything except hour and minute to be 0
		now, _ := time.Parse("15:04", time.Now().Format("15:04"))

		buttonRows := append(
			generateTimeSlotsKeyboard(
				nextIntr,
				generateTimeSlots(
					now,
					conf.TimeSlotInterval,
					time.Time(conf.FirstOrderTime),
					time.Time(conf.LastOrderTime),
				),
				o.Time,
			),
			backKeyboardButton(prevIntr),
		)
		keyboard := api.NewInlineKeyboardMarkup(buttonRows...)

		return &keyboard
	}
}

func whatKeyboardFactory(conf config.CafeConfig) keyboardFunc {
	return func(o order.Order) *api.InlineKeyboardMarkup {
		return nil
	}
}

func backKeyboardButton(prevIntr intrEndpoint) []api.InlineKeyboardButton {
	return api.NewInlineKeyboardRow(
		api.NewInlineKeyboardButtonData(
			buttonText["back"],
			string(prevIntr),
		),
	)
}

func newIntrData(e intrEndpoint, op orderOp, opval string) string {
	return fmt.Sprintf("%s?%s=%s", e, op, opval)
}

func newIntrButton(text, data string, selected bool) api.InlineKeyboardButton {
	if text == "" || data == "" {
		return api.InlineKeyboardButton{}
	}
	if selected {
		text = buttonText["selected"] + text
	}
	return api.NewInlineKeyboardButtonData(text, data)
}
