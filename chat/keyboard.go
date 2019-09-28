package chat

import (
	"time"

	"cafebot/config"
	"cafebot/order"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

func whereKeyboardFactory(conf config.CafeConfig) keyboardFunc {
	return func(intrData string, o order.Order) api.InlineKeyboardMarkup {
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

		return api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				hereButton,
				takeawayButton,
			),
		)
	}
}

func whenKeyboardFactory(conf config.CafeConfig) keyboardFunc {
	return func(intrData string, o order.Order) api.InlineKeyboardMarkup {
		nextIntr := intrWhat

		// we need everything except hour and minute to be 0
		now, _ := time.Parse("15:04", time.Now().Format("15:04"))

		interval := time.Duration(conf.TimeSlotIntervalMin) * time.Minute

		buttonRows := generateTimeSlotsKeyboard(
			nextIntr,
			generateTimeSlots(
				now,
				interval,
				time.Time(conf.FirstOrderTime),
				time.Time(conf.LastOrderTime),
			),
			o.Time,
		)
		return api.NewInlineKeyboardMarkup(buttonRows...)
	}
}

func whatKeyboardFactory(conf config.CafeConfig) keyboardFunc {
	return func(intrData string, o order.Order) api.InlineKeyboardMarkup {
		if intrData == "" {
			buttonRows := generateMenuCategoryButtonRows(
				conf.Menu.Categories,
				intrWhat,
			)
			return api.NewInlineKeyboardMarkup(buttonRows...)
		}

		meals, ok := conf.Menu.Map[intrData]
		if !ok {
			log.Errorf("invalid menu category requested: %s", intrData)
			return api.InlineKeyboardMarkup{}
		}

		buttonRows := generateMenuMealButtonRows(
			intrData,
			meals,
			intrWhat,
			o,
		)
		return api.NewInlineKeyboardMarkup(buttonRows...)
	}
}

func previewKeyboardFactory(conf config.CafeConfig) keyboardFunc {
	return func(intrData string, o order.Order) api.InlineKeyboardMarkup {
		return api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				api.NewInlineKeyboardButtonData(
					buttonText["send"],
					newIntrData(intrSent, "", opFinish, "1"),
				),
			),
		)
	}
}

func generateMenuCategoryButtonRows(
	categories []string,
	nextIntr intrEndpoint,
) [][]api.InlineKeyboardButton {
	buttonRows := make([][]api.InlineKeyboardButton, 0, len(categories))
	for _, category := range categories {
		intrData := newIntrData(nextIntr, category, "", "")
		buttonRows = append(
			buttonRows,
			[]api.InlineKeyboardButton{
				api.NewInlineKeyboardButtonData(
					category,
					intrData,
				),
			},
		)
	}
	return buttonRows
}

func generateMenuMealButtonRows(
	category string,
	meals []config.Meal,
	nextIntr intrEndpoint,
	o order.Order,
) [][]api.InlineKeyboardButton {
	buttonRows := make([][]api.InlineKeyboardButton, 0, len(meals))
	for _, meal := range meals {
		mealOp := opAddMeal
		selected := contains(o.Meal, meal.Val)
		if selected {
			mealOp = opRemoveMeal
		}

		buttonRows = append(
			buttonRows,
			[]api.InlineKeyboardButton{
				newIntrButton(
					meal.Val,
					newIntrData(
						nextIntr,
						category,
						mealOp,
						meal.Hash,
					),
					selected,
				),
			},
		)
	}
	return buttonRows
}

func backKeyboardButton(prevIntr intrEndpoint) api.InlineKeyboardButton {
	return api.NewInlineKeyboardButtonData(
		buttonText["back"],
		string(prevIntr),
	)
}

func previewButton() api.InlineKeyboardButton {
	return api.NewInlineKeyboardButtonData(
		buttonText["preview"],
		string(intrPreview),
	)
}

func newIntrButton(text, data string, selected bool) api.InlineKeyboardButton {
	if text == "" || data == "" {
		return api.InlineKeyboardButton{}
	}
	if selected {
		text = buttonText["selected"] + " " + text
	}
	return api.NewInlineKeyboardButtonData(text, data)
}

func contains(mm []string, m string) bool {
	for _, v := range mm {
		if v == m {
			return true
		}
	}
	return false
}
