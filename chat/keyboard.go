package chat

import (
	"bytes"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/config"
	"github.com/yarikbratashchuk/cafebot/order"
)

func whereKeyboardFactory(conf config.CafeConfig) keyboardFunc {
	return func(intrData string, o order.Order) *api.InlineKeyboardMarkup {
		nextIntr := intrWhen

		hereButton := newIntrButton(
			buttonText["here"],
			newIntrData(nextIntr, "", opWhere, opWhereHere),
			o.Takeaway != nil && *o.Takeaway == false,
		)
		takeawayButton := newIntrButton(
			buttonText["takeaway"],
			newIntrData(nextIntr, "", opWhere, opWhereTakeaway),
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
	return func(intrData string, o order.Order) *api.InlineKeyboardMarkup {
		nextIntr := intrWhat
		prevIntr := intrWhere

		// we need everything except hour and minute to be 0
		now, _ := time.Parse("15:04", time.Now().Format("15:04"))

		interval := time.Duration(conf.TimeSlotIntervalMin) * time.Minute

		buttonRows := append(
			generateTimeSlotsKeyboard(
				nextIntr,
				generateTimeSlots(
					now,
					interval,
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
	return func(intrData string, o order.Order) *api.InlineKeyboardMarkup {
		prevIntr := intrWhen

		if intrData == "" {
			categories := make([]string, 0, len(conf.Menu))
			for cat := range conf.Menu {
				categories = append(categories, cat)
			}
			buttonRows := append(
				generateMenuCategoryButtonRows(categories, intrWhat),
				backKeyboardButton(prevIntr),
			)
			keyboard := api.NewInlineKeyboardMarkup(buttonRows...)
			return &keyboard
		}

		// return category items and generate keyboards that add or
		// remove meal

		//buttons := append(
		//	generateMenuCategoryButtons(conf.Menu, intrWhat),
		//	backKeyboardButton(prevIntr),
		//	finishOrderButton(
		//)
		//keyboard := api.NewInlineKeyboardMarkup(buttons...)
		//return &keyboard

		return nil
	}
}

func generateMenuCategoryButtonRows(
	categories []string,
	endpoint intrEndpoint,
) [][]api.InlineKeyboardButton {
	buttonRows := make([][]api.InlineKeyboardButton, 0, len(categories))
	for _, category := range categories {
		buttonRows = append(
			buttonRows,
			[]api.InlineKeyboardButton{
				api.NewInlineKeyboardButtonData(
					category,
					newIntrData(endpoint, category, "", ""),
				),
			},
		)
	}
	return buttonRows
}

func backKeyboardButton(prevIntr intrEndpoint) []api.InlineKeyboardButton {
	return api.NewInlineKeyboardRow(
		api.NewInlineKeyboardButtonData(
			buttonText["back"],
			string(prevIntr),
		),
	)
}

func newIntrData(
	e intrEndpoint,
	intrData string,
	op orderOp,
	opval string,
) string {
	buf := new(bytes.Buffer)
	buf.WriteString(string(e))
	if intrData != "" {
		buf.WriteRune('/')
		buf.WriteString(intrData)
	}
	if op != "" && opval != "" {
		buf.WriteRune('?')
		buf.WriteString(string(op))
		buf.WriteRune('=')
		buf.WriteString(opval)
	}
	return buf.String()
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
